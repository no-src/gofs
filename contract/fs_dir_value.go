package contract

import "strconv"

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
