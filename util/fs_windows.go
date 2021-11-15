package util

import (
	"errors"
	"os"
	"syscall"
	"time"
)

func GetFileTime(path string) (cTime time.Time, aTime time.Time, mTime time.Time, err error) {
	stat, err := os.Stat(path)
	if err != nil {
		return
	}
	if stat.Sys() != nil {
		attr := stat.Sys().(*syscall.Win32FileAttributeData)
		if attr != nil {
			cTime = time.Unix(0, attr.CreationTime.Nanoseconds())
			aTime = time.Unix(0, attr.LastAccessTime.Nanoseconds())
			mTime = time.Unix(0, attr.LastWriteTime.Nanoseconds())
		}
	} else {
		err = errors.New("file sys info is nil")
	}
	return
}
