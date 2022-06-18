package sync

import (
	"fmt"

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
func NewSync(opt Option) (Sync, error) {
	if opt.Source.IsDisk() && opt.Dest.IsDisk() {
		return NewDiskSync(opt.Source, opt.Dest, opt.EnableLogicallyDelete, opt.ChunkSize, opt.CheckpointCount, opt.ForceChecksum)
	} else if opt.Source.Is(core.RemoteDisk) {
		return NewRemoteSync(opt.Source, opt.Dest, opt.EnableTLS, opt.TLSCertFile, opt.TLSKeyFile, opt.Users, opt.EnableLogicallyDelete, opt.ChunkSize, opt.CheckpointCount, opt.ForceChecksum)
	} else if opt.Dest.Is(core.RemoteDisk) {
		return NewPushClientSync(opt.Source, opt.Dest, opt.EnableTLS, opt.TLSCertFile, opt.TLSInsecureSkipVerify, opt.Users, opt.EnableLogicallyDelete, opt.ChunkSize, opt.CheckpointCount, opt.ForceChecksum)
	} else if opt.Source.IsDisk() && opt.Dest.Is(core.SFTP) {
		return NewSftpClientSync(opt.Source, opt.Dest, opt.Users, opt.EnableLogicallyDelete, opt.ChunkSize, opt.CheckpointCount, opt.ForceChecksum, opt.Retry)
	}
	return nil, fmt.Errorf("file system unsupported ! source=>%s dest=>%s", opt.Source.Type().String(), opt.Dest.Type().String())
}
