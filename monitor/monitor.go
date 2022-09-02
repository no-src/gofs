package monitor

import (
	"fmt"
	"io"
	"time"

	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/gofs/sync"
)

// Monitor file system monitor
type Monitor interface {
	// Start go to start the monitor to monitor the file change
	Start() error
	// Close stop the monitor
	Close() error
	// SyncCron register sync cron task, if spec is empty then ignore it
	SyncCron(spec string) error
	// Shutdown exit the Start
	Shutdown() error
}

// NewMonitor create a monitor instance
// syncer a Sync component
// retry a Retry component
// syncOnce tag a sync once command, the sync once command will execute when call the Start
func NewMonitor(syncer sync.Sync, retry retry.Retry, syncOnce bool, enableTLS bool, certFile string, insecureSkipVerify bool, users []*auth.User, eventWriter io.Writer, enableSyncDelay bool, syncDelayEvents int, syncDelayTime time.Duration) (Monitor, error) {
	source := syncer.Source()
	if source.IsDisk() {
		return NewFsNotifyMonitor(syncer, retry, syncOnce, eventWriter, enableSyncDelay, syncDelayEvents, syncDelayTime)
	} else if source.Is(core.RemoteDisk) && source.Server() {
		return NewRemoteServerMonitor(syncer, retry, syncOnce, eventWriter, enableSyncDelay, syncDelayEvents, syncDelayTime)
	} else if source.Is(core.RemoteDisk) && !source.Server() {
		return NewRemoteClientMonitor(syncer, retry, syncOnce, source.Host(), source.Port(), enableTLS, certFile, insecureSkipVerify, users, eventWriter, enableSyncDelay, syncDelayEvents, syncDelayTime)
	} else if source.Is(core.SFTP) {
		return NewSftpPullClientMonitor(syncer, retry, syncOnce, eventWriter, enableSyncDelay, syncDelayEvents, syncDelayTime)
	} else if source.Is(core.MinIO) {
		return NewMinIOPullClientMonitor(syncer, retry, syncOnce, eventWriter, enableSyncDelay, syncDelayEvents, syncDelayTime)
	}
	return nil, fmt.Errorf("file system unsupported ! source=>%s", source.Type().String())
}
