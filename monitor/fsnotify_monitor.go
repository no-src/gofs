package monitor

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/gofs/sync"
	"github.com/no-src/log"
	"io/fs"
	"path/filepath"
)

type fsNotifyMonitor struct {
	watcher *fsnotify.Watcher
	syncer  sync.Sync
	retry   retry.Retry
}

func NewFsNotifyMonitor(syncer sync.Sync, retry retry.Retry) (m Monitor, err error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	if syncer == nil {
		err = errors.New("syncer can't be nil")
		return nil, err
	}
	m = &fsNotifyMonitor{
		watcher: watcher,
		syncer:  syncer,
		retry:   retry,
	}
	return m, nil
}

func (m *fsNotifyMonitor) Monitor(dir string) (err error) {
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
	var last *fsnotify.Event
	for {
		select {
		case event, ok := <-m.watcher.Events:
			{
				if !ok {
					return errors.New("get watch event failed")
				}
				// in windows,will trigger twice, the same as last event,ignore it
				if last != nil && last.Name == event.Name && last.Op == event.Op {
					last = nil
					break
				} else {
					log.Debug("notify received [%s] -> [%s]", event.Op.String(), event.Name)
					last = &event
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					if m.retry != nil {
						m.retry.Do(func() error {
							return m.syncer.Write(event.Name)
						}, event.String())
					} else {
						m.syncer.Write(event.Name)
					}
				} else if event.Op&fsnotify.Create == fsnotify.Create {
					err := m.syncer.Create(event.Name)
					if err == nil {
						// if create a new dir, then monitor it
						isDir, err := m.syncer.IsDir(event.Name)
						if err == nil && isDir {
							m.Monitor(event.Name)
						}
					}
				} else if event.Op&fsnotify.Remove == fsnotify.Remove {
					m.syncer.Remove(event.Name)
				} else if event.Op&fsnotify.Rename == fsnotify.Rename {
					m.syncer.Rename(event.Name)
				} else if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					m.syncer.Chmod(event.Name)
				}
				break
			}

		case err, ok := <-m.watcher.Errors:
			{
				if !ok {
					return errors.New("get watch error failed")
				}
				log.Error(err, "watcher error")
				break
			}

		}
	}
}

func (m *fsNotifyMonitor) Close() error {
	if m.watcher != nil {
		return m.watcher.Close()
	}
	return nil
}
