package monitor

import (
	"fmt"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/gofs/sync"
)

type Monitor interface {
	Start() error
	Close() error
}

func NewMonitor(syncer sync.Sync, retry retry.Retry) (Monitor, error) {
	src := syncer.Source()
	if src.IsDisk() || (src.Is(core.RemoteDisk) && src.Server()) {
		return NewFsNotifyMonitor(syncer, retry)
	} else if src.Is(core.RemoteDisk) && !src.Server() {
		return NewRemoteMonitor(syncer, retry, src.Host(), src.Port(), src.MessageQueue())
	}
	return nil, fmt.Errorf("file system unsupported ! src=>%s", src.Type().String())

}
