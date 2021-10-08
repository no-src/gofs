//go:build !no_server
// +build !no_server

package server

import (
	"github.com/no-src/gofs/core"
	"github.com/no-src/log"
	"net/http"
)

// StartFileServer start a file server
func StartFileServer(src core.VFS, target core.VFS, addr string) error {
	if src.IsDisk() {
		http.Handle("/", http.FileServer(http.Dir(src.Path())))
		http.Handle("/src/", http.StripPrefix("/src/", http.FileServer(http.Dir(src.Path()))))
	}
	if target.IsDisk() {
		http.Handle("/target/", http.StripPrefix("/target/", http.FileServer(http.Dir(target.Path()))))
	}
	log.Log("file server [%s] starting...", addr)
	return http.ListenAndServe(addr, nil)
}
