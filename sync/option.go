package sync

import (
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/conf"
	"github.com/no-src/gofs/core"
)

// Option the sync component option
type Option struct {
	Source                core.VFS
	Dest                  core.VFS
	EnableTLS             bool
	TLSCertFile           string
	TLSKeyFile            string
	TLSInsecureSkipVerify bool
	EnableLogicallyDelete bool
	ChunkSize             int64
	CheckpointCount       int
	ForceChecksum         bool
	Users                 []*auth.User
}

// NewSyncOption create an instance of the Option, store all the sync component options
func NewSyncOption(config conf.Config, users []*auth.User) Option {
	opt := Option{
		Source:                config.Source,
		Dest:                  config.Dest,
		EnableTLS:             config.EnableTLS,
		TLSCertFile:           config.TLSCertFile,
		TLSKeyFile:            config.TLSKeyFile,
		TLSInsecureSkipVerify: config.TLSInsecureSkipVerify,
		EnableLogicallyDelete: config.EnableLogicallyDelete,
		ChunkSize:             config.ChunkSize,
		CheckpointCount:       config.CheckpointCount,
		ForceChecksum:         config.ForceChecksum,
		Users:                 users,
	}
	return opt
}
