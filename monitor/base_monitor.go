package monitor

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/no-src/gofs/eventlog"
	"github.com/no-src/gofs/retry"
	nssync "github.com/no-src/gofs/sync"
	"github.com/no-src/log"
	"github.com/robfig/cron/v3"
)

type baseMonitor struct {
	syncer          nssync.Sync
	retry           retry.Retry
	writeMap        map[string]*writeMessage
	writeList       writeMessageList
	writeChan       chan *writeMessage
	writeNotify     chan struct{}
	mu              sync.Mutex
	syncSpec        string
	cronChan        chan struct{}
	shutdown        chan bool
	syncOnce        bool
	el              eventlog.EventLog
	enableSyncDelay bool
	syncDelayEvents int
	syncDelayTime   time.Duration
	lastSyncTime    time.Time
	syncing         bool
}

func newBaseMonitor(syncer nssync.Sync, retry retry.Retry, syncOnce bool, eventWriter io.Writer, enableSyncDelay bool, syncDelayEvents int, syncDelayTime time.Duration) baseMonitor {
	return baseMonitor{
		syncer:          syncer,
		retry:           retry,
		writeMap:        make(map[string]*writeMessage),
		writeChan:       make(chan *writeMessage, 100),
		writeNotify:     make(chan struct{}, 100),
		cronChan:        make(chan struct{}, 1),
		shutdown:        make(chan bool, 1),
		syncOnce:        syncOnce,
		el:              eventlog.New(eventWriter),
		enableSyncDelay: enableSyncDelay,
		syncDelayEvents: syncDelayEvents,
		syncDelayTime:   syncDelayTime,
		lastSyncTime:    time.Now(),
		syncing:         !enableSyncDelay,
	}
}

// addWrite add or update a write message
func (m *baseMonitor) addWrite(name string) {
	m.mu.Lock()

	// If the current path's parent directory is in the writeMap, then ignore the current path.
	// For example
	// WRITE /source/workspace
	// WRITE /source/workspace/hello.txt
	// As the above says, ignore the path /source/workspace/hello.txt
	parent := filepath.Dir(name)
	pwm := m.writeMap[m.key(parent)]
	if pwm != nil {
		log.Debug("add write ignore the file path => %s", name)
		m.mu.Unlock()
		return
	}

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
	m.writeNotify <- struct{}{}
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
	m.writeNotify <- struct{}{}
}

// startReceiveWriteNotify start loop to receive write notification, and delay process
func (m *baseMonitor) startReceiveWriteNotify() {
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
					<-time.After(time.Second)
					go func() {
						m.writeNotify <- struct{}{}
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

// startSyncWrite start loop to receive a write message and process it right now
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
				if err != nil {
					log.Error(err, "write file error => [%s]", name)
				}
				return err
			}, fmt.Sprintf("write file => %s", name))
		} else {
			log.ErrorIf(m.syncer.Write(name), "write file error => [%s]", name)
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
		m.cronChan <- struct{}{}
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
	return err
}

func (m *baseMonitor) waitSyncDelay(eventLenFunc func() int) {
	for {
		if m.enableSyncDelay && !m.syncing {
			currentEvents := eventLenFunc()
			if currentEvents > 0 {
				if currentEvents < m.syncDelayEvents && time.Now().Before(m.lastSyncTime.Add(m.syncDelayTime)) {
					log.DebugSample("[sync delay] [waiting] sync delay time => %s, sync delay events => %d, last sync time => %s, current event count => %d ", m.syncDelayTime, m.syncDelayEvents, m.lastSyncTime, currentEvents)
					<-time.After(time.Second)
					continue
				}
				log.Debug("[sync delay] [starting] sync delay time => %s, sync delay events => %d, last sync time => %s, current event count => %d ", m.syncDelayTime, m.syncDelayEvents, m.lastSyncTime, currentEvents)
				m.syncing = true
			}
		}
		break
	}
}

func (m *baseMonitor) resetSyncDelay() {
	m.lastSyncTime = time.Now()
	if m.enableSyncDelay {
		syncing := m.syncing
		m.syncing = false
		if syncing {
			log.Debug("[sync delay] [reset] sync delay time => %s, sync delay events => %d, last sync time => %s ", m.syncDelayTime, m.syncDelayEvents, m.lastSyncTime)
		}
	} else {
		m.syncing = true
	}
}
