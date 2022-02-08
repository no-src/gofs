package monitor

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/eventlog"
	"github.com/no-src/gofs/ignore"
	"github.com/no-src/gofs/internal/clist"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/gofs/sync"
	"github.com/no-src/gofs/util"
	"github.com/no-src/log"
	"io"
	iofs "io/fs"
	"os"
	"path/filepath"
	"time"
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
	err = filepath.WalkDir(dir, func(path string, d iofs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			// first remove the old watch, because the volume is the same as the one before renamed,
			// then add path to watch.
			err = m.watcher.Remove(path)
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
		return errors.New("not local file system")
	}
	if err := m.monitor(source.Path()); err != nil {
		return err
	}

	go m.processWrite()
	go m.startSyncWrite()
	go m.processEvents()

	return m.listenEvents()
}

func (m *fsNotifyMonitor) listenEvents() error {
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

func (m *fsNotifyMonitor) processEvents() error {
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
			// ignore is not exist error
			if err := m.syncer.Create(event.Name); err != nil && !os.IsNotExist(err) {
				log.Error(err, "Write event execute create error => [%s]", event.Name)
			}
			m.addWrite(event.Name)
		} else if event.Op&fsnotify.Create == fsnotify.Create {
			err := m.syncer.Create(event.Name)
			if err == nil {
				// if create a new dir, then monitor it
				isDir, err := m.syncer.IsDir(event.Name)
				if err == nil && isDir {
					if err = m.monitor(event.Name); err != nil {
						log.Error(err, "Create event execute monitor error => [%s]", event.Name)
					}
				}
				if err == nil && (!isDir || (isDir && !util.IsWindows())) {
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
		} else if event.Op&fsnotify.Remove == fsnotify.Remove {
			m.removeWrite(event.Name)
			if err := m.syncer.Remove(event.Name); err != nil {
				log.Error(err, "Remove event execute error => [%s]", event.Name)
			}
		} else if event.Op&fsnotify.Rename == fsnotify.Rename {
			if err := m.syncer.Rename(event.Name); err != nil {
				log.Error(err, "Rename event execute error => [%s]", event.Name)
			}
		} else if event.Op&fsnotify.Chmod == fsnotify.Chmod {
			if err := m.syncer.Chmod(event.Name); err != nil {
				log.Error(err, "Chmod event execute error => [%s]", event.Name)
			}
		}
		m.events.Remove(element)
		m.el.Write(eventlog.NewEvent(event.Name, event.Op.String()))
	}
}

func (m *fsNotifyMonitor) Close() error {
	if m.watcher != nil {
		return m.watcher.Close()
	}
	return nil
}
