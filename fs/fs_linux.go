package fs

import (
	"errors"
	"syscall"
	"time"
)

// GetFileTimeBySys get the creation time, last access time, last modify time of the FileInfo.Sys()
func GetFileTimeBySys(sys any) (cTime time.Time, aTime time.Time, mTime time.Time, err error) {
	if sys != nil {
		attr := sys.(*syscall.Stat_t)
		if attr != nil {
			// fix compile error, syscall.Timespec's members are int32 on linux 386
			cTime = time.Unix(int64(attr.Ctim.Sec), int64(attr.Ctim.Nsec))
			aTime = time.Unix(int64(attr.Atim.Sec), int64(attr.Atim.Nsec))
			mTime = time.Unix(int64(attr.Mtim.Sec), int64(attr.Mtim.Nsec))
		}
	} else {
		err = errors.New("file sys info is nil")
	}
	return
}
