package signal

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/no-src/log"
)

var (
	defaultSendSignalTimeout = time.Second * 3
	errSendSignalTimeout     = errors.New("send signal timeout")
)

// NotifySignal sends a signal with timeout
type NotifySignal func(s os.Signal, timeout ...time.Duration) error

// Notify receive signal and try to shut down
func Notify(shutdown func() error) NotifySignal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM)
	go func() {
		for {
			s := <-c
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM:
				log.Debug("received a signal [%s], waiting to exit", s.String())
				err := shutdown()
				if err != nil {
					log.Error(err, "shutdown error")
				} else {
					signal.Stop(c)
					log.Debug("shutdown success")
					return
				}
			default:
				log.Debug("received a signal [%s], ignore it", s.String())
			}
		}
	}()
	return func(s os.Signal, timeout ...time.Duration) error {
		t := defaultSendSignalTimeout
		if len(timeout) > 0 {
			t = timeout[0]
		}
		select {
		case c <- s:
			log.Debug("[success] send a signal [%s] by user", s.String())
			return nil
		case <-time.After(t):
			log.Warn("[timeout] send a signal [%s] by user", s.String())
			return fmt.Errorf("%w => %s", errSendSignalTimeout, s.String())
		}
	}
}
