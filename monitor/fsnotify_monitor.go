package monitor

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/eventlog"
	"github.com/no-src/gofs/ignore"
	"github.com/no-src/gofs/internal/clist"
	"github.com/no-src/gofs/report"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/gofs/sync"
	"github.com/no-src/gofs/util/osutil"
	"github.com/no-src/log"
)

type fsNotifyMonitor struct {
	baseMonitor

	watcher *fsnotify.Watcher
	events  *clist.CList
}

// NewFsNotifyMonitor create an instance of fsNotifyMonitor to monitor the disk change
func NewFsNotifyMonitor(syncer sync.Sync, retry retry.Retry, syncOnce bool, eventWriter io.Writer) (m Monitor, err error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	if syncer == nil {
		err = errors.New("syncer can't be nil")
		return nil, err
	}
	m = &fsNotifyMonitor{
		watcher:     watcher,
		baseMonitor: newBaseMonitor(syncer, retry, syncOnce, eventWriter),
		events:      clist.New(),
	}
	return m, nil
}

func (m *fsNotifyMonitor) monitor(dir string) (err error) {
	dir, err = filepath.Abs(dir)
	if err != nil {
		log.Error(err, "parse dir to abs dir error")
		return err
	}
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			// first remove the old watch, because the volume is the same as the one before renamed,
			// then add path to watch.
			m.watcher.Remove(path)
			err = m.watcher.Add(path)
			if err != nil {
				log.Error(err, "watch dir error [%s]", path)
			} else {
				log.Debug("watch dir success [%s]", path)
			}
		}
		return err
	})
	if err != nil {
		log.Error(err, "monitor dir error [%s]", dir)
	}
	return err
}

func (m *fsNotifyMonitor) Start() error {
	source := m.syncer.Source()
	// execute -sync_once flag
	if m.syncOnce {
		return m.syncer.SyncOnce(source.Path())
	}

	// execute -sync_cron flag
	if err := m.startCron(func() error {
		return m.syncer.SyncOnce(source.Path())
	}); err != nil {
		return err
	}

	if !source.IsDisk() && !source.Is(core.RemoteDisk) {
		return errors.New("the source must be a disk or remote disk")
	}
	if err := m.monitor(source.Path()); err != nil {
		return err
	}

	go m.startReceiveWriteNotify()
	go m.startSyncWrite()
	go m.startProcessEvents()

	return m.startReceiveEvents()
}

// startReceiveEvents start loop to receive file change event from the fsnotify
func (m *fsNotifyMonitor) startReceiveEvents() error {
	for {
		select {
		case event, ok := <-m.watcher.Events:
			{
				if !ok {
					return errors.New("get fsnotify watch event failed")
				}
				log.Debug("notify received [%s] -> [%s]", event.Op.String(), event.Name)
				m.events.PushBack(event)
			}
		case err, ok := <-m.watcher.Errors:
			{
				if !ok {
					return errors.New("get watch error failed")
				}
				log.Error(err, "watcher error")
			}
		case shutdown := <-m.shutdown:
			{
				if shutdown {
					return nil
				}
			}
		}
	}
}

// startProcessEvents start loop to process all file change events
func (m *fsNotifyMonitor) startProcessEvents() error {
	for {
		element := m.events.Front()
		if element == nil || element.Value == nil {
			if element != nil {
				m.events.Remove(element)
			}
			<-time.After(time.Second)
			continue
		}
		event := element.Value.(fsnotify.Event)
		if ignore.MatchPath(event.Name, "monitor", event.Op.String()) {
			// ignore match
		} else if event.Op&fsnotify.Write == fsnotify.Write {
			m.write(event)
		} else if event.Op&fsnotify.Create == fsnotify.Create {
			m.create(event)
		} else if event.Op&fsnotify.Remove == fsnotify.Remove {
			m.remove(event)
		} else if event.Op&fsnotify.Rename == fsnotify.Rename {
			m.rename(event)
		} else if event.Op&fsnotify.Chmod == fsnotify.Chmod {
			m.chmod(event)
		}
		m.events.Remove(element)
		e := eventlog.NewEvent(event.Name, event.Op.String())
		m.el.Write(e)
		report.GlobalReporter.PutEvent(e)
	}
}

func (m *fsNotifyMonitor) write(event fsnotify.Event) {
	// ignore is not exist error
	if err := m.syncer.Create(event.Name); err != nil && !os.IsNotExist(err) {
		log.Error(err, "Write event execute create error => [%s]", event.Name)
	}
	m.addWrite(event.Name)
}

func (m *fsNotifyMonitor) create(event fsnotify.Event) {
	err := m.syncer.Create(event.Name)
	if err == nil {
		// if create a new dir, then monitor it
		isDir, err := m.syncer.IsDir(event.Name)
		if err == nil && isDir {
			if err = m.monitor(event.Name); err != nil {
				log.Error(err, "Create event execute monitor error => [%s]", event.Name)
			}
		}
		if err == nil && (!isDir || (isDir && !osutil.IsWindows())) {
			// rename a file, will not trigger the Write event
			// rename a dir, will not trigger the Write event on Linux, but it will trigger the Write event for parent dir on Windows
			// send a Write event manually
			log.Debug("prepare to send a Write event after Create event [%s]", event.Name)
			m.events.PushBack(fsnotify.Event{
				Name: event.Name,
				Op:   fsnotify.Write,
			})
		}
	}
}

func (m *fsNotifyMonitor) remove(event fsnotify.Event) {
	m.removeWrite(event.Name)
	log.ErrorIf(m.syncer.Remove(event.Name), "Remove event execute error => [%s]", event.Name)
}

func (m *fsNotifyMonitor) rename(event fsnotify.Event) {
	log.ErrorIf(m.syncer.Rename(event.Name), "Rename event execute error => [%s]", event.Name)
}

func (m *fsNotifyMonitor) chmod(event fsnotify.Event) {
	log.ErrorIf(m.syncer.Chmod(event.Name), "Chmod event execute error => [%s]", event.Name)
}

func (m *fsNotifyMonitor) Close() error {
	if m.watcher != nil {
		return m.watcher.Close()
	}
	return nil
}
