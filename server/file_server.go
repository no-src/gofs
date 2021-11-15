//go:build !no_server
// +build !no_server

package server

import (
	"encoding/json"
	"github.com/no-src/gofs/core"
	"github.com/no-src/log"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// StartFileServer start a file server
func StartFileServer(src core.VFS, target core.VFS, addr string) error {
	enableFileApi := false
	if src.IsDisk() || src.Is(core.RemoteDisk) {
		http.Handle(SrcRoutePrefix, http.StripPrefix(SrcRoutePrefix, http.FileServer(http.Dir(src.Path()))))
		enableFileApi = true
	}

	if target.IsDisk() {
		http.Handle(TargetRoutePrefix, http.StripPrefix(TargetRoutePrefix, http.FileServer(http.Dir(target.Path()))))
		enableFileApi = true
	}

	if enableFileApi {
		http.Handle("/query", &fileApiHandler{})
	}

	log.Log("file server [%s] starting...", addr)
	initServerAddr(addr)
	return http.ListenAndServe(addr, nil)
}

type fileApiHandler struct {
}

func (h *fileApiHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	defer func() {
		e := recover()
		if e != nil {
			writer.Write(NewErrorApiResultBytes(-7, "server internal error"))
		}
	}()
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	var remoteFiles []RemoteFile
	path := request.FormValue("path")
	srcPrefix := strings.Trim(SrcRoutePrefix, "/")
	targetPrefix := strings.Trim(TargetRoutePrefix, "/")
	if !strings.HasPrefix(strings.ToLower(path), srcPrefix) && !strings.HasPrefix(strings.ToLower(path), targetPrefix) {
		writer.Write(NewErrorApiResultBytes(-1, "must start with src or target"))
		return
	}
	path = filepath.Clean(path)
	path = filepath.ToSlash(path)
	if !strings.HasPrefix(strings.ToLower(path), srcPrefix) && !strings.HasPrefix(strings.ToLower(path), targetPrefix) {
		writer.Write(NewErrorApiResultBytes(-2, "invalid path"))
		return
	}
	f, err := os.Open(path)
	if err != nil {
		log.Error(err, "file server open path error")
		writer.Write(NewErrorApiResultBytes(-3, "open path error"))
		return
	}
	stat, err := f.Stat()
	if err != nil {
		log.Error(err, "file server get file stat error")
		writer.Write(NewErrorApiResultBytes(-4, "get file stat error"))
		return
	}
	if stat.IsDir() {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			log.Error(err, "file server read dir error")
			writer.Write(NewErrorApiResultBytes(-5, "read dir error"))
			return
		}
		for _, file := range files {
			remoteFiles = append(remoteFiles, RemoteFile{
				Path:  file.Name(),
				IsDir: file.IsDir(),
				Size:  file.Size(),
			})
		}
	}
	bytes, err := json.Marshal(NewApiResult(0, "success", remoteFiles))
	if err != nil {
		log.Error(err, "file server marshal error")
		writer.Write(NewErrorApiResultBytes(-6, "marshal error"))
		return
	}
	writer.Write(bytes)
}
