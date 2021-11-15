package sync

import (
	"bufio"
	"errors"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/server"
	"github.com/no-src/gofs/util"
	"github.com/no-src/log"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type remoteClientSync struct {
	src           core.VFS
	target        core.VFS
	srcAbsPath    string
	targetAbsPath string
	bufSize       int
}

func NewRemoteClientSync(src, target core.VFS, bufSize int) (Sync, error) {
	if len(src.Path()) == 0 {
		return nil, errors.New("src is not found")
	}
	if len(target.Path()) == 0 {
		return nil, errors.New("target is not found")
	}
	if bufSize <= 0 {
		return nil, errors.New("bufSize must greater than zero")
	}

	srcAbsPath, err := filepath.Abs(src.Path())
	if err != nil {
		return nil, err
	}

	targetAbsPath, err := filepath.Abs(target.Path())
	if err != nil {
		return nil, err
	}

	rs := &remoteClientSync{
		srcAbsPath:    srcAbsPath,
		targetAbsPath: targetAbsPath,
		bufSize:       bufSize,
		src:           src,
		target:        target,
	}
	return rs, nil
}

func (rs *remoteClientSync) Create(path string) error {
	target, err := rs.buildTargetAbsFile(path)
	if err != nil {
		return err
	}

	isDir, err := rs.IsDir(path)
	if err != nil {
		log.Error(err, "Create:check if the path is dir error")
		return err
	}
	if isDir {
		err = os.MkdirAll(target, os.ModePerm)
		if err != nil {
			log.Error(err, "Create:create dir error")
			return err
		}
	} else {
		dir := filepath.Dir(target)
		err = os.MkdirAll(dir, fs.ModePerm)
		if err != nil {
			log.Error(err, "Create:create dir error")
			return err
		}
		f, err := os.Create(target)
		defer func() {
			if err = f.Close(); err != nil {
				log.Error(err, "Create:close file error")
			}
		}()
		if err != nil {
			log.Error(err, "Create:create file error")
			return err
		}
	}
	_, _, _, aTime, mTime, err := rs.fileInfo(path)
	if err != nil {
		return err
	}
	err = os.Chtimes(target, aTime, mTime)
	if err != nil {
		return err
	}
	log.Debug("create target file success [%s] -> [%s]", path, target)
	return nil
}

func (rs *remoteClientSync) Write(path string) error {
	target, err := rs.buildTargetAbsFile(path)
	if err != nil {
		log.Error(err, "Write:build to target abs file error [%s]", path)
		return err
	}

	isDir, err := rs.IsDir(path)
	if err != nil {
		log.Error(err, "Write:check if the path is dir error")
		return err
	}

	if isDir {
		// ignored
	} else {
		resp, err := http.Get(path)
		if err != nil {
			log.Error(err, "Write:download the src file failed")
			return err
		}
		defer func() {
			if err = resp.Body.Close(); err != nil {
				log.Error(err, "Write:close the resp body error")
			}
		}()

		targetFile, err := os.OpenFile(target, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Error(err, "Write:create the target file failed")
			return err
		}
		defer func() {
			if err = targetFile.Close(); err != nil {
				log.Error(err, "Write:close the target file error")
			}
		}()
		targetStat, err := targetFile.Stat()
		if err != nil {
			log.Error(err, "Write:get the target file stat failed")
			return err
		}

		reader := bufio.NewReader(resp.Body)
		writer := bufio.NewWriter(targetFile)

		size, hash, _, aTime, mTime, err := rs.fileInfo(path)
		if err != nil {
			log.Error(err, "Write:get src file info error")
			return err
		}

		if size == 0 {
			log.Info("write to the target file success [size=%d] [%s] -> [%s]", size, path, target)
			return nil
		}

		// if src and target is the same file, ignore the following steps and return directly
		if size > 0 && size == targetStat.Size() {
			isSame, err := rs.same(hash, targetFile)
			if err == nil && isSame {
				log.Debug("Write:ignored, the file is unmodified")
				return nil
			}
		}

		block := make([]byte, rs.bufSize)
		var wc int64 = 0
		for {
			n, err := reader.Read(block)
			if err == io.EOF && n == 0 {
				break
			}
			if err != nil && err != io.EOF {
				log.Error(err, "Write:read from the src file bytes failed [%s]", path)
				return err
			}
			log.Debug("Write:read from the src file [%d] bytes [%s]", n, path)
			nn, err := writer.Write(block[:n])
			if err != nil {
				log.Error(err, "Write:write to the target file bytes failed [%s]", target)
				return err
			}
			wc += int64(nn)
			progress := float64(wc) / float64(size) * 100
			log.Debug("Write:write to the target file [%d] bytes, current progress [%d/%d][%.2f%%] [%s]", nn, wc, size, progress, target)
		}
		err = writer.Flush()
		if err == nil {
			err = os.Chtimes(target, aTime, mTime)
		}
		if err == nil {
			log.Info("write to the target file success [size=%d] [%s] -> [%s]", size, path, target)
		} else {
			log.Error(err, "Write:flush to the target file failed [%s]", target)
			return err
		}
	}
	return nil
}

func (rs *remoteClientSync) Remove(path string) error {
	target, err := rs.buildTargetAbsFile(path)
	if err != nil {
		log.Error(err, "Remove:build to target abs file error [%s]", path)
		return err
	}
	err = os.RemoveAll(target)
	if err != nil {
		log.Error(err, "Remove:remove the target file error")
	} else {
		log.Info("remove file success [%s] -> [%s]", path, target)
	}
	return err
}

func (rs *remoteClientSync) Rename(path string) error {
	// delete old file, then trigger Create
	return rs.Remove(path)
}

func (rs *remoteClientSync) Chmod(path string) error {
	log.Debug("Chmod not implemented [%s]", path)
	return nil
}

func (rs *remoteClientSync) IsDir(path string) (bool, error) {
	remoteUrl, err := url.Parse(path)
	if err != nil {
		return false, err
	}
	return remoteUrl.Query().Get("dir") == "1", nil
}

func (rs *remoteClientSync) fileInfo(path string) (size int64, hash string, cTime, aTime, mTime time.Time, err error) {
	remoteUrl, err := url.Parse(path)
	if err != nil {
		return
	}
	isDir := remoteUrl.Query().Get("dir") != "0"
	if isDir {
		return
	}

	size, err = strconv.ParseInt(remoteUrl.Query().Get("size"), 10, 64)
	if err != nil {
		return
	}
	hash = remoteUrl.Query().Get("hash")

	cTime = time.Now()
	aTime = time.Now()
	mTime = time.Now()
	cTimeStr := remoteUrl.Query().Get("ctime")
	aTimeStr := remoteUrl.Query().Get("atime")
	mTimeStr := remoteUrl.Query().Get("mtime")
	cTimeL, timeErr := strconv.ParseInt(cTimeStr, 10, 64)
	if timeErr == nil {
		cTime = time.Unix(cTimeL, 0)
	}
	aTimeL, timeErr := strconv.ParseInt(aTimeStr, 10, 64)
	if timeErr == nil {
		aTime = time.Unix(aTimeL, 0)
	}
	mTimeL, timeErr := strconv.ParseInt(mTimeStr, 10, 64)
	if timeErr == nil {
		mTime = time.Unix(mTimeL, 0)
	}
	return
}

func (rs *remoteClientSync) SyncOnce() error {
	log.Debug("SyncOnce not implemented")
	return nil
}

func (rs *remoteClientSync) buildTargetAbsFile(srcFileAbs string) (string, error) {
	remoteUrl, err := url.Parse(srcFileAbs)
	if err != nil {
		log.Error(err, "parse url error, srcFileAbs=%s", srcFileAbs)
		return "", err
	}
	target := filepath.Join(rs.targetAbsPath, strings.TrimPrefix(remoteUrl.Path, server.SrcRoutePrefix))
	return target, nil
}

func (rs *remoteClientSync) same(srcHash string, targetFile *os.File) (bool, error) {
	if len(srcHash) == 0 {
		return false, nil
	}
	targetHash, err := util.MD5FromFile(targetFile, rs.bufSize)
	if err != nil {
		log.Error(err, "calc target file md5 error [%s]", targetFile.Name())
		return false, err
	}

	if len(srcHash) > 0 && srcHash == targetHash {
		return true, nil
	}

	return false, nil
}

func (rs *remoteClientSync) Source() core.VFS {
	return rs.src
}

func (rs *remoteClientSync) Target() core.VFS {
	return rs.target
}
