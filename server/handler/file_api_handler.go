package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/no-src/gofs/contract"
	"github.com/no-src/gofs/fs"
	"github.com/no-src/gofs/server"
	"github.com/no-src/gofs/util/hashutil"
	"github.com/no-src/log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type fileApiHandler struct {
	root   http.FileSystem
	logger log.Logger
}

// NewFileApiHandler create an instance of the fileApiHandler
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
			c.JSON(http.StatusOK, server.NewServerErrorResult())
		}
	}()

	var fileList []contract.FileInfo
	path := c.Query(contract.FsPath)
	needHash := c.Query(contract.FsNeedHash)
	sourcePrefix := strings.Trim(server.SourceRoutePrefix, "/")
	destPrefix := strings.Trim(server.DestRoutePrefix, "/")
	if !strings.HasPrefix(strings.ToLower(path), sourcePrefix) && !strings.HasPrefix(strings.ToLower(path), destPrefix) {
		c.JSON(http.StatusOK, server.NewErrorApiResult(-501, "must start with source or dest"))
		return
	}

	path = filepath.Clean(path)
	path = filepath.ToSlash(path)
	if !strings.HasPrefix(strings.ToLower(path), sourcePrefix) && !strings.HasPrefix(strings.ToLower(path), destPrefix) {
		c.JSON(http.StatusOK, server.NewErrorApiResult(-502, "invalid path"))
		return
	}

	path = strings.TrimLeft(path, sourcePrefix)

	f, err := h.root.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			h.logger.Error(err, contract.NotFoundDesc)
			c.JSON(http.StatusOK, server.NewErrorApiResult(contract.NotFound, contract.NotFoundDesc))
		} else {
			h.logger.Error(err, "file server open path error")
			c.JSON(http.StatusOK, server.NewErrorApiResult(-503, "open path error"))
		}
		return
	}
	stat, err := f.Stat()
	if err != nil {
		h.logger.Error(err, "file server get file stat error")
		c.JSON(http.StatusOK, server.NewErrorApiResult(-504, "get file stat error"))
		return
	}

	if stat.IsDir() {
		dirFileList, err := h.readDir(f, needHash, path)
		if err != nil {
			c.JSON(http.StatusOK, server.NewErrorApiResult(-505, "read dir error"))
			return
		}
		fileList = append(fileList, dirFileList...)
	}

	c.JSON(http.StatusOK, server.NewApiResult(contract.Success, contract.SuccessDesc, fileList))
}

func (h *fileApiHandler) readDir(f http.File, needHash string, path string) (fileList []contract.FileInfo, err error) {
	files, err := f.Readdir(-1)
	if err != nil {
		h.logger.Error(err, "file server read dir error")
		return fileList, err
	}
	for _, file := range files {
		cTime, aTime, mTime, fsTimeErr := fs.GetFileTimeBySys(file.Sys())
		if fsTimeErr != nil {
			h.logger.Error(fsTimeErr, "get file times error => %s", file.Name())
			cTime = time.Now()
			aTime = cTime
			mTime = cTime
		}

		hash := ""
		if needHash == contract.FsNeedHashValueTrue && !file.IsDir() {
			if cf, err := h.root.Open(filepath.ToSlash(filepath.Join(path, file.Name()))); err == nil {
				hash, _ = hashutil.MD5FromFile(cf)
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
	return fileList, nil
}
