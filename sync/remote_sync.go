package sync

import (
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/core"
)

// NewRemoteSync auto create an instance of remoteServerSync or remoteClientSync according to src and dest
func NewRemoteSync(src, dest core.VFS, enableTLS bool, certFile string, keyFile string, users []*auth.User, enableLogicallyDelete bool) (Sync, error) {
	if src.Server() {
		return NewRemoteServerSync(src, dest, enableTLS, certFile, keyFile, users, enableLogicallyDelete)
	}
	return NewRemoteClientSync(src, dest, users, enableLogicallyDelete)
}
