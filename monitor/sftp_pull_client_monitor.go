package monitor

import (
	"io"
	"time"

	"github.com/no-src/gofs/retry"
	"github.com/no-src/gofs/sync"
)

type sftpPullClientMonitor struct {
	driverPullClientMonitor
}

// NewSftpPullClientMonitor create an instance of sftpPullClientMonitor to pull the files from sftp server
func NewSftpPullClientMonitor(syncer sync.Sync, retry retry.Retry, syncOnce bool, eventWriter io.Writer, enableSyncDelay bool, syncDelayEvents int, syncDelayTime time.Duration) (m Monitor, err error) {
	m = &sftpPullClientMonitor{
		driverPullClientMonitor: driverPullClientMonitor{
			baseMonitor: newBaseMonitor(syncer, retry, syncOnce, eventWriter, enableSyncDelay, syncDelayEvents, syncDelayTime),
		},
	}
	return m, nil
}
