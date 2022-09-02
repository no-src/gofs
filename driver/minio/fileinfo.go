package minio

import (
	"io/fs"
	"os"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
)

type minIOFileInfo struct {
	info minio.ObjectInfo
}

func newMinIOFileInfo(info minio.ObjectInfo) fs.FileInfo {
	return &minIOFileInfo{info: info}
}

func (fi *minIOFileInfo) Name() string {
	return fi.info.Key
}

func (fi *minIOFileInfo) Size() int64 {
	return fi.info.Size
}

func (fi *minIOFileInfo) Mode() os.FileMode {
	return 0666
}

func (fi *minIOFileInfo) ModTime() time.Time {
	return fi.info.LastModified
}

func (fi *minIOFileInfo) IsDir() bool {
	if strings.HasSuffix(fi.Name(), "/") {
		return true
	}
	return false
}

func (fi *minIOFileInfo) Sys() any {
	return nil
}
