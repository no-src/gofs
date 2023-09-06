package rate

import (
	"context"
	"io"
	"sync"

	"github.com/no-src/gofs/logger"
	"golang.org/x/time/rate"
)

type rateReader struct {
	r              io.Reader
	ra             io.ReaderAt
	l              *rate.Limiter
	bytesPerSecond int64
	c              int64
	once           sync.Once
	logger         *logger.Logger
}

func newRateReader(r io.Reader, ra io.ReaderAt, bytesPerSecond int64, logger *logger.Logger) *rateReader {
	return &rateReader{
		r:              r,
		ra:             ra,
		l:              rate.NewLimiter(1, 1),
		bytesPerSecond: bytesPerSecond,
		logger:         logger,
	}
}

func (r *rateReader) Read(p []byte) (n int, err error) {
	return r.call(func() (n int, err error) {
		return r.r.Read(p)
	})
}

func (r *rateReader) ReadAt(p []byte, off int64) (n int, err error) {
	return r.call(func() (n int, err error) {
		return r.ra.ReadAt(p, off)
	})
}

func (r *rateReader) call(f func() (n int, err error)) (n int, err error) {
	r.once.Do(func() {
		r.l.Reserve()
	})
	n, err = f()
	if err != nil {
		return n, err
	}
	r.c += int64(n)
	if r.c >= r.bytesPerSecond {
		r.logger.ErrorIf(r.l.Wait(context.Background()), "limiter wait error")
		r.c = 0
	}
	return n, err
}
