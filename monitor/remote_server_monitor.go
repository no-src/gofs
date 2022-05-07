package monitor

import (
	"io"
	"time"

	"github.com/no-src/gofs/retry"
	"github.com/no-src/gofs/sync"
)

// NewRemoteServerMonitor create an instance of the fsNotifyMonitor
func NewRemoteServerMonitor(syncer sync.Sync, retry retry.Retry, syncOnce bool, eventWriter io.Writer, enableSyncDelay bool, syncDelayEvents int, syncDelayTime time.Duration) (m Monitor, err error) {
	return NewFsNotifyMonitor(syncer, retry, syncOnce, eventWriter, enableSyncDelay, syncDelayEvents, syncDelayTime)
}
