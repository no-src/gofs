package monitor

import (
	"fmt"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/gofs/sync"
	"github.com/no-src/log"
	"sort"
	goSync "sync"
	"time"
)

type baseMonitor struct {
	syncer      sync.Sync
	retry       retry.Retry
	writeMap    map[string]*writeMessage
	writeList   writeMessageList
	writeChan   chan string
	writeNotify chan bool
	mu          goSync.Mutex
}

func newBaseMonitor(syncer sync.Sync, retry retry.Retry) baseMonitor {
	return baseMonitor{
		syncer:      syncer,
		retry:       retry,
		writeMap:    make(map[string]*writeMessage),
		writeChan:   make(chan string, 100),
		writeNotify: make(chan bool, 100),
	}
}

// addWrite add or update a write message
func (m *baseMonitor) addWrite(name string) {
	m.mu.Lock()
	wm := m.writeMap[name]
	if wm == nil {
		wm = newDefaultWriteMessage(name)
		m.writeMap[name] = wm
		m.writeList = append(m.writeList, wm)
	} else {
		wm.count++
		wm.last = time.Now().UnixNano()
		if len(m.writeList) > 0 {
			sort.Sort(m.writeList)
		}
	}
	m.mu.Unlock()
	m.writeNotify <- true
}

// processWrite process write list
func (m *baseMonitor) processWrite() {
	for {
		select {
		case <-m.writeNotify:
		case <-time.After(time.Second):
		}
		m.mu.Lock()
		now := time.Now().UnixNano()
		if len(m.writeList) > 0 {
			wm := m.writeList[0]
			if wm != nil {
				if (wm.count <= 2 && now-wm.last <= time.Second.Nanoseconds()) || (wm.count > 2 && now-wm.last <= 3*time.Second.Nanoseconds()) {
					m.mu.Unlock()
					continue
				}
				go func() {
					m.writeChan <- wm.name
				}()
			}

			m.writeList = m.writeList[1:]
			delete(m.writeMap, wm.name)
		}
		m.mu.Unlock()
	}
}

// startSyncWrite write file to sync
func (m *baseMonitor) startSyncWrite() {
	for {
		name := <-m.writeChan
		if m.retry != nil {
			m.retry.Do(func() error {
				return m.syncer.Write(name)
			}, fmt.Sprintf("write file => %s", name))
		} else {
			if err := m.syncer.Write(name); err != nil {
				log.Error(err, "write file error => [%s]", name)
			}
		}
	}
}
