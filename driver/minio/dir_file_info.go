package minio

import (
	"io/fs"
	"time"
)

type dirMinIOFileInfo struct {
	name string
}

func newDirMinIOFileInfo(name string) fs.FileInfo {
	return &dirMinIOFileInfo{
		name: name,
	}
}

func (fi *dirMinIOFileInfo) Name() string {
	return fi.name
}

func (fi *dirMinIOFileInfo) Size() int64 {
	return 0
}

func (fi *dirMinIOFileInfo) Mode() fs.FileMode {
	return 0755
}

func (fi *dirMinIOFileInfo) ModTime() time.Time {
	return time.Now()
}

func (fi *dirMinIOFileInfo) IsDir() bool {
	return true
}

func (fi *dirMinIOFileInfo) Sys() any {
	return nil
}
