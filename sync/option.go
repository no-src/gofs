package sync

import (
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/conf"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/encrypt"
	"github.com/no-src/gofs/retry"
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
	ChecksumAlgorithm     string
	Progress              bool
	MaxTranRate           int64
	SSHKey                string
	Users                 []*auth.User
	Retry                 retry.Retry
	EncOpt                encrypt.Option
}

// NewSyncOption create an instance of the Option, store all the sync component options
func NewSyncOption(config conf.Config, users []*auth.User, r retry.Retry) Option {
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
		ChecksumAlgorithm:     config.ChecksumAlgorithm,
		Progress:              config.Progress,
		MaxTranRate:           config.MaxTranRate,
		SSHKey:                config.SSHKey,
		Users:                 users,
		Retry:                 r,
		EncOpt:                encrypt.NewOption(config),
	}
	return opt
}
