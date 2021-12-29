package handler

import (
	"github.com/gin-gonic/gin"
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
	root   http.FileSystem
	logger log.Logger
}

func NewFileApiHandler(root http.FileSystem, logger log.Logger) GinHandler {
	return &fileApiHandler{
		root:   root,
		logger: logger,
	}
}

func (h *fileApiHandler) Handle(c *gin.Context) {
	defer func() {
		e := recover()
		if e != nil {
			c.JSON(http.StatusOK, server.NewErrorApiResult(-7, "server internal error"))
		}
	}()

	var fileList []contract.FileInfo
	path := c.Query(contract.FsPath)
	needHash := c.Query(contract.FsNeedHash)
	srcPrefix := strings.Trim(server.SrcRoutePrefix, "/")
	targetPrefix := strings.Trim(server.TargetRoutePrefix, "/")
	if !strings.HasPrefix(strings.ToLower(path), srcPrefix) && !strings.HasPrefix(strings.ToLower(path), targetPrefix) {
		c.JSON(http.StatusOK, server.NewErrorApiResult(-1, "must start with src or target"))
		return
	}

	path = filepath.Clean(path)
	path = filepath.ToSlash(path)
	if !strings.HasPrefix(strings.ToLower(path), srcPrefix) && !strings.HasPrefix(strings.ToLower(path), targetPrefix) {
		c.JSON(http.StatusOK, server.NewErrorApiResult(-2, "invalid path"))
		return
	}

	path = strings.TrimLeft(path, srcPrefix)

	f, err := h.root.Open(path)
	if err != nil {
		h.logger.Error(err, "file server open path error")
		c.JSON(http.StatusOK, server.NewErrorApiResult(-3, "open path error"))
		return
	}
	stat, err := f.Stat()
	if err != nil {
		h.logger.Error(err, "file server get file stat error")
		c.JSON(http.StatusOK, server.NewErrorApiResult(-4, "get file stat error"))
		return
	}
	if stat.IsDir() {
		files, err := f.Readdir(-1)
		if err != nil {
			h.logger.Error(err, "file server read dir error")
			c.JSON(http.StatusOK, server.NewErrorApiResult(-5, "read dir error"))
			return
		}
		for _, file := range files {
			cTime, aTime, mTime, fsTimeErr := util.GetFileTimeBySys(file.Sys())
			if fsTimeErr != nil {
				h.logger.Error(fsTimeErr, "get file times error => %s", file.Name())
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

	c.JSON(http.StatusOK, server.NewApiResult(0, "success", fileList))
}
