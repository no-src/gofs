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
	return GetFileTimeBySys(stat.Sys())
}

func GetFileTimeBySys(sys interface{}) (cTime time.Time, aTime time.Time, mTime time.Time, err error) {
	if sys != nil {
		attr := sys.(*syscall.Win32FileAttributeData)
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
