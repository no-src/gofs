//go:build no_server
// +build no_server

package fs

import (
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/retry"
)

import "errors"

// StartFileServer start a file server
func StartFileServer(src core.VFS, target core.VFS, addr string, init retry.WaitDone, enableTLS bool, certFile string, keyFile string) error {
	init.Done()
	return errors.New("file server is not supported. if you need a file server, try to rebuild without 'no_server' tags")
}
