package fs

import (
	"os"
)

// FileExist is file Exist
func FileExist(path string) (exist bool, err error) {
	_, err = os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// CreateFile create a file without truncate
func CreateFile(name string) (*os.File, error) {
	return os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0666)
}

// OpenRWFile open a file with read write mode
func OpenRWFile(name string) (*os.File, error) {
	return os.OpenFile(name, os.O_RDWR, 0666)
}

// IsDir the path is directory or not
func IsDir(path string) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return f.IsDir(), nil
}
