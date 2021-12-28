package monitor

import (
	"fmt"
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
}

// NewMonitor create a monitor instance
// syncer a Sync component
// retry a Retry component
// syncOnce tag a sync once command, the sync once command will execute when call the Start
func NewMonitor(syncer sync.Sync, retry retry.Retry, syncOnce bool, enableTLS bool, certFile string, keyFile string, users []*auth.User) (Monitor, error) {
	src := syncer.Source()
	if src.IsDisk() {
		return NewFsNotifyMonitor(syncer, retry, syncOnce)
	} else if src.Is(core.RemoteDisk) && src.Server() {
		return NewRemoteServerMonitor(syncer, retry, syncOnce)
	} else if src.Is(core.RemoteDisk) && !src.Server() {
		return NewRemoteClientMonitor(syncer, retry, syncOnce, src.Host(), src.Port(), enableTLS, users)
	}
	return nil, fmt.Errorf("file system unsupported ! src=>%s", src.Type().String())

}
