package sync

import (
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/logger"
)

type baseSync struct {
	source core.VFS
	dest   core.VFS
	logger *logger.Logger
}

func newBaseSync(source, dest core.VFS, logger *logger.Logger) baseSync {
	return baseSync{
		source: source,
		dest:   dest,
		logger: logger,
	}
}

func (s *baseSync) Source() core.VFS {
	return s.source
}

func (s *baseSync) Dest() core.VFS {
	return s.dest
}

func (s *baseSync) Close() {
}
