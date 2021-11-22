//go:build !no_server
// +build !no_server

package server

import (
	"github.com/no-src/gofs/contract"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/gofs/util"
	"github.com/no-src/log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

// StartFileServer start a file server
func StartFileServer(src core.VFS, target core.VFS, addr string, init retry.WaitDone) error {
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
		http.Handle(QueryRoute, &fileApiHandler{
			root: http.Dir(src.Path()),
		})
	}

	log.Log("file server [%s] starting...", addr)
	initServerAddr(addr)
	init.Done()
	return http.ListenAndServe(addr, nil)
}

type fileApiHandler struct {
	root http.FileSystem
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
	path := request.FormValue(contract.FsPath)
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

	path = strings.TrimLeft(path, srcPrefix)

	f, err := h.root.Open(path)
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
		files, err := f.Readdir(-1)
		if err != nil {
			log.Error(err, "file server read dir error")
			writer.Write(NewErrorApiResultBytes(-5, "read dir error"))
			return
		}
		for _, file := range files {
			cTime, aTime, mTime, fsTimeErr := util.GetFileTimeBySys(file.Sys())
			if fsTimeErr != nil {
				cTime = time.Now()
				aTime = cTime
				mTime = cTime
			}
			remoteFiles = append(remoteFiles, RemoteFile{
				Path:  file.Name(),
				IsDir: file.IsDir(),
				Size:  file.Size(),
				ATime: aTime.Unix(),
				CTime: cTime.Unix(),
				MTime: mTime.Unix(),
			})
		}
	}
	bytes, err := util.Marshal(NewApiResult(0, "success", remoteFiles))
	if err != nil {
		log.Error(err, "file server marshal error")
		writer.Write(NewErrorApiResultBytes(-6, "marshal error"))
		return
	}
	writer.Write(bytes)
}
