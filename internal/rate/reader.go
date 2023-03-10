package rate

import "io"

type reader struct {
	rate *rateReader
}

// NewReader create a limit io.Reader that wrap the real io.Reader.
// The bytesPerSecond must be greater than defaultBufSize of io.Reader.
func NewReader(r io.Reader, bytesPerSecond int64) io.Reader {
	if bytesPerSecond <= 0 {
		return r
	}
	return &reader{
		rate: newRateReader(r, nil, bytesPerSecond),
	}
}

func (r *reader) Read(p []byte) (n int, err error) {
	return r.rate.Read(p)
}
