//go:build no_server
// +build no_server

package fs

import (
	"errors"
	"github.com/no-src/gofs/server"
)

// StartFileServer start a file server
func StartFileServer(opt server.Option) error {
	opt.Init.Done()
	return errors.New("file server is not supported. if you need a file server, try to rebuild without 'no_server' tags")
}
