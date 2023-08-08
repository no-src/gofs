package handler

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/no-src/gofs/action"
	"github.com/no-src/gofs/contract"
	"github.com/no-src/gofs/contract/push"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/fs"
	"github.com/no-src/gofs/server"
	"github.com/no-src/gofs/util/hashutil"
	"github.com/no-src/gofs/util/jsonutil"
	"github.com/no-src/log"
)

type pushHandler struct {
	logger                log.Logger
	storagePath           string
	enableLogicallyDelete bool
	hash                  hashutil.Hash
}

// NewPushHandlerFunc returns a gin.HandlerFunc that to manage the files
func NewPushHandlerFunc(logger log.Logger, source core.VFS, enableLogicallyDelete bool, hash hashutil.Hash) gin.HandlerFunc {
	return (&pushHandler{
		logger:                logger,
		storagePath:           source.Path(),
		enableLogicallyDelete: enableLogicallyDelete,
		hash:                  hash,
	}).Handle
}

func (h *pushHandler) Handle(c *gin.Context) {
	defer func() {
		e := recover()
		if e != nil {
			c.JSON(http.StatusOK, server.NewServerErrorResult())
		}
	}()

	pushDataStr := c.PostForm(push.ParamPushData)
	var pushData push.PushData
	err := jsonutil.Unmarshal([]byte(pushDataStr), &pushData)
	if err != nil {
		msg := "unmarshal push data error"
		c.JSON(http.StatusOK, server.NewErrorApiResult(-501, msg))
		h.logger.Error(err, "%s => %s", msg, pushDataStr)
		return
	}

	h.logger.Debug("receive action %s => %s", pushData.Action.String(), pushDataStr)

	if pushData.Action.Valid() == action.UnknownAction {
		c.JSON(http.StatusOK, server.NewErrorApiResult(-502, fmt.Sprintf("unknown action => %d", pushData.Action.Int())))
		return
	}
	fi := pushData.FileInfo
	switch pushData.Action {
	case action.CreateAction:
		err = h.create(fi)
	case action.SymlinkAction:
		err = h.symlink(fi)
	case action.RemoveAction:
		err = h.remove(fi)
	case action.RenameAction:
		err = h.rename(fi)
	case action.ChmodAction:
		err = h.chmod(fi)
	case action.WriteAction:
		r, _ := h.write(pushData, c)
		c.JSON(http.StatusOK, r)
		return
	default:
		err = fmt.Errorf("unsupported action => [%d:%s]", pushData.Action.Int(), pushData.Action.String())
	}
	if err != nil {
		h.logger.Error(err, "process action error %s => %s", pushData.Action.String(), fi.Path)
		c.JSON(http.StatusOK, server.NewErrorApiResult(-503, fmt.Sprintf("process action error => %s", err.Error())))
	} else {
		c.JSON(http.StatusOK, server.NewApiResult(contract.Success, contract.SuccessDesc, nil))
	}
}

func (h *pushHandler) buildAbsPath(path string) string {
	return filepath.Join(h.storagePath, path)
}

func (h *pushHandler) create(fi contract.FileInfo) error {
	path := h.buildAbsPath(fi.Path)
	exist, err := fs.FileExist(path)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}
	if fi.IsDir.Bool() {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	} else {
		dir := filepath.Dir(path)
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
		f, err := fs.CreateFile(path)
		defer func() {
			if err = f.Close(); err != nil {
				h.logger.Error(err, "close file error")
			}
		}()
		if err != nil {
			return err
		}
	}

	err = h.chtimes(path, fi)
	if err != nil {
		return err
	}
	h.logger.Info("create the dest file success [%s]", path)
	return nil
}

func (h *pushHandler) symlink(fi contract.FileInfo) error {
	path := h.buildAbsPath(fi.Path)
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}
	if err = fs.Symlink(fi.LinkTo, path); err != nil {
		return err
	}
	h.logger.Info("create symlink success [%s] -> [%s]", path, fi.LinkTo)
	return nil
}

func (h *pushHandler) remove(fi contract.FileInfo) (err error) {
	path := h.buildAbsPath(fi.Path)
	if h.enableLogicallyDelete {
		err = fs.LogicallyDelete(path)
	} else {
		err = os.RemoveAll(path)
	}
	if err == nil {
		h.logger.Info("remove file success [%s]", path)
	}
	return err
}

func (h *pushHandler) rename(fi contract.FileInfo) (err error) {
	path := h.buildAbsPath(fi.Path)
	err = os.RemoveAll(path)
	if err == nil {
		h.logger.Info("remove file success [%s]", path)
	}
	return err
}

func (h *pushHandler) chmod(fi contract.FileInfo) (err error) {
	path := h.buildAbsPath(fi.Path)
	h.logger.Debug("chmod is unimplemented [%s]", path)
	return nil
}

