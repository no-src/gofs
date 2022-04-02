package sync

import (
	"fmt"
	"github.com/no-src/gofs/auth"
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
	// SyncOnce sync the path to dest once
	SyncOnce(path string) error
	// Source the source file system
	Source() core.VFS
	// Dest the destination file system
	Dest() core.VFS
}

// NewSync auto create an instance of the expected sync according to source and dest
func NewSync(source core.VFS, dest core.VFS, enableTLS bool, certFile string, keyFile string, users []*auth.User, enableLogicallyDelete bool, chunkSize int64, checkpointCount int) (Sync, error) {
	if source.IsDisk() && dest.IsDisk() {
		return NewDiskSync(source, dest, enableLogicallyDelete)
	} else if source.Is(core.RemoteDisk) {
		return NewRemoteSync(source, dest, enableTLS, certFile, keyFile, users, enableLogicallyDelete, chunkSize, checkpointCount)
	} else if dest.Is(core.RemoteDisk) {
		return NewPushClientSync(source, dest, enableTLS, users, enableLogicallyDelete, chunkSize, checkpointCount)
	}
	return nil, fmt.Errorf("file system unsupported ! source=>%s dest=>%s", source.Type().String(), dest.Type().String())
}
