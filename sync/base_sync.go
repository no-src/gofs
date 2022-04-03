package sync

import (
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/fs"
	"github.com/no-src/gofs/util/hashutil"
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

func (s *baseSync) Source() core.VFS {
	return s.source
}

func (s *baseSync) Dest() core.VFS {
	return s.dest
}

func (s *baseSync) compareHashValues(dstPath string, sourceSize int64, sourceHash string, chunkSize int64, hvs hashutil.HashValues) (equal bool, hv *hashutil.HashValue) {
	if sourceSize > 0 {
		// calculate the entire file hash value
		if len(hvs) == 0 || hvs.Last().Offset < sourceSize {
			hvs = append(hvs, hashutil.NewHashValue(sourceSize, sourceHash))
		}
		hv, err := hashutil.CompareHashValuesWithFileName(dstPath, chunkSize, hvs)
		if err == nil && hv != nil {
			return hv.Offset == sourceSize && hv.Hash == sourceHash && len(sourceHash) > 0, hv
		}
	}
	return false, nil
}
