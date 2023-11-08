package rate

import (
	"net/http"

	"github.com/no-src/gofs/logger"
)

type dir struct {
	fs             http.FileSystem
	bytesPerSecond int64
	logger         *logger.Logger
}

// NewHTTPDir create a limit http.FileSystem that wrap the real http.Dir.
func NewHTTPDir(path string, bytesPerSecond int64, logger *logger.Logger) http.FileSystem {
	return NewDir(http.Dir(path), bytesPerSecond, logger)
}

// NewDir create a limit http.FileSystem that wrap the real http.FileSystem.
func NewDir(fs http.FileSystem, bytesPerSecond int64, logger *logger.Logger) http.FileSystem {
	if bytesPerSecond <= 0 {
		return fs
	}
	return &dir{
		fs:             fs,
		bytesPerSecond: bytesPerSecond,
		logger:         logger,
	}
}

func (d *dir) Open(name string) (http.File, error) {
	f, err := d.fs.Open(name)
	if err != nil {
		return f, err
	}
	return NewFile(f, d.bytesPerSecond, d.logger), err
}
