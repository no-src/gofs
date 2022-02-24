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
	Source                core.VFS
	Dest                  core.VFS
	Addr                  string
	Init                  wait.WaitDone
	EnableTLS             bool
	CertFile              string
	KeyFile               string
	Users                 []*auth.User
	EnableCompress        bool
	Logger                log.Logger
	EnablePprof           bool
	PprofPrivate          bool
	EnableLogicallyDelete bool
	EnablePushServer      bool
}

// NewServerOption create an instance of the Option, store all the web server options
func NewServerOption(config conf.Config, init wait.WaitDone, users []*auth.User, logger log.Logger) Option {
	opt := Option{
		Source:                config.Source,
		Dest:                  config.Dest,
		Addr:                  config.FileServerAddr,
		Init:                  init,
		EnableTLS:             config.EnableTLS,
		CertFile:              config.TLSCertFile,
		KeyFile:               config.TLSKeyFile,
		Users:                 users,
		EnableCompress:        config.EnableFileServerCompress,
		Logger:                logger,
		EnablePprof:           config.EnablePProf,
		PprofPrivate:          config.PProfPrivate,
		EnableLogicallyDelete: config.EnableLogicallyDelete,
		EnablePushServer:      config.EnablePushServer,
	}
	return opt
}
