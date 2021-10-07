package retry

import (
	"github.com/no-src/log"
	"time"
)

type defaultRetry struct {
	retryCount int
	retryWait  time.Duration
	retryAsync bool
}

func NewRetry(retryCount int, retryWait time.Duration, retryAsync bool) Retry {
	r := &defaultRetry{
		retryCount: retryCount,
		retryWait:  retryWait,
		retryAsync: retryAsync,
	}
	return r
}

// Do if execute failed, retry retryCount times, per wait Duration Sleep
func (r *defaultRetry) Do(f func() error, desc string) {
	if f == nil || f() == nil || r.retryCount <= 0 {
		return
	}
	log.Warn("execute failed, wait to retry [%s]", desc)
	if r.retryAsync {
		go r.retry(f, desc)
	} else {
		r.retry(f, desc)
	}
}

func (r *defaultRetry) retry(f func() error, desc string) {
	for i := 0; i < r.retryCount; i++ {
		err := f()
		if err == nil {
			if i > 0 {
				log.Debug("retry [%d] success [%s] ", i+1, desc)
			}
			break
		} else {
			log.Debug("retry [%d] after %s [%s]", i+1, r.retryWait.String(), desc)
			if i == r.retryCount-1 {
				log.Error(err, "retry [%d] times, and aborted [%s]", r.retryCount, desc)
			} else {
				time.Sleep(r.retryWait)
			}
		}
	}
}
