package sftp

import (
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/no-src/gofs/driver"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/log"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// sftpClient a sftp client component, support auto reconnect
type sftpClient struct {
	*sftp.Client

	remoteAddr    string
	userName      string
	password      string
	r             retry.Retry
	mu            sync.RWMutex
	online        bool
	autoReconnect bool
}

// NewSFTPClient get a sftp client
func NewSFTPClient(remoteAddr string, userName string, password string, autoReconnect bool) driver.Driver {
	return newSFTPClient(remoteAddr, userName, password, autoReconnect)
}

func newSFTPClient(remoteAddr string, userName string, password string, autoReconnect bool) *sftpClient {
	return &sftpClient{
		remoteAddr:    remoteAddr,
		userName:      userName,
		password:      password,
		r:             retry.New(100, time.Second, false),
		autoReconnect: autoReconnect,
	}
}

func (sc *sftpClient) Connect() error {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	if sc.online {
		return nil
	}
	c, err := net.Dial("tcp", sc.remoteAddr)
	if err != nil {
		return err
	}
	conf := ssh.ClientConfig{
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		User:            sc.userName,
	}
	conf.Auth = append(conf.Auth, ssh.Password(sc.password))
	cc, chans, reqs, err := ssh.NewClientConn(c, "", &conf)
	if err != nil {
		return err
	}
	sc.Client, err = sftp.NewClient(ssh.NewClient(cc, chans, reqs))
	if err == nil {
		sc.online = true
		log.Debug("connect to sftp server success => %s", sc.remoteAddr)
	}
	return err
}

func (sc *sftpClient) reconnect() error {
	log.Debug("reconnect to sftp server => %s", sc.remoteAddr)
	return sc.r.Do(sc.Connect, "sftp reconnect").Wait()
}

func (sc *sftpClient) reconnectIfLost(f func() error) error {
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

func (sc *sftpClient) isClosed(err error) bool {
	return err == sftp.ErrSSHFxConnectionLost
}

func (sc *sftpClient) MkdirAll(path string) error {
	return sc.reconnectIfLost(func() error {
		return sc.Client.MkdirAll(path)
	})
}

func (sc *sftpClient) Create(path string) (rwc io.ReadWriteCloser, err error) {
	err = sc.reconnectIfLost(func() error {
		rwc, err = sc.Client.Create(path)
		return err
	})
	return rwc, err
}

func (sc *sftpClient) Remove(path string) error {
	return sc.reconnectIfLost(func() error {
		f, err := sc.Client.Stat(path)
		if os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			return err
		}
		var files []string
		var dirs []string

		if f.IsDir() {
			walker := sc.Client.Walk(path)
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
			err = sc.Client.Remove(p)
			if err != nil && !os.IsNotExist(err) {
				return err
			}
		}

		// remove all dirs
		for i := len(dirs) - 1; i >= 0; i-- {
			err = sc.Client.RemoveDirectory(dirs[i])
			if err != nil && !os.IsNotExist(err) {
				return err
			}
		}

		return nil
	})
}

func (sc *sftpClient) Rename(oldPath, newPath string) error {
	return sc.reconnectIfLost(func() error {
		return sc.Client.Rename(oldPath, newPath)
	})
}

func (sc *sftpClient) Chtimes(path string, aTime time.Time, mTime time.Time) error {
	return sc.reconnectIfLost(func() error {
		return sc.Client.Chtimes(path, aTime, mTime)
	})
}

func (sc *sftpClient) Open(path string) (f http.File, err error) {
	err = sc.reconnectIfLost(func() error {
		var sftpFile *sftp.File
		sftpFile, err = sc.Client.Open(path)
		if err == nil {
			f = newFile(sftpFile, sc, path)
		}
		return err
	})
	return f, err
}

func (sc *sftpClient) ReadDir(path string) (fis []os.FileInfo, err error) {
	err = sc.reconnectIfLost(func() error {
		fis, err = sc.Client.ReadDir(path)
		return err
	})
	return fis, err
}
