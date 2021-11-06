package sync

import "github.com/no-src/gofs/core"

func NewRemoteSync(src, target core.VFS, bufSize int) (Sync, error) {
	if src.Server() {
		return NewRemoteServerSync(src, target, bufSize)
	}
	return NewRemoteClientSync(src, target, bufSize)
}
