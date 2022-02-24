package signal

import (
	"github.com/no-src/log"
	"os"
	"os/signal"
	"syscall"
)

// Notify receive signal and try to shut down
func Notify(shutdown func() error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGKILL, syscall.SIGTERM)
	for {
		s := <-c
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGKILL, syscall.SIGTERM:
			log.Debug("received a signal [%s], waiting to exit", s.String())
			err := shutdown()
			if err != nil {
				log.Error(err, "shutdown error")
				break
			} else {
				signal.Stop(c)
				close(c)
				return
			}
		default:
			log.Debug("received a signal [%s], ignore it", s.String())
			break
		}
	}
}
