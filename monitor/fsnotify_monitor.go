package monitor

import (
	"errors"
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
	"github.com/no-src/gofs/wait"
	"github.com/no-src/log"
)

type fsNotifyMonitor struct {
	baseMonitor

	watcher *fsnotify.Watcher
	events  *clist.CList
	pi      ignore.PathIgnore
}

// NewFsNotifyMonitor create an instance of fsNotifyMonitor to monitor the disk change
func NewFsNotifyMonitor(opt Option) (m Monitor, err error) {
	pi := opt.PathIgnore

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	if opt.Syncer == nil {
		err = errors.New("syncer can't be nil")
		return nil, err
	}

	m = &fsNotifyMonitor{
		watcher:     watcher,
		baseMonitor: newBaseMonitor(opt),
		events:      clist.New(),
		pi:          pi,
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

func (m *fsNotifyMonitor) Start() (wait.Wait, error) {
	source := m.syncer.Source()
	wd := wait.NewWaitDone()

	// execute -sync_once flag
	if m.syncOnce {
		wd.Done()
		return wd, m.syncer.SyncOnce(source.Path())
	}

	// execute -sync_cron flag
	if err := m.startCron(func() error {
		return m.syncer.SyncOnce(source.Path())
	}); err != nil {
		return nil, err
	}

	if !source.IsDisk() && !source.Is(core.RemoteDisk) {
		return nil, errors.New("the source must be a disk or remote disk")
	}
	if err := m.monitor(source.Path()); err != nil {
		return nil, err
	}

	go m.startReceiveWriteNotify()
	go m.startSyncWrite()
	go m.startProcessEvents()
	go m.startReceiveEvents(wd)
	return wd, nil
}

// startReceiveEvents start loop to receive file change event from the fsnotify
func (m *fsNotifyMonitor) startReceiveEvents(wd wait.Done) error {
	for {
		select {
		case event, ok := <-m.watcher.Events:
			{
				if !ok {
					err := errors.New("get fsnotify watch event failed")
					wd.DoneWithError(err)
					return err
				}
				log.Debug("notify received [%s] -> [%s]", event.Op.String(), event.Name)
				m.events.PushBack(event)
			}
		case err, ok := <-m.watcher.Errors:
			{
				if !ok {
					err = errors.New("get watch error failed")
					wd.DoneWithError(err)
					return err
				}
				log.Error(err, "watcher error")
			}
		case shutdown := <-m.shutdown:
			{
				if shutdown {
					wd.Done()
					return nil
				}
			}
		}
	}
}

// startProcessEvents start loop to process all file change events
func (m *fsNotifyMonitor) startProcessEvents() error {
	for {
		m.waitSyncDelay(m.events.Len)

		element := m.events.Front()
		if element == nil || element.Value == nil {
			if element != nil {
				m.events.Remove(element)
			}
			m.resetSyncDelay()
			<-time.After(time.Second)
			continue
		}
		event := element.Value.(fsnotify.Event)
		if m.pi.MatchPath(event.Name, "monitor", event.Op.String()) {
			// if the rule is matched, then ignore the event except create a directory, because of the subdirectory maybe not match the ignore rule.
			// so we should monitor the current directory here, otherwise we will lose some data.
			// for example, we define an ignore rule "/home/logs/*" and create a directory "/home/logs" to trigger Create event, then create a file "/home/logs/2022/info.log".
			// the file "info.log" does not match the ignore rule and should be synchronized to destination directory.
			m.monitorDirIfCreate(event)
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
		log.Error(err, "[write] event execute create error => [%s]", event.Name)
	}

	var size int64
	stat, err := os.Stat(event.Name)
	if err == nil {
		size = stat.Size()
	}
	m.addWrite(event.Name, size)
}

func (m *fsNotifyMonitor) create(event fsnotify.Event) {
	err := m.syncer.Create(event.Name)
	if err == nil {
		// if create a new dir, then monitor it
		isDir, err := m.syncer.IsDir(event.Name)
		if err == nil && isDir {
			if err = m.monitor(event.Name); err != nil {
				log.Error(err, "[create] event execute monitor error => [%s]", event.Name)
			}
		}
		if err == nil {
			// rename a file will not trigger the Write event
			// rename a dir will not trigger the Write event on Linux and some Windows environments
			// in some cases it will trigger the Write event for the parent dir on Windows
			// send a Write event manually
			log.Debug("prepare to send a [write] event after [create] event [%s]", event.Name)
			m.events.PushBack(fsnotify.Event{
				Name: event.Name,
				Op:   fsnotify.Write,
			})
		}
	}
}

func (m *fsNotifyMonitor) remove(event fsnotify.Event) {
	m.removeWrite(event.Name)
	log.ErrorIf(m.syncer.Remove(event.Name), "[remove] event execute error => [%s]", event.Name)
}

func (m *fsNotifyMonitor) rename(event fsnotify.Event) {
	log.ErrorIf(m.syncer.Rename(event.Name), "[rename] event execute error => [%s]", event.Name)
}

func (m *fsNotifyMonitor) chmod(event fsnotify.Event) {
	log.ErrorIf(m.syncer.Chmod(event.Name), "[chmod] event execute error => [%s]", event.Name)
}

// monitorDirIfCreate monitor the directory if you create a new directory
func (m *fsNotifyMonitor) monitorDirIfCreate(event fsnotify.Event) {
	if event.Op&fsnotify.Create == fsnotify.Create {
		isDir, err := m.syncer.IsDir(event.Name)
		if log.ErrorIf(err, "check path is dir or not error") != nil {
			return
		}
		if isDir {
			log.ErrorIf(m.monitor(event.Name), "monitor the directory error that matched ignore rule => [%s]", event.Name)
		}
	}
}

func (m *fsNotifyMonitor) Close() error {
	if m.watcher != nil {
		return m.watcher.Close()
	}
	return nil
}
