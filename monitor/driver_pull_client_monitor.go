package monitor

import (
	"github.com/no-src/gofs/internal/cbool"
	"github.com/no-src/gofs/wait"
	"github.com/no-src/log"
)

type driverPullClientMonitor struct {
	baseMonitor
}

func (m *driverPullClientMonitor) Start() (wait.Wait, error) {
	wd := wait.NewWaitDone()
	shutdown := cbool.New(false)
	go m.waitShutdown(shutdown, wd)

	// execute -sync_once flag
	if m.syncOnce {
		return wd, m.syncAndShutdown()
	}

	// execute -sync_cron flag
	if err := m.startCron(m.sync); err != nil {
		return nil, err
	}

	return wd, nil
}

// syncAndShutdown execute sync and then try to shut down, the caller should wait for shutdown by wait.Wait()
func (m *driverPullClientMonitor) syncAndShutdown() (err error) {
	if err = m.sync(); err != nil {
		return err
	}
	if err = m.Shutdown(); err != nil {
		return err
	}
	return nil
}

// waitShutdown wait for the shutdown notify then mark the work done
func (m *driverPullClientMonitor) waitShutdown(st *cbool.CBool, wd wait.Done) {
	select {
	case <-st.SetC(<-m.shutdown):
		{
			if st.Get() {
				log.ErrorIf(m.Close(), "close driver pull client monitor error")
				m.syncer.Close()
				wd.Done()
			}
		}
	}
}

// sync try to sync all the files once
func (m *driverPullClientMonitor) sync() (err error) {
	// source.Path() and source.RemotePath() are equivalent here, and source.RemotePath() has higher priority
	source := m.syncer.Source()
	path := source.RemotePath()
	if len(path) == 0 {
		path = source.Path()
	}
	return m.syncer.SyncOnce(path)
}

func (m *driverPullClientMonitor) Close() error {
	return nil
}
