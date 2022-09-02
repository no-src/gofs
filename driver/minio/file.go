package minio

import (
	"context"
	"io/fs"
	"net/http"

	"github.com/minio/minio-go/v7"
)

type file struct {
	*minio.Object

	client     *minio.Client
	bucketName string
	name       string
}

func newFile(obj *minio.Object, client *minio.Client, bucketName string, name string) http.File {
	return &file{
		Object:     obj,
		client:     client,
		bucketName: bucketName,
		name:       name,
	}
}

func (f *file) Readdir(count int) (fis []fs.FileInfo, err error) {
	infoChan := f.client.ListObjects(context.Background(), f.bucketName, minio.ListObjectsOptions{Prefix: f.name})
	for info := range infoChan {
		fis = append(fis, newMinIOFileInfo(info))
	}
	if count > 0 && len(fis) > count {
		fis = fis[:count]
	}
	return fis, nil
}

func (f *file) Stat() (fs.FileInfo, error) {
	info, err := f.Object.Stat()
	if err != nil {
		return nil, err
	}
	return newMinIOFileInfo(info), nil
}
