package fs

import (
	"errors"
	"syscall"
	"time"
)

// GetFileTimeBySys get the creation time, last access time, last modify time of the FileInfo.Sys()
func GetFileTimeBySys(sys any) (cTime time.Time, aTime time.Time, mTime time.Time, err error) {
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
