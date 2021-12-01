package sync

import (
	"fmt"
	"github.com/no-src/gofs/contract"
	"github.com/no-src/gofs/core"
)

// Sync a file sync interface
type Sync interface {
	// Create create the path
	Create(path string) error
	// Write write the data to path
	Write(path string) error
	// Remove remove the path
	Remove(path string) error
	// Rename rename the path
	Rename(path string) error
	// Chmod change the mode of path
	Chmod(path string) error
	// IsDir is a dir the path
	IsDir(path string) (bool, error)
	// SyncOnce sync the path to target once
	SyncOnce(path string) error
	// Source the source file system
	Source() core.VFS
	// Target the target file system
	Target() core.VFS
}

// NewSync auto create an instance of the expected sync according to src and target
func NewSync(src core.VFS, target core.VFS, bufSize int, users []*contract.User) (Sync, error) {
	if src.IsDisk() && target.IsDisk() {
		return NewDiskSync(src, target, bufSize)
	} else if src.Is(core.RemoteDisk) {
		return NewRemoteSync(src, target, bufSize, users)
	}
	return nil, fmt.Errorf("file system unsupported ! src=>%s target=>%s", src.Type().String(), target.Type().String())
}
