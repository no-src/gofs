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
	"github.com/no-src/gofs/retry"
	"github.com/no-src/log"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// sftpClient a sftp client component, support auto reconnect
type sftpClient struct {
	*sftp.Client

	driverName    string
	remoteAddr    string
	userName      string
	password      string
	sshKey        string
	r             retry.Retry
	mu            sync.RWMutex
	online        bool
	autoReconnect bool
}

// NewSFTPClient get a sftp client
func NewSFTPClient(remoteAddr string, userName string, password string, sshKey string, autoReconnect bool, r retry.Retry) driver.Driver {
	return newSFTPClient(remoteAddr, userName, password, sshKey, autoReconnect, r)
}

func newSFTPClient(remoteAddr string, userName string, password string, sshKey string, autoReconnect bool, r retry.Retry) *sftpClient {
	return &sftpClient{
		driverName:    "sftp",
		remoteAddr:    remoteAddr,
		userName:      userName,
		password:      password,
		sshKey:        sshKey,
		r:             r,
		autoReconnect: autoReconnect,
	}
}

func (sc *sftpClient) DriverName() string {
	return sc.driverName
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
	sc.Client, err = sftp.NewClient(ssh.NewClient(cc, chans, reqs))
	if err == nil {
		sc.online = true
		log.Debug("connect to sftp server success => %s", sc.remoteAddr)
	}
	return err
}

func (sc *sftpClient) getHostKeyCallback() (ssh.HostKeyCallback, error) {
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

func (sc *sftpClient) Create(path string) (err error) {
	err = sc.reconnectIfLost(func() error {
		var f *sftp.File
		f, err = sc.Client.Create(path)
		if err == nil {
			log.ErrorIf(f.Close(), "close sftp file err => %s", path)
		}
		return err
	})
	return err
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

func (sc *sftpClient) Stat(path string) (fi os.FileInfo, err error) {
	err = sc.reconnectIfLost(func() error {
		fi, err = sc.Client.Stat(path)
		return err
	})
	return fi, err
}

func (sc *sftpClient) GetFileTime(path string) (cTime time.Time, aTime time.Time, mTime time.Time, err error) {
	err = sc.reconnectIfLost(func() error {
		var fi os.FileInfo
		fi, err = sc.Client.Stat(path)
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

func (sc *sftpClient) WalkDir(root string, fn fs.WalkDirFunc) error {
	return sc.reconnectIfLost(func() error {
		walker := sc.Client.Walk(root)
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

func (sc *sftpClient) Write(src string, dest string) (err error) {
	err = sc.reconnectIfLost(func() error {
		var srcFile *os.File
		srcFile, err = os.Open(src)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		var destFile *sftp.File
		destFile, err = sc.Client.OpenFile(dest, os.O_WRONLY|os.O_CREATE)
		if err != nil {
			return err
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
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
