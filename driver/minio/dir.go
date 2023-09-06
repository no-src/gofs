package minio

import (
	"errors"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/no-src/gofs/internal/logger"
	"github.com/no-src/gofs/retry"
)

// Dir an implementation of http.FileSystem for MinIO
type Dir struct {
	bucketName string
	driver     *minIODriver
}

// NewDir returns a http.FileSystem instance for MinIO
func NewDir(bucketName string, endpoint string, secure bool, userName string, password string, r retry.Retry, maxTranRate int64, logger *logger.Logger) (http.FileSystem, error) {
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
	driver := newMinIODriver(endpoint, bucketName, secure, userName, password, true, r, maxTranRate, logger)
	return &Dir{
		driver:     driver,
		bucketName: bucketName,
	}, driver.Connect()
}

// Open opens the named file for reading
func (d *Dir) Open(name string) (http.File, error) {
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
		return nil, errors.New("http: invalid character in file path")
	}
	fullName := filepath.ToSlash(filepath.FromSlash(path.Clean("/" + name)))
	httpFile, err := d.driver.openFileOrDir(fullName)
	if err != nil {
		var respErr minio.ErrorResponse
		if errors.As(err, &respErr) && len(respErr.Key) == 0 {
			return newDirFile(d.driver.Client(), d.bucketName, name), nil
		}
		return nil, err
	}
	return httpFile, err
}
