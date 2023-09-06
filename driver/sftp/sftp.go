package sftp

import (
	"errors"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/driver"
	"github.com/no-src/gofs/internal/rate"
	"github.com/no-src/gofs/logger"
	"github.com/no-src/gofs/retry"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// sftpDriver a sftp driver component, support auto reconnect
type sftpDriver struct {
	client        *sftp.Client
	driverName    string
	remoteAddr    string
	sshConfig     core.SSHConfig
	r             retry.Retry
	mu            sync.RWMutex
	online        bool
	autoReconnect bool
	maxTranRate   int64
	logger        *logger.Logger
}

// NewSFTPDriver get a sftp driver
func NewSFTPDriver(remoteAddr string, sshConfig core.SSHConfig, autoReconnect bool, r retry.Retry, maxTranRate int64, logger *logger.Logger) driver.Driver {
	return newSFTPDriver(remoteAddr, sshConfig, autoReconnect, r, maxTranRate, logger)
}

func newSFTPDriver(remoteAddr string, sshConfig core.SSHConfig, autoReconnect bool, r retry.Retry, maxTranRate int64, logger *logger.Logger) *sftpDriver {
	return &sftpDriver{
		driverName:    "sftp",
		remoteAddr:    remoteAddr,
		sshConfig:     sshConfig,
		r:             r,
		autoReconnect: autoReconnect,
		maxTranRate:   maxTranRate,
		logger:        logger,
	}
}

func (sd *sftpDriver) DriverName() string {
	return sd.driverName
}

func (sd *sftpDriver) Connect() error {
	sd.mu.Lock()
	defer sd.mu.Unlock()
	if sd.online {
		return nil
	}
	if len(sd.sshConfig.Username) == 0 {
		return errors.New("sftp: the username is required")
	}
	c, err := net.Dial("tcp", sd.remoteAddr)
	if err != nil {
		return err
	}
	hostKeyCallback, err := sd.getHostKeyCallback()
	if err != nil {
		return err
	}
	conf := ssh.ClientConfig{
		HostKeyCallback: hostKeyCallback,
		User:            sd.sshConfig.Username,
	}
	if len(sd.sshConfig.Key) > 0 {
		signer, err := sd.getSigner()
		if err != nil {
			return err
		}
		conf.Auth = append(conf.Auth, ssh.PublicKeys(signer))
	}
	if len(sd.sshConfig.Password) > 0 {
		conf.Auth = append(conf.Auth, ssh.Password(sd.sshConfig.Password))
	}

	if len(conf.Auth) == 0 {
		return errors.New("sftp: must set the ssh password or ssh key file")
	}

	cc, chans, reqs, err := ssh.NewClientConn(c, "", &conf)
	if err != nil {
		return err
	}
	sd.client, err = sftp.NewClient(ssh.NewClient(cc, chans, reqs))
	if err == nil {
		sd.online = true
		sd.logger.Debug("connect to sftp server success => %s", sd.remoteAddr)
	}
	return err
}

func (sd *sftpDriver) getSigner() (signer ssh.Signer, err error) {
	pemBytes, err := os.ReadFile(sd.sshConfig.Key)
	if err != nil {
		return nil, err
	}

	if len(sd.sshConfig.KeyPass) > 0 {
		return ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(sd.sshConfig.KeyPass))
	}
	return ssh.ParsePrivateKey(pemBytes)
}

func (sd *sftpDriver) getHostKeyCallback() (ssh.HostKeyCallback, error) {
	if len(sd.sshConfig.HostKey) == 0 {
		sd.logger.Warn("sftp: the ssh host key file is not set, the server's host key validation will be disabled, it may cause a MitM attack")
		return ssh.InsecureIgnoreHostKey(), nil
	}
	// ~/.ssh/known_hosts
	keyData, err := os.ReadFile(sd.sshConfig.HostKey)
	if err != nil {
		return nil, err
	}
	_, _, pk, _, _, err := ssh.ParseKnownHosts(keyData)
	if err != nil {
		return nil, err
	}
	return ssh.FixedHostKey(pk), nil
}

func (sd *sftpDriver) reconnect() error {
	sd.logger.Debug("reconnect to sftp server => %s", sd.remoteAddr)
	return sd.r.Do(sd.Connect, "sftp reconnect").Wait()
}

func (sd *sftpDriver) reconnectIfLost(f func() error) error {
	if !sd.autoReconnect {
		return f()
	}
	sd.mu.RLock()
	if !sd.online {
		sd.mu.RUnlock()
		return errors.New("sftp offline")
	}
	sd.mu.RUnlock()

	err := f()
	if sd.isClosed(err) {
		sd.logger.Error(err, "connect to sftp server failed")
		sd.mu.Lock()
		sd.online = false
		sd.mu.Unlock()
		if sd.reconnect() == nil {
			err = f()
		}
	}
	return err
}

