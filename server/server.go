package server

import (
	"net/http"
)

// StartFileServer start a file server
func StartFileServer(src string, target string, addr string) error {
	http.Handle("/", http.FileServer(http.Dir(src)))
	http.Handle("/src/", http.StripPrefix("/src/", http.FileServer(http.Dir(src))))
	http.Handle("/target/", http.StripPrefix("/target/", http.FileServer(http.Dir(target))))
	return http.ListenAndServe(addr, nil)
}
