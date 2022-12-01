package monitor

import (
	"fmt"

	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/wait"
)

// Monitor file system monitor
type Monitor interface {
	// Start go to start the monitor to monitor the file change
	Start() (wait.Wait, error)
	// Close stop the monitor
	Close() error
	// SyncCron register sync cron task, if spec is empty then ignore it
	SyncCron(spec string) error
	// Shutdown exit the Start
	Shutdown() error
}

// NewMonitor create a monitor instance
func NewMonitor(opt Option) (Monitor, error) {
	source := opt.Syncer.Source()
	if source.IsDisk() {
		return NewFsNotifyMonitor(opt)
	} else if source.Is(core.RemoteDisk) && source.Server() {
		return NewRemoteServerMonitor(opt)
	} else if source.Is(core.RemoteDisk) && !source.Server() {
		return NewRemoteClientMonitor(opt)
	} else if source.Is(core.SFTP) {
		return NewSftpPullClientMonitor(opt)
	} else if source.Is(core.MinIO) {
		return NewMinIOPullClientMonitor(opt)
	}
	return nil, fmt.Errorf("file system unsupported ! source=>%s", source.Type().String())
}
