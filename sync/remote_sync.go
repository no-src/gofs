package sync

import (
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/server/middleware/auth"
)

// NewRemoteSync auto create an instance of remoteServerSync or remoteClientSync according to src and target
func NewRemoteSync(src, target core.VFS, bufSize int, users []*auth.User) (Sync, error) {
	if src.Server() {
		return NewRemoteServerSync(src, target, bufSize)
	}
	return NewRemoteClientSync(src, target, bufSize, users)
}
