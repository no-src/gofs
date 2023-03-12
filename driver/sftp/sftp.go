package sftp

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/no-src/gofs/driver"
	"github.com/no-src/gofs/internal/rate"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/log"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// sftpDriver a sftp driver component, support auto reconnect
type sftpDriver struct {
	client        *sftp.Client
	driverName    string
	remoteAddr    string
	userName      string
	password      string
	sshKey        string
	r             retry.Retry
	mu            sync.RWMutex
	online        bool
	autoReconnect bool
	maxTranRate   int64
}

// NewSFTPDriver get a sftp driver
func NewSFTPDriver(remoteAddr string, userName string, password string, sshKey string, autoReconnect bool, r retry.Retry, maxTranRate int64) driver.Driver {
	return newSFTPDriver(remoteAddr, userName, password, sshKey, autoReconnect, r, maxTranRate)
}

func newSFTPDriver(remoteAddr string, userName string, password string, sshKey string, autoReconnect bool, r retry.Retry, maxTranRate int64) *sftpDriver {
	return &sftpDriver{
		driverName:    "sftp",
		remoteAddr:    remoteAddr,
		userName:      userName,
		password:      password,
		sshKey:        sshKey,
		r:             r,
		autoReconnect: autoReconnect,
		maxTranRate:   maxTranRate,
	}
}

func (sc *sftpDriver) DriverName() string {
	return sc.driverName
}

func (sc *sftpDriver) Connect() error {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	if sc.online {
		return nil
	}
	c, err := net.Dial("tcp", sc.remoteAddr)
	if err != nil {
		return err
	}
	hostKeyCallback, err := sc.getHostKeyCallback()
	if err != nil {
		return err
	}
	conf := ssh.ClientConfig{
		HostKeyCallback: hostKeyCallback,
		User:            sc.userName,
	}
	conf.Auth = append(conf.Auth, ssh.Password(sc.password))
	cc, chans, reqs, err := ssh.NewClientConn(c, "", &conf)
	if err != nil {
		return err
	}
	sc.client, err = sftp.NewClient(ssh.NewClient(cc, chans, reqs))
	if err == nil {
		sc.online = true
		log.Debug("connect to sftp server success => %s", sc.remoteAddr)
	}
	return err
}

func (sc *sftpDriver) getHostKeyCallback() (ssh.HostKeyCallback, error) {
	sc.sshKey = strings.TrimSpace(sc.sshKey)
	if len(sc.sshKey) == 0 {
		return ssh.InsecureIgnoreHostKey(), nil
	}
	keyFile, err := os.Open(sc.sshKey)
	if err != nil {
		return nil, err
	}
	defer keyFile.Close()
	keyData, err := io.ReadAll(keyFile)
	if err != nil {
		return nil, err
	}
	keyStat, err := keyFile.Stat()
	if err != nil {
		return nil, err
	}
	keyFileName := strings.ToLower(keyStat.Name())
	var pk ssh.PublicKey
	// ~/.ssh/known_hosts
	if keyFileName == "known_hosts" {
		_, _, pk, _, _, err = ssh.ParseKnownHosts(keyData)
		if err != nil {
			return nil, err
		}
	}
	if pk == nil {
		return nil, fmt.Errorf("invalid ssh key file => %s", keyFileName)
	}
	return ssh.FixedHostKey(pk), nil
}

func (sc *sftpDriver) reconnect() error {
	log.Debug("reconnect to sftp server => %s", sc.remoteAddr)
	return sc.r.Do(sc.Connect, "sftp reconnect").Wait()
}

func (sc *sftpDriver) reconnectIfLost(f func() error) error {
	if !sc.autoReconnect {
		return f()
	}
	sc.mu.RLock()
	if !sc.online {
		sc.mu.RUnlock()
		return errors.New("sftp offline")
	}
	sc.mu.RUnlock()

	err := f()
	if sc.isClosed(err) {
		log.Error(err, "connect to sftp server failed")
		sc.mu.Lock()
		sc.online = false
		sc.mu.Unlock()
		if sc.reconnect() == nil {
			err = f()
		}
	}
	return err
}

func (sc *sftpDriver) isClosed(err error) bool {
	return err == sftp.ErrSSHFxConnectionLost
}

func (sc *sftpDriver) MkdirAll(path string) error {
	return sc.reconnectIfLost(func() error {
		return sc.client.MkdirAll(path)
	})
}

