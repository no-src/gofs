package sync

import (
	"fmt"
	"github.com/no-src/gofs/core"
)

type Sync interface {
	Create(path string) error
	Write(path string) error
	Remove(path string) error
	Rename(path string) error
	Chmod(path string) error
	IsDir(path string) (bool, error)
	SyncOnce(path string) error
	Source() core.VFS
	Target() core.VFS
}

func NewSync(src core.VFS, target core.VFS, bufSize int) (Sync, error) {
	if src.IsDisk() && target.IsDisk() {
		return NewDiskSync(src, target, bufSize)
	} else if src.Is(core.RemoteDisk) {
		return NewRemoteSync(src, target, bufSize)
	}
	return nil, fmt.Errorf("file system unsupported ! src=>%s target=>%s", src.Type().String(), target.Type().String())
}
