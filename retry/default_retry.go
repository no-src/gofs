package retry

import (
	"context"
	"time"

	"github.com/no-src/gofs/logger"
	"github.com/no-src/gofs/wait"
)

type defaultRetry struct {
	count  int
	wait   time.Duration
	async  bool
	logger *logger.Logger
}

// New get a default retry instance
// count the retry execute count
// wait execute once per wait interval
// async is async or sync to execute retry
func New(count int, wait time.Duration, async bool, logger *logger.Logger) Retry {
	r := &defaultRetry{
		count:  count,
		wait:   wait,
		async:  async,
		logger: logger,
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
			r.logger.Warn("retry do recover from => [%s] error => %v", desc, e)
			wd.Done()
			w = wd
		}
	}()

	if f == nil || f() == nil || r.count <= 0 {
		wd.Done()
		return wd
	}
	r.logger.Warn("execute failed, wait to retry [%s] %d times, execute once per %s", desc, r.count, r.wait)
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
			r.logger.Debug("retry [%d] [%s] done => %s", i+1, desc, ctx.Err())
			return
		default:

		}
		err := f()
		if err == nil {
			if i > 0 {
				r.logger.Debug("retry [%d] success [%s] ", i+1, desc)
			}
			break
		} else {
			r.logger.Debug("retry [%d] after %s [%s]", i+1, r.wait.String(), desc)
			if i == r.count-1 {
				r.logger.Error(err, "retry [%d] times, and aborted [%s]", r.count, desc)
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
