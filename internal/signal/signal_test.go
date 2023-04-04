package signal

import (
	"errors"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestNotify(t *testing.T) {
	st := Notify(func() error {
		return nil
	})

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
			st(tc.signal, time.Second)
		})
	}
}

func TestNotify_ShutdownError(t *testing.T) {
	st := Notify(func() error {
		return errors.New("shutdown error mock")
	})

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
			st(tc.signal)
		})
	}
}

func TestNotify_IgnoreSignal(t *testing.T) {
	st := Notify(func() error {
		return nil
	})

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
			st(tc.signal)
		})
	}
}
