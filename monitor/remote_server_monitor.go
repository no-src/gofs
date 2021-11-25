package monitor

import (
	"github.com/no-src/gofs/retry"
	"github.com/no-src/gofs/sync"
)

func NewRemoteServerMonitor(syncer sync.Sync, retry retry.Retry, syncOnce bool) (m Monitor, err error) {
	return NewFsNotifyMonitor(syncer, retry, syncOnce)
}
