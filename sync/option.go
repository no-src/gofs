package sync

import (
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/conf"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/encrypt"
	"github.com/no-src/gofs/ignore"
	"github.com/no-src/gofs/report"
	"github.com/no-src/gofs/retry"
)

// Option the sync component option
type Option struct {
	Source                core.VFS
	Dest                  core.VFS
	EnableHTTP3           bool
	FileServerAddr        string
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
	DryRun                bool
	CopyLink              bool
	CopyUnsafeLink        bool
	SSHKey                string
	TokenSecret           string
	Users                 []*auth.User
	Retry                 retry.Retry
	EncOpt                encrypt.Option
	PathIgnore            ignore.PathIgnore
	Reporter              report.Reporter
	TaskConf              string
}

// NewSyncOption create an instance of the Option, store all the sync component options
func NewSyncOption(config conf.Config, users []*auth.User, r retry.Retry, pi ignore.PathIgnore, reporter report.Reporter) Option {
	opt := Option{
		Source:                config.Source,
		Dest:                  config.Dest,
		EnableHTTP3:           config.EnableHTTP3,
		FileServerAddr:        config.FileServerAddr,
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
		DryRun:                config.DryRun,
		CopyLink:              config.CopyLink,
		CopyUnsafeLink:        config.CopyUnsafeLink,
		SSHKey:                config.SSHKey,
		TokenSecret:           config.TokenSecret,
		Users:                 users,
		Retry:                 r,
		EncOpt:                encrypt.NewOption(config),
		PathIgnore:            pi,
		Reporter:              reporter,
		TaskConf:              config.TaskConf,
	}
	return opt
}
