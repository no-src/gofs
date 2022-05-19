package server

import (
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/conf"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/wait"
	"github.com/no-src/log"
)

// Option the web server option
type Option struct {
	Source                   core.VFS
	Dest                     core.VFS
	Addr                     string
	EnableTLS                bool
	TLSCertFile              string
	TLSKeyFile               string
	EnableFileServerCompress bool
	EnableManage             bool
	ManagePrivate            bool
	EnableLogicallyDelete    bool
	EnablePushServer         bool
	EnableReport             bool
	SessionMode              int
	SessionConnection        string
	ChunkSize                int64
	CheckpointCount          int
	Init                     wait.WaitDone
	Users                    []*auth.User
	Logger                   log.Logger
}

// NewServerOption create an instance of the Option, store all the web server options
func NewServerOption(config conf.Config, init wait.WaitDone, users []*auth.User, logger log.Logger) Option {
	opt := Option{
		Source:                   config.Source,
		Dest:                     config.Dest,
		Addr:                     config.FileServerAddr,
		EnableTLS:                config.EnableTLS,
		TLSCertFile:              config.TLSCertFile,
		TLSKeyFile:               config.TLSKeyFile,
		EnableFileServerCompress: config.EnableFileServerCompress,
		EnableManage:             config.EnableManage,
		ManagePrivate:            config.ManagePrivate,
		EnableLogicallyDelete:    config.EnableLogicallyDelete,
		EnablePushServer:         config.EnablePushServer,
		EnableReport:             config.EnableReport,
		SessionMode:              config.SessionMode,
		SessionConnection:        config.SessionConnection,
		ChunkSize:                config.ChunkSize,
		CheckpointCount:          config.CheckpointCount,
		Init:                     init,
		Users:                    users,
		Logger:                   logger,
	}
	return opt
}
