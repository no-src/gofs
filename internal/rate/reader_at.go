package rate

import (
	"io"

	"github.com/no-src/gofs/logger"
)

type readerAt struct {
	rate *rateReader
}

// NewReaderAt create a limit io.ReaderAt that wrap the real io.ReaderAt.
// The bytesPerSecond must be greater than defaultBufSize of io.ReaderAt.
func NewReaderAt(ra io.ReaderAt, bytesPerSecond int64, logger *logger.Logger) io.ReaderAt {
	if bytesPerSecond <= 0 {
		return ra
	}
	return &readerAt{
		rate: newRateReader(nil, ra, bytesPerSecond, logger),
	}
}

func (r *readerAt) ReadAt(p []byte, off int64) (n int, err error) {
	return r.rate.ReadAt(p, off)
}