func (h *pushHandler) write(pushData push.PushData, c *gin.Context) (server.ApiResult, error) {
	fi := pushData.FileInfo
	if fi.IsDir.Bool() {
		err := errors.New("can't write a directory")
		h.logger.Error(err, "write upload file error")
		return server.NewErrorApiResult(-504, err.Error()), err
	}
	path := h.buildAbsPath(fi.Path)
	fh, err := c.FormFile(push.ParamUpFile)
	if err != nil {
		msg := "get upload file error"
		h.logger.Error(err, msg)
		return server.NewErrorApiResult(-505, msg), err
	}

	code, hv, err := h.Save(fh, path, pushData)
	if err != nil {
		h.logger.Error(err, fmt.Sprintf("save upload file error => [%s]", path))
		return server.NewErrorApiResult(-506, fmt.Sprintf("save upload file error => [%s]", fi.Path)), err
	} else if code != contract.Unknown {
		h.logger.Debug("upload a file that is %s => %s", code.String(), fi.Path)
		return server.NewApiResult(code, code.String(), hv), nil
	}

	// change file times
	if err = h.chtimes(path, fi); err != nil {
		log.Error(err, "change file times error after write file => [%s]", path)
		return server.NewErrorApiResult(-507, fmt.Sprintf("change file times error => [%s]", fi.Path)), err
	}

	return server.NewApiResult(contract.Success, contract.SuccessDesc, nil), nil
}

func (h *pushHandler) chtimes(absPath string, fi contract.FileInfo) error {
	return os.Chtimes(absPath, time.Unix(fi.ATime, 0), time.Unix(fi.MTime, 0))
}

func (h *pushHandler) Save(file *multipart.FileHeader, dst string, pushData push.PushData) (code contract.Code, hv *hashutil.HashValue, err error) {
	offset := pushData.Chunk.Offset
	if pushData.PushAction < push.WritePushAction {
		code, hv = h.compare(dst, pushData)
		return code, hv, nil
	}
	src, err := file.Open()
	if err != nil {
		return code, nil, err
	}
	defer src.Close()

	var out *os.File
	if offset > 0 {
		out, err = fs.CreateFile(dst)
	} else {
		out, err = os.Create(dst)
	}

	if err != nil {
		return code, nil, err
	}
	defer out.Close()

	if offset > 0 {
		if pushData.PushAction == push.TruncatePushAction {
			err = out.Truncate(offset)
		} else {
			_, err = out.Seek(offset, io.SeekStart)
		}
		if err != nil {
			return code, nil, err
		}
	}
	if pushData.PushAction == push.WritePushAction {
		_, err = io.Copy(out, src)
	}
	return code, nil, err
}

func (h *pushHandler) compare(dst string, pushData push.PushData) (contract.Code, *hashutil.HashValue) {
	fileSize := pushData.FileInfo.Size
	chunkSize := pushData.Chunk.Size
	hvs := pushData.FileInfo.HashValues
	pushAction := pushData.PushAction

	if pushAction == push.CompareFilePushAction || pushAction == push.CompareFileAndChunkPushAction {
		destStat, err := os.Stat(dst)
		if err == nil && h.quickCompare(fileSize, destStat.Size(), pushData.FileInfo.MTime, destStat.ModTime().Unix(), pushData.ForceChecksum) {
			log.Debug("[push handler] the file size and file modification time are both unmodified => %s", pushData.FileInfo.Path)
			return contract.NotModified, nil
		}
	}

	switch pushAction {
	case push.CompareFilePushAction:
		if h.compareFile(dst, pushData.FileInfo.Hash, fileSize) {
			return contract.NotModified, nil
		}
	case push.CompareChunkPushAction:
		if h.compareChunk(dst, pushData.Chunk.Offset, pushData.Chunk.Hash, chunkSize) {
			return contract.ChunkNotModified, nil
		}
	case push.CompareFileAndChunkPushAction:
		hv := h.compareHashValues(dst, fileSize, chunkSize, hvs)
		if hv != nil && hv.Offset == fileSize {
			return contract.NotModified, nil
		} else if hv != nil {
			return contract.ChunkNotModified, hv
		}
	}
	if pushAction == push.CompareChunkPushAction {
		return contract.ChunkModified, nil
	}
	return contract.Modified, nil
}

func (h *pushHandler) quickCompare(sourceSize, destSize int64, sourceUnixModTime, destUnixModTime int64, forceChecksum bool) (equal bool) {
	if !forceChecksum && sourceSize == destSize && sourceUnixModTime == destUnixModTime {
		return true
	}
	return false
}

// compareFile compare file size and hash value
func (h *pushHandler) compareFile(dstPath string, sourceHash string, sourceSize int64) (equal bool) {
	if sourceSize > 0 && len(sourceHash) > 0 {
		destStat, err := os.Stat(dstPath)
		if err == nil && destStat.Size() == sourceSize {
			destHash, err := h.hash.HashFromFileName(dstPath)
			if err == nil && destHash == sourceHash {
				return true
			}
		}
	}
	return false
}

// compareChunk compare chunk size and hash value
func (h *pushHandler) compareChunk(dstPath string, offset int64, chunkHash string, chunkSize int64) (equal bool) {
	if chunkSize > 0 && len(chunkHash) > 0 {
		destStat, err := os.Stat(dstPath)
		if err == nil && destStat.Size() >= offset+chunkSize {
			destHash, err := h.hash.HashFromFileChunk(dstPath, offset, chunkSize)
			if err == nil && destHash == chunkHash {
				return true
			}
		}
	}
	return false
}

func (h *pushHandler) compareHashValues(dstPath string, sourceSize int64, chunkSize int64, hvs hashutil.HashValues) *hashutil.HashValue {
	if sourceSize > 0 {
		hv, err := h.hash.CompareHashValuesWithFileName(dstPath, chunkSize, hvs)
		if err == nil && hv != nil {
			return hv
		}
	}
	return nil
}
