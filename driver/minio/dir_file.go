package minio

import (
	"context"
	"io"
	"io/fs"
	"net/http"

	"github.com/minio/minio-go/v7"
)

type dirFile struct {
	io.Reader // never called
	io.Seeker // never called

	client     *minio.Client
	bucketName string
	name       string
}

func newDirFile(client *minio.Client, bucketName string, name string) http.File {
	return &dirFile{
		client:     client,
		bucketName: bucketName,
		name:       name,
	}
}

func (f *dirFile) Close() error {
	return nil
}

func (f *dirFile) Readdir(count int) (fis []fs.FileInfo, err error) {
	infoChan := f.client.ListObjects(context.Background(), f.bucketName, minio.ListObjectsOptions{Prefix: f.name})
	for info := range infoChan {
		fis = append(fis, newMinIOFileInfoWithRoot(info, f.name))
	}
	if count > 0 && len(fis) > count {
		fis = fis[:count]
	}
	return fis, nil
}

func (f *dirFile) Stat() (fs.FileInfo, error) {
	return newDirMinIOFileInfo(f.name), nil
}
