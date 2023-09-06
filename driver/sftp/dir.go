package sftp

import (
	"errors"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/internal/logger"
	"github.com/no-src/gofs/retry"
)

// Dir an implementation of http.FileSystem for sftp
type Dir struct {
	root   string
	driver *sftpDriver
}

// NewDir returns a http.FileSystem instance for sftp
func NewDir(root string, address string, sshConfig core.SSHConfig, r retry.Retry, maxTranRate int64, logger *logger.Logger) (http.FileSystem, error) {
	root = strings.TrimSpace(root)
	if len(root) == 0 {
		root = "."
	}
	driver := newSFTPDriver(address, sshConfig, true, r, maxTranRate, logger)
	return &Dir{
		driver: driver,
		root:   root,
	}, driver.Connect()
}

// Open opens the named file for reading
func (d *Dir) Open(name string) (http.File, error) {
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
		return nil, errors.New("http: invalid character in file path")
	}
	fullName := filepath.ToSlash(filepath.Join(d.root, filepath.FromSlash(path.Clean("/"+name))))
	return d.driver.Open(fullName)
}
