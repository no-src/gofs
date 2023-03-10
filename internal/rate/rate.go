package rate

import (
	"context"
	"io"
	"sync"

	"github.com/no-src/log"
	"golang.org/x/time/rate"
)

type reader struct {
	r              io.Reader
	ra             io.ReaderAt
	l              *rate.Limiter
	bytesPerSecond int64
	c              int64
	once           sync.Once
}

// NewReader create a limit io.Reader that wrap the real io.Reader.
// The bytesPerSecond must be more than defaultBufSize of io.Reader.
func NewReader(r io.Reader, bytesPerSecond int64) io.Reader {
	if bytesPerSecond <= 0 {
		return r
	}
	return newReader(r, nil, bytesPerSecond)
}

// NewReaderAt create a limit io.ReaderAt that wrap the real io.ReaderAt.
// The bytesPerSecond must be more than defaultBufSize of io.ReaderAt.
func NewReaderAt(ra io.ReaderAt, bytesPerSecond int64) io.ReaderAt {
	if bytesPerSecond <= 0 {
		return ra
	}
	return newReader(nil, ra, bytesPerSecond)
}

func newReader(r io.Reader, ra io.ReaderAt, bytesPerSecond int64) *reader {
	return &reader{
		r:              r,
		ra:             ra,
		l:              rate.NewLimiter(1, 1),
		bytesPerSecond: bytesPerSecond,
	}
}

func (r *reader) Read(p []byte) (n int, err error) {
	return r.call(func() (n int, err error) {
		return r.r.Read(p)
	})
}

func (r *reader) ReadAt(p []byte, off int64) (n int, err error) {
	return r.call(func() (n int, err error) {
		return r.ra.ReadAt(p, off)
	})
}

func (r *reader) call(f func() (n int, err error)) (n int, err error) {
	r.once.Do(func() {
		r.l.Reserve()
	})
	n, err = f()
	if err != nil {
		return n, err
	}
	r.c += int64(n)
	if r.c >= r.bytesPerSecond {
		log.ErrorIf(r.l.Wait(context.Background()), "limiter wait error")
		r.c = 0
	}
	return n, err
}
