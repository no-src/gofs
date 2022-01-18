package sync

import (
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/core"
)

// NewRemoteSync auto create an instance of remoteServerSync or remoteClientSync according to src and target
func NewRemoteSync(src, target core.VFS, enableTLS bool, certFile string, keyFile string, users []*auth.User, enableLogicallyDelete bool) (Sync, error) {
	if src.Server() {
		return NewRemoteServerSync(src, target, enableTLS, certFile, keyFile, users, enableLogicallyDelete)
	}
	return NewRemoteClientSync(src, target, users, enableLogicallyDelete)
}
