package fs

import (
	"syscall"
	"time"
)

// GetFileTimeBySys get the creation time, last access time, last modify time of the FileInfo.Sys()
func GetFileTimeBySys(sys any) (cTime time.Time, aTime time.Time, mTime time.Time, err error) {
	if sys != nil {
		attr := sys.(*syscall.Stat_t)
		if attr != nil {
			cTime = time.Unix(attr.Ctimespec.Sec, attr.Ctimespec.Nsec)
			aTime = time.Unix(attr.Atimespec.Sec, attr.Atimespec.Nsec)
			mTime = time.Unix(attr.Mtimespec.Sec, attr.Mtimespec.Nsec)
		}
	} else {
		err = errFileSysInfoIsNil
	}
	return
}
