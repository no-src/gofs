package util

import (
	"errors"
	"os"
)

// FileExist is file Exist
func FileExist(path string) (exist bool, err error) {
	_, err = os.Stat(path)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
