package sync

import (
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/core"
)

// NewRemoteSync auto create an instance of remoteServerSync or remoteClientSync according to src and target
func NewRemoteSync(src, target core.VFS, bufSize int, enableTLS bool, certFile string, keyFile string, users []*auth.User) (Sync, error) {
	if src.Server() {
		return NewRemoteServerSync(src, target, bufSize, enableTLS, certFile, keyFile, users)
	}
	return NewRemoteClientSync(src, target, bufSize, users)
}
