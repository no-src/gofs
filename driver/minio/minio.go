package minio

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/no-src/gofs/driver"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/log"
)

// minIOClient a MinIO client component, support auto reconnect
type minIOClient struct {
	*minio.Client

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
}

// NewMinIOClient get a MinIO client
func NewMinIOClient(endpoint string, bucketName string, secure bool, userName string, password string, autoReconnect bool, r retry.Retry) driver.Driver {
	return newMinIOClient(endpoint, bucketName, secure, userName, password, autoReconnect, r)
}

func newMinIOClient(endpoint string, bucketName string, secure bool, userName string, password string, autoReconnect bool, r retry.Retry) *minIOClient {
	return &minIOClient{
		driverName:    "MinIO",
		endpoint:      endpoint,
		bucketName:    bucketName,
		secure:        secure,
		userName:      userName,
		password:      password,
		r:             r,
		autoReconnect: autoReconnect,
		ctx:           context.Background(),
	}
}

func (c *minIOClient) DriverName() string {
	return c.driverName
}

func (c *minIOClient) Connect() (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.online {
		return nil
	}

	c.Client, err = minio.New(c.endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.userName, c.password, ""),
		Secure: c.secure,
	})

	if err != nil {
		return err
	}

	bucketExist, err := c.Client.BucketExists(c.ctx, c.bucketName)
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

func (c *minIOClient) reconnect() error {
	log.Debug("reconnect to MinIO server => %s", c.endpoint)
	return c.r.Do(c.Connect, "MinIO reconnect").Wait()
}

func (c *minIOClient) reconnectIfLost(f func() error) error {
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

func (c *minIOClient) isClosed(err error) bool {
	return minio.IsNetworkOrHostDown(err, false)
}

func (c *minIOClient) MkdirAll(path string) error {
	return nil
}

func (c *minIOClient) Create(path string) (err error) {
	err = c.reconnectIfLost(func() error {
		_, err = c.Client.PutObject(c.ctx, c.bucketName, path, bytes.NewReader(nil), 0, minio.PutObjectOptions{})
		return err
	})
	return err
}

func (c *minIOClient) Remove(path string) (err error) {
	return c.reconnectIfLost(func() error {
		infoChan := c.Client.ListObjects(c.ctx, c.bucketName, minio.ListObjectsOptions{
			Recursive: true,
			Prefix:    path,
		})
		pathWithSlash := path
		if !strings.HasSuffix(path, "/") {
			pathWithSlash += "/"
		}
		for info := range infoChan {
			if path == info.Key || strings.HasPrefix(info.Key, pathWithSlash) {
				err = c.Client.RemoveObject(c.ctx, c.bucketName, info.Key, minio.RemoveObjectOptions{})
				if err != nil {
					return err
				}
			}
		}
		return err
	})
}

func (c *minIOClient) Rename(oldPath, newPath string) error {
	return c.reconnectIfLost(func() error {
		// copy the object then remove the old object
		_, err := c.Client.CopyObject(c.ctx, minio.CopyDestOptions{Bucket: c.bucketName, Object: newPath}, minio.CopySrcOptions{Bucket: c.bucketName, Object: oldPath})
		if err == nil {
			err = c.Client.RemoveObject(c.ctx, c.bucketName, oldPath, minio.RemoveObjectOptions{})
		}
		return err
	})
}

func (c *minIOClient) Chtimes(path string, aTime time.Time, mTime time.Time) error {
	return nil
}

func (c *minIOClient) Open(path string) (f http.File, err error) {
	err = c.reconnectIfLost(func() error {
		var obj *minio.Object
		obj, err = c.Client.GetObject(c.ctx, c.bucketName, path, minio.GetObjectOptions{})
		if err == nil {
			f = newFile(obj, c.Client, c.bucketName, path)
		}
		return err
	})
	return f, err
}

func (c *minIOClient) ReadDir(path string) (fis []os.FileInfo, err error) {
	err = c.reconnectIfLost(func() error {
		infoChan := c.Client.ListObjects(c.ctx, c.bucketName, minio.ListObjectsOptions{Recursive: true})
		for info := range infoChan {
			fis = append(fis, newMinIOFileInfo(info))
		}
		return nil
	})
	return fis, err
}

func (c *minIOClient) Stat(path string) (fi os.FileInfo, err error) {
	err = c.reconnectIfLost(func() error {
		var info minio.ObjectInfo
		info, err = c.Client.StatObject(c.ctx, c.bucketName, path, minio.StatObjectOptions{})
		if err != nil {
			return err
		}
		fi = newMinIOFileInfo(info)
		return nil
	})
	return fi, err
}

func (c *minIOClient) GetFileTime(path string) (cTime time.Time, aTime time.Time, mTime time.Time, err error) {
	err = c.reconnectIfLost(func() error {
		var info minio.ObjectInfo
		info, err = c.Client.StatObject(c.ctx, c.bucketName, path, minio.StatObjectOptions{})
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

func (c *minIOClient) WalkDir(root string, fn fs.WalkDirFunc) error {
	return c.reconnectIfLost(func() error {
		infoChan := c.Client.ListObjects(c.ctx, c.bucketName, minio.ListObjectsOptions{Recursive: true})
		for info := range infoChan {
			if err := fn(info.Key, &statDirEntry{newMinIOFileInfo(info)}, info.Err); err != nil {
				return err
			}
		}
		return nil
	})
}

func (c *minIOClient) Write(src string, dest string) (err error) {
	return c.reconnectIfLost(func() error {
		_, err = c.Client.FPutObject(c.ctx, c.bucketName, dest, src, minio.PutObjectOptions{})
		return err
	})
}

type statDirEntry struct {
	info fs.FileInfo
}

func (d *statDirEntry) Name() string               { return d.info.Name() }
func (d *statDirEntry) IsDir() bool                { return d.info.IsDir() }
func (d *statDirEntry) Type() fs.FileMode          { return d.info.Mode().Type() }
func (d *statDirEntry) Info() (fs.FileInfo, error) { return d.info, nil }
