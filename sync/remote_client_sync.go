package sync

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/contract"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/server"
	"github.com/no-src/gofs/util"
	"github.com/no-src/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type remoteClientSync struct {
	baseSync
	src           core.VFS
	target        core.VFS
	targetAbsPath string
	currentUser   *auth.User
	cookies       []*http.Cookie
}

// NewRemoteClientSync create an instance of remoteClientSync to receive the file change message and execute it
func NewRemoteClientSync(src, target core.VFS, users []*auth.User, enableLogicallyDelete bool) (Sync, error) {
	if len(target.Path()) == 0 {
		return nil, errors.New("target is not found")
	}

	targetAbsPath, err := filepath.Abs(target.Path())
	if err != nil {
		return nil, err
	}

	rs := &remoteClientSync{
		targetAbsPath: targetAbsPath,
		src:           src,
		target:        target,
		baseSync:      newBaseSync(enableLogicallyDelete),
	}
	if len(users) > 0 {
		rs.currentUser = users[0]
	}
	return rs, nil
}

func (rs *remoteClientSync) Create(path string) error {
	target, err := rs.buildTargetAbsFile(path)
	if err != nil {
		return err
	}

	exist, err := util.FileExist(target)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}

	isDir, err := rs.IsDir(path)
	if err != nil {
		return err
	}
	if isDir {
		err = os.MkdirAll(target, os.ModePerm)
		if err != nil {
			return err
		}
	} else {
		dir := filepath.Dir(target)
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
		f, err := util.CreateFile(target)
		defer func() {
			if err = f.Close(); err != nil {
				log.Error(err, "Create:close file error")
			}
		}()
		if err != nil {
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
	log.Info("create the target file success [%s] -> [%s]", path, target)
	return nil
}

func (rs *remoteClientSync) Write(path string) error {
	target, err := rs.buildTargetAbsFile(path)
	if err != nil {
		return err
	}

	isDir, err := rs.IsDir(path)
	if err != nil {
		return err
	}

	if isDir {
		return rs.SyncOnce(path)
	} else {
		resp, err := rs.httpGetWithAuth(path)
		if err != nil {
			return err
		}
		defer func() {
			if err = resp.Body.Close(); err != nil {
				log.Error(err, "Write:close the resp body error")
			}
		}()

		targetFile, err := util.OpenRWFile(target)
		if err != nil {
			return err
		}
		defer func() {
			if err = targetFile.Close(); err != nil {
				log.Error(err, "Write:close the target file error")
			}
		}()
		targetStat, err := targetFile.Stat()
		if err != nil {
			return err
		}

		reader := bufio.NewReader(resp.Body)
		writer := bufio.NewWriter(targetFile)

		size, hash, _, aTime, mTime, err := rs.fileInfo(path)
		if err != nil {
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
				log.Debug("Write:ignored, the file is unmodified => %s", path)
				return nil
			}
			// reset the offset
			if _, err = targetFile.Seek(0, 0); err != nil {
				return err
			}
		}

		// truncate first before write to file
		err = targetFile.Truncate(0)
		if err != nil {
			return err
		}

		n, err := reader.WriteTo(writer)
		if err != nil {
			return err
		}

		err = writer.Flush()

		if err == nil {
			log.Info("write to the target file success, size[%d => %d] [%s] => [%s]", size, n, path, target)

			// change file times
			if err := os.Chtimes(target, aTime, mTime); err != nil {
				log.Warn("Write:change file times error => %s =>[%s]", err.Error(), target)
			}
		} else {
			return err
		}
	}
	return nil
}

func (rs *remoteClientSync) Remove(path string) error {
	return rs.remove(path, false)
}

func (rs *remoteClientSync) remove(path string, forceDelete bool) error {
	target, err := rs.buildTargetAbsFile(path)
	if err != nil {
		return err
	}
	if !forceDelete && rs.enableLogicallyDelete {
		err = rs.logicallyDelete(target)
	} else {
		err = os.RemoveAll(target)
	}
	if err == nil {
		log.Info("remove file success [%s] -> [%s]", path, target)
	}
	return err
}

func (rs *remoteClientSync) Rename(path string) error {
	// delete old file, then trigger Create
	return rs.remove(path, true)
}

func (rs *remoteClientSync) Chmod(path string) error {
	log.Debug("Chmod is unimplemented [%s]", path)
	return nil
}

func (rs *remoteClientSync) IsDir(path string) (bool, error) {
	remoteUrl, err := url.Parse(path)
	if err != nil {
		return false, err
	}
	return contract.FsIsDir.Is(remoteUrl.Query().Get(contract.FsDir)), nil
}

