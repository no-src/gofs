package server

import (
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/conf"
	"github.com/no-src/gofs/report"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/gofs/wait"
	"github.com/no-src/log"
)

// Option the web server option
type Option struct {
	conf.Config

	Init     wait.Done
	Users    []*auth.User
	Logger   log.Logger
	Retry    retry.Retry
	Reporter *report.Reporter
}

// NewServerOption create an instance of the Option, store all the web server options
func NewServerOption(c conf.Config, init wait.Done, users []*auth.User, logger log.Logger, r retry.Retry, reporter *report.Reporter) Option {
	opt := Option{
		Config:   c,
		Init:     init,
		Users:    users,
		Logger:   logger,
		Retry:    r,
		Reporter: reporter,
	}
	return opt
}
