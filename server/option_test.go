package server

import (
	"testing"
	"time"

	"github.com/no-src/gofs/conf"
	"github.com/no-src/gofs/report"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/gofs/wait"
)

func TestNewServerOption(t *testing.T) {
	retryWait := time.Second
	opt := NewServerOption(conf.Config{}, wait.NewWaitDone(), nil, nil, retry.New(1, retryWait, false, nil), report.NewReporter())
	if opt.Users != nil || opt.Logger != nil || opt.Retry.WaitTime() != retryWait {
		t.Errorf("NewServerOption() error, option => %v", opt)
	}
}
