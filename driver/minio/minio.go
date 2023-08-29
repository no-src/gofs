package minio

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/s3utils"
	"github.com/no-src/gofs/driver"
	nsfs "github.com/no-src/gofs/fs"
	"github.com/no-src/gofs/internal/rate"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/log"
)

// minIODriver a MinIO driver component, support auto reconnect
type minIODriver struct {
	client        *minio.Client
	driverName    string
	endpoint      string
	bucketName    string
	secure        bool
	userName      string
	password      string
	r             retry.Retry
	mu            sync.RWMutex
	online        bool
	autoReconnect bool
	ctx           context.Context
	maxTranRate   int64
}

// NewMinIODriver get a MinIO driver
func NewMinIODriver(endpoint string, bucketName string, secure bool, userName string, password string, autoReconnect bool, r retry.Retry, maxTranRate int64) driver.Driver {
	return newMinIODriver(endpoint, bucketName, secure, userName, password, autoReconnect, r, maxTranRate)
}

func newMinIODriver(endpoint string, bucketName string, secure bool, userName string, password string, autoReconnect bool, r retry.Retry, maxTranRate int64) *minIODriver {
	return &minIODriver{
		driverName:    "minio",
		endpoint:      endpoint,
		bucketName:    bucketName,
		secure:        secure,
		userName:      userName,
		password:      password,
		r:             r,
		autoReconnect: autoReconnect,
		ctx:           context.Background(),
		maxTranRate:   maxTranRate,
	}
}

func (c *minIODriver) DriverName() string {
	return c.driverName
}

func (c *minIODriver) Connect() (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.online {
		return nil
	}

	c.client, err = minio.New(c.endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.userName, c.password, ""),
		Secure: c.secure,
	})

	if err != nil {
		return err
	}

	bucketExist, err := c.client.BucketExists(c.ctx, c.bucketName)
	if err != nil {
		return err
	}

	if !bucketExist {
		return fmt.Errorf("bucket %s is not exist", c.bucketName)
	}

	c.online = true
	log.Debug("connect to MinIO server success => %s", c.endpoint)

	return nil
}

func (c *minIODriver) reconnect() error {
	log.Debug("reconnect to MinIO server => %s", c.endpoint)
	return c.r.Do(c.Connect, "MinIO reconnect").Wait()
}

func (c *minIODriver) reconnectIfLost(f func() error) error {
	if !c.autoReconnect {
		return f()
	}
	c.mu.RLock()
	if !c.online {
		c.mu.RUnlock()
		return errors.New("MinIO offline")
	}
	c.mu.RUnlock()

	err := f()
	if c.isClosed(err) {
		log.Error(err, "connect to MinIO server failed")
		c.mu.Lock()
		c.online = false
		c.mu.Unlock()
		if c.reconnect() == nil {
			err = f()
		}
	}
	return err
}

func (c *minIODriver) isClosed(err error) bool {
	return minio.IsNetworkOrHostDown(err, false)
}

func (c *minIODriver) MkdirAll(path string) error {
	return nil
}

func (c *minIODriver) Create(path string) (err error) {
	err = c.reconnectIfLost(func() error {
		_, err = c.client.StatObject(c.ctx, c.bucketName, path, minio.StatObjectOptions{})
		var respErr minio.ErrorResponse
		if err != nil && errors.As(err, &respErr) && respErr.StatusCode == http.StatusNotFound {
			_, err = c.client.PutObject(c.ctx, c.bucketName, path, bytes.NewReader(nil), 0, minio.PutObjectOptions{})
		}
		return err
	})
	return err
}

func (c *minIODriver) Symlink(oldname, newname string) (err error) {
	if err = c.Remove(newname); err != nil {
		return err
	}
	err = c.reconnectIfLost(func() error {
		content := nsfs.SymlinkText(oldname)
		_, err = c.client.PutObject(c.ctx, c.bucketName, newname, bytes.NewReader([]byte(content)), int64(len(content)), minio.PutObjectOptions{})
		return err
	})
	return err
}

func (c *minIODriver) Remove(path string) (err error) {
	return c.reconnectIfLost(func() error {
		infoChan := c.client.ListObjects(c.ctx, c.bucketName, minio.ListObjectsOptions{
			Recursive: true,
			Prefix:    path,
		})
		pathWithSlash := path
		if !strings.HasSuffix(path, "/") {
			pathWithSlash += "/"
		}
		for info := range infoChan {
			if path == info.Key || strings.HasPrefix(info.Key, pathWithSlash) {
				err = c.client.RemoveObject(c.ctx, c.bucketName, info.Key, minio.RemoveObjectOptions{})
				if err != nil {
					return err
				}
			}
		}
		return err
	})
}

func (c *minIODriver) Rename(oldPath, newPath string) error {
	return c.reconnectIfLost(func() error {
		// copy the object then remove the old object
		_, err := c.client.CopyObject(c.ctx, minio.CopyDestOptions{Bucket: c.bucketName, Object: newPath}, minio.CopySrcOptions{Bucket: c.bucketName, Object: oldPath})
		if err == nil {
			err = c.client.RemoveObject(c.ctx, c.bucketName, oldPath, minio.RemoveObjectOptions{})
		}
		return err
	})
}

