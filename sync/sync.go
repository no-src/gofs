package sync

import (
	"errors"
	"fmt"

	"github.com/no-src/gofs/core"
)

var (
	errUserIsRequired         = errors.New("user account is required")
	errInvalidChunkSize       = errors.New("chunk size must greater than zero")
	errSourceNotFound         = errors.New("source is not found")
	errDestNotFound           = errors.New("dest is not found")
	errFileServerUnauthorized = errors.New("file server is unauthorized")
	errAuthResponse           = errors.New("receive auth command response error")
	errInfoResponse           = errors.New("receive info command response error")
	errAuthTimeout            = errors.New("auth timeout")
	errInfoTimeout            = errors.New("info timeout")
	errFileSystemUnsupported  = errors.New("file system unsupported")
	errPushServerUnsupported  = errors.New("the push server is unsupported")
	errSendToPushServer       = errors.New("send a request to the push server error")
	errCallQueryAPI           = errors.New("call the query api error")
	errNilRemoteSyncServer    = errors.New("remote sync server is nil")
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
	// the fields of option
	source := opt.Source
	dest := opt.Dest

	if source.IsDisk() && dest.IsDisk() {
		return NewDiskSync(opt)
	} else if source.Is(core.RemoteDisk) {
		return NewRemoteSync(opt)
	} else if dest.Is(core.RemoteDisk) {
		return NewPushClientSync(opt)
	} else if source.IsDisk() && dest.Is(core.SFTP) {
		return NewSftpPushClientSync(opt)
	} else if source.Is(core.SFTP) && dest.IsDisk() {
		return NewSftpPullClientSync(opt)
	} else if source.IsDisk() && dest.Is(core.MinIO) {
		return NewMinIOPushClientSync(opt)
	} else if source.Is(core.MinIO) && dest.IsDisk() {
		return NewMinIOPullClientSync(opt)
	}
	return nil, fmt.Errorf("%w source=>%s dest=>%s", errFileSystemUnsupported, source.Type().String(), dest.Type().String())
}
