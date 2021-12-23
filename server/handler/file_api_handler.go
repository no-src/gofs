package handler

import (
	"github.com/no-src/gofs/contract"
	"github.com/no-src/gofs/server"
	"github.com/no-src/gofs/util"
	"github.com/no-src/log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type fileApiHandler struct {
	root http.FileSystem
}

func NewFileApiHandler(root http.FileSystem) http.Handler {
	return &fileApiHandler{
		root: root,
	}
}

func (h *fileApiHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	defer func() {
		e := recover()
		if e != nil {
			writer.Write(server.NewErrorApiResultBytes(-7, "server internal error"))
		}
	}()
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	var fileList []contract.FileInfo
	path := request.FormValue(contract.FsPath)
	needHash := request.FormValue(contract.FsNeedHash)
	srcPrefix := strings.Trim(server.SrcRoutePrefix, "/")
	targetPrefix := strings.Trim(server.TargetRoutePrefix, "/")
	if !strings.HasPrefix(strings.ToLower(path), srcPrefix) && !strings.HasPrefix(strings.ToLower(path), targetPrefix) {
		writer.Write(server.NewErrorApiResultBytes(-1, "must start with src or target"))
		return
	}

	path = filepath.Clean(path)
	path = filepath.ToSlash(path)
	if !strings.HasPrefix(strings.ToLower(path), srcPrefix) && !strings.HasPrefix(strings.ToLower(path), targetPrefix) {
		writer.Write(server.NewErrorApiResultBytes(-2, "invalid path"))
		return
	}

	path = strings.TrimLeft(path, srcPrefix)

	f, err := h.root.Open(path)
	if err != nil {
		log.Error(err, "file server open path error")
		writer.Write(server.NewErrorApiResultBytes(-3, "open path error"))
		return
	}
	stat, err := f.Stat()
	if err != nil {
		log.Error(err, "file server get file stat error")
		writer.Write(server.NewErrorApiResultBytes(-4, "get file stat error"))
		return
	}
	if stat.IsDir() {
		files, err := f.Readdir(-1)
		if err != nil {
			log.Error(err, "file server read dir error")
			writer.Write(server.NewErrorApiResultBytes(-5, "read dir error"))
			return
		}
		for _, file := range files {
			cTime, aTime, mTime, fsTimeErr := util.GetFileTimeBySys(file.Sys())
			if fsTimeErr != nil {
				cTime = time.Now()
				aTime = cTime
				mTime = cTime
			}

			hash := ""
			if needHash == contract.FsNeedHashValueTrue && !file.IsDir() {
				if cf, err := h.root.Open(filepath.ToSlash(filepath.Join(path, file.Name()))); err == nil {
					hash, _ = util.MD5FromFile(cf, 1024)
				}
			}

			fileList = append(fileList, contract.FileInfo{
				Path:  file.Name(),
				IsDir: contract.ParseFsDirValue(file.IsDir()),
				Size:  file.Size(),
				Hash:  hash,
				ATime: aTime.Unix(),
				CTime: cTime.Unix(),
				MTime: mTime.Unix(),
			})
		}
	}
	bytes, err := util.Marshal(server.NewApiResult(0, "success", fileList))
	if err != nil {
		log.Error(err, "file server marshal error")
		writer.Write(server.NewErrorApiResultBytes(-6, "marshal error"))
		return
	}
	writer.Write(bytes)
}