func (sd *sftpDriver) isClosed(err error) bool {
	return err == sftp.ErrSSHFxConnectionLost
}

func (sd *sftpDriver) MkdirAll(path string) error {
	return sd.reconnectIfLost(func() error {
		return sd.client.MkdirAll(path)
	})
}

func (sd *sftpDriver) Create(path string) (err error) {
	err = sd.reconnectIfLost(func() error {
		var f *sftp.File
		f, err = sd.client.Create(path)
		if err == nil {
			sd.logger.ErrorIf(f.Close(), "close sftp file err => %s", path)
		}
		return err
	})
	return err
}

func (sd *sftpDriver) Symlink(oldname, newname string) error {
	if err := sd.Remove(newname); err != nil {
		return err
	}
	return sd.reconnectIfLost(func() error {
		return sd.client.Symlink(oldname, newname)
	})
}

func (sd *sftpDriver) Remove(path string) error {
	return sd.reconnectIfLost(func() error {
		f, err := sd.client.Lstat(path)
		if os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			return err
		}
		var files []string
		var dirs []string

		if f.IsDir() {
			walker := sd.client.Walk(path)
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
			err = sd.client.Remove(p)
			if err != nil && !os.IsNotExist(err) {
				return err
			}
		}

		// remove all dirs
		for i := len(dirs) - 1; i >= 0; i-- {
			err = sd.client.RemoveDirectory(dirs[i])
			if err != nil && !os.IsNotExist(err) {
				return err
			}
		}

		return nil
	})
}

func (sd *sftpDriver) Rename(oldPath, newPath string) error {
	return sd.reconnectIfLost(func() error {
		return sd.client.Rename(oldPath, newPath)
	})
}

func (sd *sftpDriver) Chtimes(path string, aTime time.Time, mTime time.Time) error {
	return sd.reconnectIfLost(func() error {
		return sd.client.Chtimes(path, aTime, mTime)
	})
}

func (sd *sftpDriver) Open(path string) (f http.File, err error) {
	err = sd.reconnectIfLost(func() error {
		var sftpFile *sftp.File
		sftpFile, err = sd.client.Open(path)
		if err == nil {
			f = rate.NewFile(newFile(sftpFile, sd, path), sd.maxTranRate, sd.logger)
		}
		return err
	})
	return f, err
}

func (sd *sftpDriver) ReadDir(path string) (fis []fs.FileInfo, err error) {
	err = sd.reconnectIfLost(func() error {
		fis, err = sd.client.ReadDir(path)
		return err
	})
	return fis, err
}

func (sd *sftpDriver) Stat(path string) (fi fs.FileInfo, err error) {
	err = sd.reconnectIfLost(func() error {
		fi, err = sd.client.Stat(path)
		return err
	})
	return fi, err
}

func (sd *sftpDriver) Lstat(path string) (fi fs.FileInfo, err error) {
	err = sd.reconnectIfLost(func() error {
		fi, err = sd.client.Lstat(path)
		return err
	})
	return fi, err
}

func (sd *sftpDriver) GetFileTime(path string) (cTime time.Time, aTime time.Time, mTime time.Time, err error) {
	err = sd.reconnectIfLost(func() error {
		var fi fs.FileInfo
		fi, err = sd.client.Lstat(path)
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

func (sd *sftpDriver) WalkDir(root string, fn fs.WalkDirFunc) error {
	return sd.reconnectIfLost(func() error {
		walker := sd.client.Walk(root)
		for {
			next := walker.Step()
			if err := walker.Err(); err != nil {
				return err
			}
			if !next {
				return nil
			}
			if err := fn(walker.Path(), fs.FileInfoToDirEntry(walker.Stat()), walker.Err()); err != nil {
				return err
			}
		}
	})
}

func (sd *sftpDriver) Write(src string, dest string) (err error) {
	err = sd.reconnectIfLost(func() error {
		var srcFile *os.File
		srcFile, err = os.Open(src)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		var destFile *sftp.File
		destFile, err = sd.client.OpenFile(dest, os.O_WRONLY|os.O_CREATE)
		if err != nil {
			return err
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, rate.NewReader(srcFile, sd.maxTranRate, sd.logger))
		return err
	})
	return err
}

func (sd *sftpDriver) ReadLink(path string) (realPath string, err error) {
	err = sd.reconnectIfLost(func() error {
		realPath, err = sd.client.ReadLink(path)
		return err
	})
	return realPath, err
}
