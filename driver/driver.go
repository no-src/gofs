package driver

import (
	"io"
	"time"
)

// Driver a data source client driver
type Driver interface {
	// Connect connects the server
	Connect() error
	// MkdirAll creates a directory named path
	MkdirAll(path string) error
	// Create creates the named file
	Create(path string) (rwc io.ReadWriteCloser, err error)
	// Remove removes the specified file or directory
	Remove(path string) error
	// Rename renames a file
	Rename(oldPath, newPath string) error
	// Chtimes changes the access and modification times of the named file
	Chtimes(path string, aTime time.Time, mTime time.Time) error
}
