package driver

import (
	"io/fs"
	"net/http"
	"os"
	"time"
)

// Driver a data source client driver
type Driver interface {
	// DriverName return driver name
	DriverName() string
	// Connect connects the server
	Connect() error
	// MkdirAll creates a directory named path
	MkdirAll(path string) error
	// Create creates the named file
	Create(path string) (err error)
	// Remove removes the specified file or directory
	Remove(path string) error
	// Rename renames a file
	Rename(oldPath, newPath string) error
	// Chtimes changes the access and modification times of the named file
	Chtimes(path string, aTime time.Time, mTime time.Time) error
	// WalkDir walks the file tree rooted at root, calling fn for each file or directory in the tree, including root
	WalkDir(root string, fn fs.WalkDirFunc) error
	// Open opens the named file for reading
	Open(path string) (f http.File, err error)
	// Stat returns the os.FileInfo describing the named file
	Stat(path string) (fi os.FileInfo, err error)
	// GetFileTime get the creation time, last access time, last modify time of the path
	GetFileTime(path string) (cTime time.Time, aTime time.Time, mTime time.Time, err error)
	// Write write src file to dest file
	Write(src string, dest string) error
}
