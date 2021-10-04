package server

import (
	"github.com/no-src/gofs/core"
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
	return http.ListenAndServe(addr, nil)
}
