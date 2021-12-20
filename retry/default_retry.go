package retry

import (
	"context"
	"github.com/no-src/log"
	"time"
)

type defaultRetry struct {
	retryCount int
	retryWait  time.Duration
	retryAsync bool
}

// NewRetry get a default retry instance
// retryCount retry execute count
// retryWait execute once per retryWait interval
// retryAsync async or sync to execute retry
func NewRetry(retryCount int, retryWait time.Duration, retryAsync bool) Retry {
	r := &defaultRetry{
		retryCount: retryCount,
		retryWait:  retryWait,
		retryAsync: retryAsync,
	}
	return r
}

func (r *defaultRetry) Do(f func() error, desc string) Wait {
	return r.DoWithContext(context.Background(), f, desc)
}

func (r *defaultRetry) DoWithContext(ctx context.Context, f func() error, desc string) Wait {
	defer func() {
		e := recover()
		if e != nil {
			log.Warn("retry do recover from => [%s] error => %v", desc, e)
		}
	}()
	wait := NewWaitDone()
	if f == nil || f() == nil || r.retryCount <= 0 {
		wait.Done()
		return wait
	}
	log.Warn("execute failed, wait to retry [%s] %d times, execute once per %s", desc, r.retryCount, r.retryWait)
	if r.retryAsync {
		go r.retry(ctx, wait, f, desc)
	} else {
		r.retry(ctx, wait, f, desc)
	}
	return wait
}

func (r *defaultRetry) retry(ctx context.Context, wait WaitDone, f func() error, desc string) {
	defer func() {
		wait.Done()
	}()
	for i := 0; i < r.retryCount; i++ {
		select {
		case <-ctx.Done():
			log.Debug("retry [%d] [%s] done => %s", i+1, desc, ctx.Err())
			return
		default:

		}
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

func (r *defaultRetry) RetryCount() int {
	return r.retryCount
}

func (r *defaultRetry) RetryWait() time.Duration {
	return r.retryWait
}
