package sync

import (
	"github.com/no-src/gofs/fs"
)

type baseSync struct {
	enableLogicallyDelete bool
}

func newBaseSync(enableLogicallyDelete bool) baseSync {
	return baseSync{
		enableLogicallyDelete: enableLogicallyDelete,
	}
}

// logicallyDelete delete the path logically
func (s *baseSync) logicallyDelete(path string) error {
	return fs.LogicallyDelete(path)
}
