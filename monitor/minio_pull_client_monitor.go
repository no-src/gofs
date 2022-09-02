package monitor

import (
	"io"
	"time"

	"github.com/no-src/gofs/retry"
	"github.com/no-src/gofs/sync"
)

type minIOPullClientMonitor struct {
	driverPullClientMonitor
}

// NewMinIOPullClientMonitor create an instance of minIOPullClientMonitor to pull the files from MinIO server
func NewMinIOPullClientMonitor(syncer sync.Sync, retry retry.Retry, syncOnce bool, eventWriter io.Writer, enableSyncDelay bool, syncDelayEvents int, syncDelayTime time.Duration) (m Monitor, err error) {
	m = &minIOPullClientMonitor{
		driverPullClientMonitor: driverPullClientMonitor{
			baseMonitor: newBaseMonitor(syncer, retry, syncOnce, eventWriter, enableSyncDelay, syncDelayEvents, syncDelayTime),
		},
	}
	return m, nil
}
