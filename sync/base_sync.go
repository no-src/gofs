package sync

import (
	"github.com/no-src/gofs/core"
)

type baseSync struct {
	source core.VFS
	dest   core.VFS
}

func newBaseSync(source, dest core.VFS) baseSync {
	return baseSync{
		source: source,
		dest:   dest,
	}
}

func (s *baseSync) Source() core.VFS {
	return s.source
}

func (s *baseSync) Dest() core.VFS {
	return s.dest
}
