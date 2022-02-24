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
	// FsNeedHash return file hash or not
	FsNeedHash = "need_hash"
)

const (
	// ParamValueTrue the parameter value means true
	ParamValueTrue = "1"
	// ParamValueFalse the parameter value means false
	ParamValueFalse = "0"
	// FsNeedHashValueTrue the optional value of the FsNeedHash parameter, that means let file server return file hash value
	FsNeedHashValueTrue = ParamValueTrue
)

// FsDirValue the optional value of FsDir
type FsDirValue int

const (
	// FsIsDir current path is a dir
	FsIsDir FsDirValue = 1
	// FsNotDir current path is not a dir
	FsNotDir FsDirValue = 0
	// FsUnknown current path is unknown file type
	FsUnknown FsDirValue = -1
)

// String parse the current value to string
func (v FsDirValue) String() string {
	return strconv.Itoa(int(v))
}

// Is is current value equal to dest
func (v FsDirValue) Is(t string) bool {
	return v.String() == t
}

// Not is current value not equal to dest
func (v FsDirValue) Not(t string) bool {
	return v.String() != t
}

// Bool current path is a dir or not
func (v FsDirValue) Bool() bool {
	return v == FsIsDir
}

// ParseFsDirValue parse boolean to FsDirValue
func ParseFsDirValue(isDir bool) FsDirValue {
	if isDir {
		return FsIsDir
	}
	return FsNotDir
}
