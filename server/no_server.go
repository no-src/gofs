//go:build no_server
// +build no_server

package server

import (
	"github.com/no-src/gofs/core"
)

import "errors"

// StartFileServer start a file server
func StartFileServer(src core.VFS, target core.VFS, addr string) error {
	return errors.New("file server is not supported. if you need a file server, try to rebuild without 'no_server' tags")
}
