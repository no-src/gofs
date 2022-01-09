package monitor

import (
	"errors"
	"fmt"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/gofs/sync"
	"github.com/no-src/log"
	"github.com/robfig/cron/v3"
	"net/url"
	"os"
	"sort"
	"strings"
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
	syncSpec    string
	cronChan    chan bool
	shutdown    chan bool
}

func newBaseMonitor(syncer sync.Sync, retry retry.Retry) baseMonitor {
	return baseMonitor{
		syncer:      syncer,
		retry:       retry,
		writeMap:    make(map[string]*writeMessage),
		writeChan:   make(chan *writeMessage, 100),
		writeNotify: make(chan bool, 100),
		cronChan:    make(chan bool, 1),
		shutdown:    make(chan bool, 1),
	}
}

// addWrite add or update a write message
func (m *baseMonitor) addWrite(name string) {
	m.mu.Lock()
	wm := m.writeMap[m.key(name)]
	if wm == nil {
		wm = newDefaultWriteMessage(name)
		m.writeMap[m.key(name)] = wm
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
	wm := m.writeMap[m.key(name)]
	if wm != nil {
		wm.cancel = true
		delete(m.writeMap, m.key(name))
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

				delete(m.writeMap, m.key(wm.name))
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
				// if file or directory is not exist, ignore it
				if os.IsNotExist(err) {
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

// key return file identity as hash key, that removes the query section if the file name is an url
func (m *baseMonitor) key(name string) string {
	u, err := url.Parse(name)
	if err != nil || u == nil || len(u.RawQuery) == 0 {
		return name
	}
	return strings.ReplaceAll(name, "?"+u.RawQuery, "")
}

func (m *baseMonitor) startCron(f func() error) error {
	if len(m.syncSpec) == 0 {
		return nil
	}
	c := cron.New(cron.WithSeconds())
	id, err := c.AddFunc(m.syncSpec, func() {
		defer func() {
			<-m.cronChan
			if e := recover(); e != nil {
				log.Error(errors.New("cron task execute panic"), "%v", e)
			}
		}()
		m.cronChan <- true
		log.Info("start execute cron task, spec=[%s]", m.syncSpec)
		err := f()
		if err != nil {
			log.Error(err, "execute cron error spec=[%s]", m.syncSpec)
		} else {
			log.Info("execute cron task finished, spec=[%s]", m.syncSpec)
		}
	})
	if err != nil {
		return err
	}
	log.Info("cron task starting, spec=[%s] id=[%d]", m.syncSpec, id)
	c.Start()
	return nil
}

func (m *baseMonitor) SyncCron(spec string) error {
	spec = strings.TrimSpace(spec)
	if len(spec) == 0 {
		return nil
	}
	parser := cron.NewParser(
		cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)
	_, err := parser.Parse(spec)
	if err == nil {
		m.syncSpec = spec
	}
	return err
}

func (m *baseMonitor) Shutdown() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	m.shutdown <- true
	close(m.shutdown)
	return err
}
