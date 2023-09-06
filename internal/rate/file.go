package rate

import (
	"io"
	"net/http"

	"github.com/no-src/gofs/logger"
)

type file struct {
	http.File

	r io.Reader
}

// NewFile create a limit http.File that wrap the real http.File.
func NewFile(f http.File, bytesPerSecond int64, logger *logger.Logger) http.File {
	return &file{
		File: f,
		r:    NewReader(f, bytesPerSecond, logger),
	}
}

func (f *file) Read(p []byte) (n int, err error) {
	return f.r.Read(p)
}
