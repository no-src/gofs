package server

import (
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/log"
)

type Option struct {
	Src            core.VFS
	Dest           core.VFS
	Addr           string
	Init           retry.WaitDone
	EnableTLS      bool
	CertFile       string
	KeyFile        string
	Users          []*auth.User
	EnableCompress bool
	Logger         log.Logger
	EnablePprof    bool
	PprofPrivate   bool
}

func NewServerOption(src core.VFS, dest core.VFS, addr string, init retry.WaitDone, enableTLS bool, certFile string, keyFile string, users []*auth.User, enableCompress bool, logger log.Logger, enablePprof bool, pprofPrivate bool) Option {
	opt := Option{
		Src:            src,
		Dest:           dest,
		Addr:           addr,
		Init:           init,
		EnableTLS:      enableTLS,
		CertFile:       certFile,
		KeyFile:        keyFile,
		Users:          users,
		EnableCompress: enableCompress,
		Logger:         logger,
		EnablePprof:    enablePprof,
		PprofPrivate:   pprofPrivate,
	}
	return opt
}
