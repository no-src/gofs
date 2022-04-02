package sync

import (
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/core"
)

// NewRemoteSync auto create an instance of remoteServerSync or remoteClientSync according to source and dest
func NewRemoteSync(source, dest core.VFS, enableTLS bool, certFile string, keyFile string, users []*auth.User, enableLogicallyDelete bool, chunkSize int64, checkpointCount int) (Sync, error) {
	if source.Server() {
		return NewRemoteServerSync(source, dest, enableTLS, certFile, keyFile, users, enableLogicallyDelete, chunkSize, checkpointCount)
	}
	return NewRemoteClientSync(source, dest, users, enableLogicallyDelete, chunkSize)
}
