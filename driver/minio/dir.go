package minio

import (
	"errors"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/no-src/gofs/retry"
)

// Dir an implementation of http.FileSystem for MinIO
type Dir struct {
	bucketName string
	client     *minIOClient
}

// NewDir returns a http.FileSystem instance for MinIO
func NewDir(bucketName string, endpoint string, secure bool, userName string, password string, r retry.Retry) (http.FileSystem, error) {
	bucketName = strings.TrimSpace(bucketName)
	if len(bucketName) == 0 {
		return nil, errors.New("the bucket can't be empty")
	}
	userName = strings.TrimSpace(userName)
	if len(userName) == 0 {
		return nil, errors.New("invalid username for MinIO")
	}
	password = strings.TrimSpace(password)
	if len(password) == 0 {
		return nil, errors.New("invalid password for MinIO")
	}
	client := newMinIOClient(endpoint, bucketName, secure, userName, password, true, r)
	return &Dir{
		client:     client,
		bucketName: bucketName,
	}, client.Connect()
}

// Open opens the named file for reading
func (d *Dir) Open(name string) (http.File, error) {
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
		return nil, errors.New("http: invalid character in file path")
	}
	fullName := filepath.ToSlash(filepath.FromSlash(path.Clean("/" + name)))
	httpFile, err := d.client.Open(fullName)
	if err != nil {
		return nil, err
	}
	_, err = httpFile.Stat()
	if err != nil {
		var respErr minio.ErrorResponse
		if errors.As(err, &respErr) && len(respErr.Key) == 0 {
			return nil, errors.New("to list directory is unsupported")
		}
		return nil, err
	}
	return httpFile, err
}