func (sc *sftpDriver) Create(path string) (err error) {
	err = sc.reconnectIfLost(func() error {
		var f *sftp.File
		f, err = sc.client.Create(path)
		if err == nil {
			log.ErrorIf(f.Close(), "close sftp file err => %s", path)
		}
		return err
	})
	return err
}

func (sc *sftpDriver) Remove(path string) error {
	return sc.reconnectIfLost(func() error {
		f, err := sc.client.Stat(path)
		if os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			return err
		}
		var files []string
		var dirs []string

		if f.IsDir() {
			walker := sc.client.Walk(path)
			for walker.Step() {
				if walker.Err() != nil {
					continue
				}
				if walker.Stat().IsDir() {
					dirs = append(dirs, walker.Path())
				} else {
					files = append(files, walker.Path())
				}
			}
		} else {
			files = append(files, path)
		}

		// remove all files
		for _, p := range files {
			err = sc.client.Remove(p)
			if err != nil && !os.IsNotExist(err) {
				return err
			}
		}

		// remove all dirs
		for i := len(dirs) - 1; i >= 0; i-- {
			err = sc.client.RemoveDirectory(dirs[i])
			if err != nil && !os.IsNotExist(err) {
				return err
			}
		}

		return nil
	})
}

func (sc *sftpDriver) Rename(oldPath, newPath string) error {
	return sc.reconnectIfLost(func() error {
		return sc.client.Rename(oldPath, newPath)
	})
}

func (sc *sftpDriver) Chtimes(path string, aTime time.Time, mTime time.Time) error {
	return sc.reconnectIfLost(func() error {
		return sc.client.Chtimes(path, aTime, mTime)
	})
}

func (sc *sftpDriver) Open(path string) (f http.File, err error) {
	err = sc.reconnectIfLost(func() error {
		var sftpFile *sftp.File
		sftpFile, err = sc.client.Open(path)
		if err == nil {
			f = rate.NewFile(newFile(sftpFile, sc, path), sc.maxTranRate)
		}
		return err
	})
	return f, err
}

func (sc *sftpDriver) ReadDir(path string) (fis []os.FileInfo, err error) {
	err = sc.reconnectIfLost(func() error {
		fis, err = sc.client.ReadDir(path)
		return err
	})
	return fis, err
}

func (sc *sftpDriver) Stat(path string) (fi os.FileInfo, err error) {
	err = sc.reconnectIfLost(func() error {
		fi, err = sc.client.Stat(path)
		return err
	})
	return fi, err
}

func (sc *sftpDriver) GetFileTime(path string) (cTime time.Time, aTime time.Time, mTime time.Time, err error) {
	err = sc.reconnectIfLost(func() error {
		var fi os.FileInfo
		fi, err = sc.client.Stat(path)
		if err != nil {
			return err
		}
		cTime = fi.ModTime()
		aTime = fi.ModTime()
		mTime = fi.ModTime()
		return nil
	})
	return
}

func (sc *sftpDriver) WalkDir(root string, fn fs.WalkDirFunc) error {
	return sc.reconnectIfLost(func() error {
		walker := sc.client.Walk(root)
		for {
			next := walker.Step()
			if err := walker.Err(); err != nil {
				return err
			}
			if !next {
				return nil
			}
			if err := fn(walker.Path(), &statDirEntry{walker.Stat()}, walker.Err()); err != nil {
				return err
			}
		}
	})
}

func (sc *sftpDriver) Write(src string, dest string) (err error) {
	err = sc.reconnectIfLost(func() error {
		var srcFile *os.File
		srcFile, err = os.Open(src)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		var destFile *sftp.File
		destFile, err = sc.client.OpenFile(dest, os.O_WRONLY|os.O_CREATE)
		if err != nil {
			return err
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, rate.NewReader(srcFile, sc.maxTranRate))
		return err
	})
	return err
}

type statDirEntry struct {
	info fs.FileInfo
}

func (d *statDirEntry) Name() string               { return d.info.Name() }
func (d *statDirEntry) IsDir() bool                { return d.info.IsDir() }
func (d *statDirEntry) Type() fs.FileMode          { return d.info.Mode().Type() }
func (d *statDirEntry) Info() (fs.FileInfo, error) { return d.info, nil }