func (rs *remoteClientSync) fileInfo(path string) (size int64, hash string, cTime, aTime, mTime time.Time, err error) {
	remoteUrl, err := url.Parse(path)
	if err != nil {
		return
	}
	isDir := contract.FsNotDir.Not(remoteUrl.Query().Get(contract.FsDir))
	if isDir {
		return
	}
	size, err = util.Int64(remoteUrl.Query().Get(contract.FsSize))
	if err != nil {
		return
	}
	hash = remoteUrl.Query().Get(contract.FsHash)

	cTime = time.Now()
	aTime = time.Now()
	mTime = time.Now()
	cTimeStr := remoteUrl.Query().Get(contract.FsCtime)
	aTimeStr := remoteUrl.Query().Get(contract.FsAtime)
	mTimeStr := remoteUrl.Query().Get(contract.FsMtime)
	cTimeL, timeErr := util.Int64(cTimeStr)
	if timeErr == nil {
		cTime = time.Unix(cTimeL, 0)
	}
	aTimeL, timeErr := util.Int64(aTimeStr)
	if timeErr == nil {
		aTime = time.Unix(aTimeL, 0)
	}
	mTimeL, timeErr := util.Int64(mTimeStr)
	if timeErr == nil {
		mTime = time.Unix(mTimeL, 0)
	}
	return
}

func (rs *remoteClientSync) SyncOnce(path string) error {
	remoteUrl, err := url.Parse(path)
	if err != nil {
		return err
	}
	syncPath := strings.Trim(remoteUrl.Path, "/")
	return rs.sync(fmt.Sprintf("%s://%s", remoteUrl.Scheme, remoteUrl.Host), syncPath)
}

func (rs *remoteClientSync) sync(serverAddr, path string) error {
	log.Debug("remote client sync path => %s", path)
	reqValues := url.Values{}
	reqValues.Add(contract.FsPath, path)
	reqValues.Add(contract.FsNeedHash, contract.FsNeedHashValueTrue)
	queryUrl := fmt.Sprintf("%s%s?%s", serverAddr, server.QueryRoute, reqValues.Encode())
	resp, err := rs.httpGetWithAuth(queryUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var apiResult server.ApiResult
	err = util.Unmarshal(data, &apiResult)
	if err != nil {
		return err
	}
	if apiResult.Code != 0 {
		return errors.New(fmt.Sprintf("query error:%s", apiResult.Message))
	}
	if apiResult.Data == nil {
		return nil
	}
	dataBytes, err := util.Marshal(apiResult.Data)
	if err != nil {
		return err
	}
	var files []contract.FileInfo
	err = util.Unmarshal(dataBytes, &files)
	if err != nil {
		return err
	}
	for _, file := range files {
		currentPath := path + "/" + file.Path
		values := url.Values{}
		values.Add(contract.FsDir, file.IsDir.String())
		values.Add(contract.FsSize, util.String(file.Size))
		values.Add(contract.FsHash, file.Hash)
		values.Add(contract.FsCtime, util.String(file.CTime))
		values.Add(contract.FsAtime, util.String(file.ATime))
		values.Add(contract.FsMtime, util.String(file.MTime))
		syncPath := fmt.Sprintf("%s/%s?%s", serverAddr, currentPath, values.Encode())

		// create directory or file
		if err = rs.Create(syncPath); err != nil {
			log.Error(err, "sync create directory or file error => [syncPath=%s]", syncPath)
		}

		if file.IsDir.Bool() {
			// sync current directory content
			if err = rs.sync(serverAddr, currentPath); err != nil {
				log.Error(err, "sync current directory content error => [serverAddr=%s] [currentPath=%s]", serverAddr, currentPath)
			}
		} else {
			// sync remote file to local disk
			if err = rs.Write(syncPath); err != nil {
				log.Error(err, "sync remote file to local disk error => [syncPath=%s]", syncPath)
			}
		}
	}
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
	targetHash, err := util.MD5FromFile(targetFile)
	if err != nil {
		log.Error(err, "calculate md5 hash of the target file error [%s]", targetFile.Name())
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

func (rs *remoteClientSync) httpGetWithAuth(rawURL string) (resp *http.Response, err error) {
	resp, err = util.HttpGetWithCookie(rawURL, rs.cookies...)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusUnauthorized && rs.currentUser != nil {
		// auto login
		parseUrl, err := url.Parse(rawURL)
		if err != nil {
			return nil, err
		}
		loginUrl := fmt.Sprintf("%s://%s%s", parseUrl.Scheme, parseUrl.Host, server.LoginSignInFullRoute)
		form := url.Values{}
		user := rs.currentUser
		form.Set(server.ServerParamUserName, user.UserName())
		form.Set(server.ServerParamPassword, user.Password())
		log.Debug("try to auto login file server %s=%s %s=%s", server.ServerParamUserName, user.UserName(), server.ServerParamPassword, user.Password())
		loginResp, err := util.HttpPostWithoutRedirect(loginUrl, form)
		if err != nil {
			return nil, err
		}
		if loginResp.StatusCode == http.StatusFound {
			rs.cookies = loginResp.Cookies()
			if len(rs.cookies) > 0 {
				log.Debug("try to auto login file server success maybe, retry to get resource => %s", rawURL)
				return util.HttpGetWithCookie(rawURL, rs.cookies...)
			}
		}
		return nil, errors.New("file server is unauthorized")
	}
	return resp, err
}
