package signal

import (
	"errors"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/no-src/gofs/internal/logger"
)

func TestNotify(t *testing.T) {
	logger := logger.NewTestLogger()
	defer logger.Close()

	ns, ss := Notify(func() error {
		return nil
	}, logger)

	testCases := []struct {
		name   string
		signal os.Signal
	}{
		{"SIGHUP", syscall.SIGHUP},
		{"SIGINT", syscall.SIGINT},
		{"SIGQUIT", syscall.SIGQUIT},
		{"SIGABRT", syscall.SIGABRT},
		{"SIGTERM", syscall.SIGTERM},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ns(tc.signal, time.Second)
			ss()
		})
	}
}

func TestNotify_ShutdownError(t *testing.T) {
	logger := logger.NewTestLogger()
	defer logger.Close()

	ns, ss := Notify(func() error {
		return errors.New("shutdown error mock")
	}, logger)

	testCases := []struct {
		name   string
		signal os.Signal
	}{
		{"SIGHUP", syscall.SIGHUP},
		{"SIGINT", syscall.SIGINT},
		{"SIGQUIT", syscall.SIGQUIT},
		{"SIGABRT", syscall.SIGABRT},
		{"SIGTERM", syscall.SIGTERM},

		{"SIGALRM", syscall.SIGALRM},
		{"SIGPIPE", syscall.SIGPIPE},
		{"SIGFPE", syscall.SIGFPE},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ns(tc.signal)
			ss()
		})
	}
}

func TestNotify_IgnoreSignal(t *testing.T) {
	logger := logger.NewTestLogger()
	defer logger.Close()

	ns, ss := Notify(func() error {
		return nil
	}, logger)

	testCases := []struct {
		name   string
		signal os.Signal
	}{
		{"SIGALRM", syscall.SIGALRM},
		{"SIGPIPE", syscall.SIGPIPE},
		{"SIGFPE", syscall.SIGFPE},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ns(tc.signal)
			ss()
		})
	}
}
