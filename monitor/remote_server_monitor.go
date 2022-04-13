package monitor

import (
	"io"

	"github.com/no-src/gofs/retry"
	"github.com/no-src/gofs/sync"
)

// NewRemoteServerMonitor create an instance of the fsNotifyMonitor
func NewRemoteServerMonitor(syncer sync.Sync, retry retry.Retry, syncOnce bool, eventWriter io.Writer) (m Monitor, err error) {
	return NewFsNotifyMonitor(syncer, retry, syncOnce, eventWriter)
}
