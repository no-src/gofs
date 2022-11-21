package monitor

import (
	"io"
	"time"

	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/conf"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/gofs/sync"
)

// Option the monitor component option
type Option struct {
	SyncOnce              bool
	EnableTLS             bool
	TLSCertFile           string
	TLSInsecureSkipVerify bool
	EnableSyncDelay       bool
	SyncDelayEvents       int
	SyncDelayTime         time.Duration
	SyncWorkers           int
	Users                 []*auth.User
	EventWriter           io.Writer
	Syncer                sync.Sync
	Retry                 retry.Retry
}

// NewMonitorOption create an instance of the Option, store all the monitor component options
func NewMonitorOption(config conf.Config, syncer sync.Sync, retry retry.Retry, users []*auth.User, eventWriter io.Writer) Option {
	opt := Option{
		SyncOnce:              config.SyncOnce,
		EnableTLS:             config.EnableTLS,
		TLSCertFile:           config.TLSCertFile,
		TLSInsecureSkipVerify: config.TLSInsecureSkipVerify,
		EnableSyncDelay:       config.EnableSyncDelay,
		SyncDelayEvents:       config.SyncDelayEvents,
		SyncDelayTime:         config.SyncDelayTime.Duration(),
		SyncWorkers:           config.SyncWorkers,
		Syncer:                syncer,
		Retry:                 retry,
		Users:                 users,
		EventWriter:           eventWriter,
	}
	return opt
}
