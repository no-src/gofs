package retry

import (
	"context"
	"time"

	"github.com/no-src/gofs/wait"
	"github.com/no-src/log"
)

type defaultRetry struct {
	count int
	wait  time.Duration
	async bool
}

// New get a default retry instance
// count the retry execute count
// wait execute once per wait interval
// async is async or sync to execute retry
func New(count int, wait time.Duration, async bool) Retry {
	r := &defaultRetry{
		count: count,
		wait:  wait,
		async: async,
	}
	return r
}

func (r *defaultRetry) Do(f func() error, desc string) wait.Wait {
	return r.DoWithContext(context.Background(), f, desc)
}

func (r *defaultRetry) DoWithContext(ctx context.Context, f func() error, desc string) (w wait.Wait) {
	wd := wait.NewWaitDone()
	defer func() {
		e := recover()
		if e != nil {
			log.Warn("retry do recover from => [%s] error => %v", desc, e)
			wd.Done()
			w = wd
		}
	}()

	if f == nil || f() == nil || r.count <= 0 {
		wd.Done()
		return wd
	}
	log.Warn("execute failed, wait to retry [%s] %d times, execute once per %s", desc, r.count, r.wait)
	if r.async {
		go r.retry(ctx, wd, f, desc)
	} else {
		r.retry(ctx, wd, f, desc)
	}
	return wd
}

func (r *defaultRetry) retry(ctx context.Context, wd wait.Done, f func() error, desc string) {
	defer func() {
		wd.Done()
	}()
	for i := 0; i < r.count; i++ {
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
			log.Debug("retry [%d] after %s [%s]", i+1, r.wait.String(), desc)
			if i == r.count-1 {
				log.Error(err, "retry [%d] times, and aborted [%s]", r.count, desc)
			} else {
				time.Sleep(r.wait)
			}
		}
	}
}

func (r *defaultRetry) Count() int {
	return r.count
}

func (r *defaultRetry) WaitTime() time.Duration {
	return r.wait
}
