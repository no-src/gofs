package retry

import (
	"github.com/no-src/log"
	"time"
)

type defaultRetry struct {
	retryCount int
	retryWait  time.Duration
}

func NewRetry(retryCount int, retryWait time.Duration) Retry {
	r := &defaultRetry{
		retryCount: retryCount,
		retryWait:  retryWait,
	}
	return r
}

// Do if execute failed, retry retryCount times,per wait Duration Sleep
func (r *defaultRetry) Do(f func() error, desc string) {
	if f == nil || f() == nil || r.retryCount <= 0 {
		return
	}
	log.Warn("execute failed, wait to retry [%s]", desc)
	go func() {
		for i := 0; i < r.retryCount; i++ {
			err := f()
			if err == nil {
				if i > 0 {
					log.Log("retry [%d] success [%s] ", i+1, desc)
				}
				break
			} else {
				log.Log("retry [%d] after %s [%s]", i+1, r.retryWait.String(), desc)
				if i == r.retryCount-1 {
					log.Error(err, "retry [%d] times, and aborted [%s]", r.retryCount, desc)
				} else {
					time.Sleep(r.retryWait)
				}
			}
		}
	}()
}
