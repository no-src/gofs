package minio

import (
	"io/fs"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
)

type minIOFileInfo struct {
	info minio.ObjectInfo
	root string
}

func newMinIOFileInfo(info minio.ObjectInfo) fs.FileInfo {
	return newMinIOFileInfoWithRoot(info, "")
}

func newMinIOFileInfoWithRoot(info minio.ObjectInfo, root string) fs.FileInfo {
	return &minIOFileInfo{info: info, root: root}
}

func (fi *minIOFileInfo) Name() string {
	return strings.TrimSuffix(strings.TrimPrefix(fi.info.Key, fi.root), "/")
}

func (fi *minIOFileInfo) Size() int64 {
	return fi.info.Size
}

func (fi *minIOFileInfo) Mode() fs.FileMode {
	return 0666
}

func (fi *minIOFileInfo) ModTime() time.Time {
	return fi.info.LastModified
}

func (fi *minIOFileInfo) IsDir() bool {
	return strings.HasSuffix(fi.info.Key, "/")
}

func (fi *minIOFileInfo) Sys() any {
	return nil
}
