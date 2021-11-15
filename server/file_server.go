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
	if src.IsDisk() || src.Is(core.RemoteDisk) {
		http.Handle(SrcRoutePrefix, http.StripPrefix(SrcRoutePrefix, http.FileServer(http.Dir(src.Path()))))
	}
	if target.IsDisk() {
		http.Handle(TargetRoutePrefix, http.StripPrefix(TargetRoutePrefix, http.FileServer(http.Dir(target.Path()))))
	}
	log.Log("file server [%s] starting...", addr)
	initServerAddr(addr)
	return http.ListenAndServe(addr, nil)
}
