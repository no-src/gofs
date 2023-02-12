package server

import (
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/conf"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/gofs/wait"
	"github.com/no-src/log"
)

// Option the web server option
type Option struct {
	Source                   core.VFS
	Dest                     core.VFS
	Addr                     string
	EnableHTTP3              bool
	EnableTLS                bool
	TLSCertFile              string
	TLSKeyFile               string
	EnableFileServerCompress bool
	EnableManage             bool
	ManagePrivate            bool
	EnableLogicallyDelete    bool
	EnablePushServer         bool
	EnableReport             bool
	SessionConnection        string
	ChunkSize                int64
	CheckpointCount          int
	Init                     wait.Done
	Users                    []*auth.User
	SSHKey                   string
	Logger                   log.Logger
	Retry                    retry.Retry
}

// NewServerOption create an instance of the Option, store all the web server options
func NewServerOption(config conf.Config, init wait.Done, users []*auth.User, logger log.Logger, r retry.Retry) Option {
	opt := Option{
		Source:                   config.Source,
		Dest:                     config.Dest,
		Addr:                     config.FileServerAddr,
		EnableHTTP3:              config.EnableHTTP3,
		EnableTLS:                config.EnableTLS,
		TLSCertFile:              config.TLSCertFile,
		TLSKeyFile:               config.TLSKeyFile,
		EnableFileServerCompress: config.EnableFileServerCompress,
		EnableManage:             config.EnableManage,
		ManagePrivate:            config.ManagePrivate,
		EnableLogicallyDelete:    config.EnableLogicallyDelete,
		EnablePushServer:         config.EnablePushServer,
		EnableReport:             config.EnableReport,
		SessionConnection:        config.SessionConnection,
		ChunkSize:                config.ChunkSize,
		CheckpointCount:          config.CheckpointCount,
		Init:                     init,
		Users:                    users,
		SSHKey:                   config.SSHKey,
		Logger:                   logger,
		Retry:                    r,
	}
	return opt
}