func (c *minIODriver) Chtimes(path string, aTime time.Time, mTime time.Time) error {
	return nil
}

func (c *minIODriver) Open(path string) (f http.File, err error) {
	err = c.reconnectIfLost(func() error {
		var obj *minio.Object
		obj, err = c.client.GetObject(c.ctx, c.bucketName, path, minio.GetObjectOptions{})
		if err == nil {
			f = rate.NewFile(newFile(obj, c.client, c.bucketName, path), c.maxTranRate)
		}
		return err
	})
	return f, err
}

func (c *minIODriver) openFileOrDir(path string) (f http.File, err error) {
	err = c.reconnectIfLost(func() error {
		infoChan := c.client.ListObjects(c.ctx, c.bucketName, minio.ListObjectsOptions{
			Prefix: path,
		})
		for info := range infoChan {
			if info.Err != nil {
				return err
			}

			if strings.Trim(info.Key, "/") == strings.Trim(path, "/") {
				if strings.HasSuffix(info.Key, "/") {
					f = newDirFile(c.client, c.bucketName, info.Key)
					return nil
				}
				var obj *minio.Object
				obj, err = c.client.GetObject(c.ctx, c.bucketName, info.Key, minio.GetObjectOptions{})
				if err != nil {
					return err
				}
				f = newFile(obj, c.client, c.bucketName, path)
				return nil
			}
			// not matched means path is directory
			err = minio.ErrorResponse{}
			return err
		}
		err = fs.ErrNotExist
		return err
	})
	return f, err
}

func (c *minIODriver) ReadDir(path string) (fis []fs.FileInfo, err error) {
	err = c.reconnectIfLost(func() error {
		infoChan := c.client.ListObjects(c.ctx, c.bucketName, minio.ListObjectsOptions{Recursive: true})
		for info := range infoChan {
			fis = append(fis, newMinIOFileInfo(info))
		}
		return nil
	})
	return fis, err
}

func (c *minIODriver) Stat(path string) (fi fs.FileInfo, err error) {
	err = c.reconnectIfLost(func() error {
		var info minio.ObjectInfo
		info, err = c.client.StatObject(c.ctx, c.bucketName, path, minio.StatObjectOptions{})
		if err != nil {
			return err
		}
		fi = newMinIOFileInfo(info)
		return nil
	})
	return fi, err
}

func (c *minIODriver) Lstat(path string) (fi fs.FileInfo, err error) {
	return c.Stat(path)
}

func (c *minIODriver) GetFileTime(path string) (cTime time.Time, aTime time.Time, mTime time.Time, err error) {
	err = c.reconnectIfLost(func() error {
		var info minio.ObjectInfo
		info, err = c.client.StatObject(c.ctx, c.bucketName, path, minio.StatObjectOptions{})
		if err != nil {
			return err
		}
		cTime = info.LastModified
		aTime = info.LastModified
		mTime = info.LastModified
		return nil
	})
	return
}

func (c *minIODriver) WalkDir(root string, fn fs.WalkDirFunc) error {
	return c.reconnectIfLost(func() error {
		infoChan := c.client.ListObjects(c.ctx, c.bucketName, minio.ListObjectsOptions{Recursive: true})
		for info := range infoChan {
			if err := fn(info.Key, fs.FileInfoToDirEntry(newMinIOFileInfo(info)), info.Err); err != nil {
				return err
			}
		}
		return nil
	})
}

func (c *minIODriver) Write(src string, dest string) (err error) {
	return c.reconnectIfLost(func() error {
		_, err = c.fPutObject(c.ctx, c.bucketName, dest, src, minio.PutObjectOptions{})
		return err
	})
}

func (c *minIODriver) Client() *minio.Client {
	return c.client
}

func (c *minIODriver) ReadLink(path string) (string, error) {
	return path, nil
}

// fPutObject - Create an object in a bucket, with contents from file at filePath. Allows request cancellation.
// Keep up to date with the minio.Client.FPutObject.
func (c *minIODriver) fPutObject(ctx context.Context, bucketName, objectName, filePath string, opts minio.PutObjectOptions) (info minio.UploadInfo, err error) {
	// Input validation.
	if err := s3utils.CheckValidBucketName(bucketName); err != nil {
		return minio.UploadInfo{}, err
	}
	if err := s3utils.CheckValidObjectName(objectName); err != nil {
		return minio.UploadInfo{}, err
	}

	// Open the referenced file.
	fileReader, err := os.Open(filePath)
	// If any error fail quickly here.
	if err != nil {
		return minio.UploadInfo{}, err
	}
	defer fileReader.Close()

	// Save the file stat.
	fileStat, err := fileReader.Stat()
	if err != nil {
		return minio.UploadInfo{}, err
	}

	// Save the file size.
	fileSize := fileStat.Size()

	// Set contentType based on filepath extension if not given or default
	// value of "application/octet-stream" if the extension has no associated type.
	if opts.ContentType == "" {
		if opts.ContentType = mime.TypeByExtension(filepath.Ext(filePath)); opts.ContentType == "" {
			opts.ContentType = "application/octet-stream"
		}
	}
	return c.client.PutObject(ctx, bucketName, objectName, rate.NewReader(fileReader, c.maxTranRate), fileSize, opts)
}
