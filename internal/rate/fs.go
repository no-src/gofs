package rate

import (
	"net/http"
)

type dir struct {
	fs             http.FileSystem
	bytesPerSecond int64
}

// NewHTTPDir create a limit http.FileSystem that wrap the real http.Dir.
func NewHTTPDir(path string, bytesPerSecond int64) http.FileSystem {
	return NewDir(http.Dir(path), bytesPerSecond)
}

// NewDir create a limit http.FileSystem that wrap the real http.FileSystem.
func NewDir(fs http.FileSystem, bytesPerSecond int64) http.FileSystem {
	if bytesPerSecond <= 0 {
		return fs
	}
	return &dir{
		fs:             fs,
		bytesPerSecond: bytesPerSecond,
	}
}

func (d *dir) Open(name string) (http.File, error) {
	f, err := d.fs.Open(name)
	if err != nil {
		return f, err
	}
	return NewFile(f, d.bytesPerSecond), err
}
