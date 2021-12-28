package monitor

import (
	"fmt"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/gofs/sync"
	"github.com/no-src/log"
	"os"
	"sort"
	goSync "sync"
	"time"
)

type baseMonitor struct {
	syncer      sync.Sync
	retry       retry.Retry
	writeMap    map[string]*writeMessage
	writeList   writeMessageList
	writeChan   chan *writeMessage
	writeNotify chan bool
	mu          goSync.Mutex
}

func newBaseMonitor(syncer sync.Sync, retry retry.Retry) baseMonitor {
	return baseMonitor{
		syncer:      syncer,
		retry:       retry,
		writeMap:    make(map[string]*writeMessage),
		writeChan:   make(chan *writeMessage, 100),
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

// removeWrite remove write message
func (m *baseMonitor) removeWrite(name string) {
	m.mu.Lock()
	wm := m.writeMap[name]
	if wm != nil {
		wm.cancel = true
		delete(m.writeMap, name)
		log.Debug("removeWrite => [%s]", name)
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
			if wm != nil && !wm.cancel {
				if (wm.count <= 2 && now-wm.last <= time.Second.Nanoseconds()) || (wm.count > 2 && now-wm.last <= 3*time.Second.Nanoseconds()) {
					m.mu.Unlock()
					go func() {
						<-time.After(time.Second)
						m.writeNotify <- true
					}()
					continue
				}
				go func() {
					m.writeChan <- wm
				}()

				delete(m.writeMap, wm.name)
			}
			m.writeList = m.writeList[1:]
		}
		m.mu.Unlock()
	}
}

// startSyncWrite write file to sync
func (m *baseMonitor) startSyncWrite() {
	for {
		wm := <-m.writeChan
		if wm == nil || wm.cancel {
			continue
		}
		name := wm.name
		if m.retry != nil {
			m.retry.Do(func() error {
				err := m.syncer.Write(name)
				// if file or directory is not exist, ignore it and warning
				if os.IsNotExist(err) {
					log.Warn("write file failed => [%s]", err.Error())
					return nil
				}
				return err
			}, fmt.Sprintf("write file => %s", name))
		} else {
			if err := m.syncer.Write(name); err != nil {
				log.Error(err, "write file error => [%s]", name)
			}
		}
	}
}
