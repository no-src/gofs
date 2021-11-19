package contract

import (
	"strconv"
)

const (
	// FsDir is dir, see FsDirValue
	FsDir = "dir"
	// FsSize file size, bytes
	FsSize = "size"
	// FsHash file hash value
	FsHash = "hash"
	// FsCtime file create time
	FsCtime = "ctime"
	// FsAtime file last access time
	FsAtime = "atime"
	// FsMtime file last modify time
	FsMtime = "mtime"
	// FsPath file path
	FsPath = "path"
)

type FsDirValue int

const (
	FsIsDir   FsDirValue = 1
	FsNotDir  FsDirValue = 0
	FsUnknown FsDirValue = -1
)

func (v FsDirValue) String() string {
	return strconv.Itoa(int(v))
}

func (v FsDirValue) Is(t string) bool {
	return v.String() == t
}

func (v FsDirValue) Not(t string) bool {
	return v.String() != t
}
