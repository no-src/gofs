package rate

import (
	"io"

	"github.com/no-src/gofs/logger"
)

type reader struct {
	rate *rateReader
}

// NewReader create a limit io.Reader that wrap the real io.Reader.
// The bytesPerSecond must be greater than defaultBufSize of io.Reader.
func NewReader(r io.Reader, bytesPerSecond int64, logger *logger.Logger) io.Reader {
	if bytesPerSecond <= 0 {
		return r
	}
	return &reader{
		rate: newRateReader(r, nil, bytesPerSecond, logger),
	}
}

func (r *reader) Read(p []byte) (n int, err error) {
	return r.rate.Read(p)
}
