package sync

import (
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/fs"
)

type baseSync struct {
	source                core.VFS
	dest                  core.VFS
	enableLogicallyDelete bool
}

func newBaseSync(source, dest core.VFS, enableLogicallyDelete bool) baseSync {
	return baseSync{
		source:                source,
		dest:                  dest,
		enableLogicallyDelete: enableLogicallyDelete,
	}
}

// logicallyDelete delete the path logically
func (s *baseSync) logicallyDelete(path string) error {
	return fs.LogicallyDelete(path)
}
