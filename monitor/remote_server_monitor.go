package monitor

import (
	"github.com/no-src/gofs/retry"
	"github.com/no-src/gofs/sync"
	"io"
)

func NewRemoteServerMonitor(syncer sync.Sync, retry retry.Retry, syncOnce bool, eventWriter io.Writer) (m Monitor, err error) {
	return NewFsNotifyMonitor(syncer, retry, syncOnce, eventWriter)
}
