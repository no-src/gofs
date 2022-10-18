package sync

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/contract"
	"github.com/no-src/gofs/fs"
	"github.com/no-src/gofs/ignore"
	"github.com/no-src/gofs/server"
	"github.com/no-src/gofs/server/client"
	"github.com/no-src/gofs/util/hashutil"
	"github.com/no-src/gofs/util/httputil"
	"github.com/no-src/gofs/util/jsonutil"
	"github.com/no-src/gofs/util/stringutil"
	"github.com/no-src/log"
)

type remoteClientSync struct {
	baseSync

	destAbsPath string
	currentUser *auth.User
	cookies     []*http.Cookie
	chunkSize   int64
}

// NewRemoteClientSync create an instance of remoteClientSync to receive the file change message and execute it
func NewRemoteClientSync(opt Option) (Sync, error) {
	// the fields of option
	source := opt.Source
	dest := opt.Dest
	users := opt.Users
	chunkSize := opt.ChunkSize
	forceChecksum := opt.ForceChecksum
	enableLogicallyDelete := opt.EnableLogicallyDelete

	if dest.IsEmpty() {
		return nil, errors.New("dest is not found")
	}

	destAbsPath, err := dest.Abs()
	if err != nil {
		return nil, err
	}

	rs := &remoteClientSync{
		destAbsPath: destAbsPath,
		baseSync:    newBaseSync(source, dest, enableLogicallyDelete, forceChecksum),
		chunkSize:   chunkSize,
	}
	if len(users) > 0 {
		rs.currentUser = users[0]
	}
	return rs, nil
}

func (rs *remoteClientSync) Create(path string) error {
	dest, err := rs.buildDestAbsFile(path)
	if err != nil {
		return err
	}

	exist, err := fs.FileExist(dest)
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
		err = os.MkdirAll(dest, os.ModePerm)
		if err != nil {
			return err
		}
	} else {
		dir := filepath.Dir(dest)
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
		f, err := fs.CreateFile(dest)
		defer func() {
			log.ErrorIf(f.Close(), "[create] close the dest file error")
		}()
		if err != nil {
			return err
		}
	}
	_, _, _, _, aTime, mTime, err := rs.fileInfo(path)
	if err != nil {
		return err
	}
	err = os.Chtimes(dest, aTime, mTime)
	if err != nil {
		return err
	}
	log.Info("create the dest file success [%s] -> [%s]", path, dest)
	return nil
}

func (rs *remoteClientSync) Write(path string) error {
	dest, err := rs.buildDestAbsFile(path)
	if err != nil {
		return err
	}

	isDir, err := rs.IsDir(path)
	if err != nil {
		return err
	}

	// process directory
	if isDir {
		return rs.SyncOnce(path)
	}

	// process file
	return rs.write(path, dest)
}

// write try to write a file to the destination
func (rs *remoteClientSync) write(path, dest string) error {
	size, hash, hvs, _, aTime, mTime, err := rs.fileInfo(path)
	if err != nil {
		return err
	}

	destStat, err := os.Stat(dest)
	if err == nil && rs.quickCompare(size, destStat.Size(), mTime, destStat.ModTime()) {
		log.Debug("[remote client sync] [write] [ignored], the file size and file modification time are both unmodified => %s", path)
		return nil
	}

	// if source and dest is the same file, ignore the following steps and return directly
	equal, hv := rs.compareHashValues(dest, size, hash, rs.chunkSize, hvs)
	if equal {
		log.Debug("[remote client sync] [write] [ignored], the file is unmodified => %s", path)
		return nil
	}
	var offset int64
	rangeHeader := make(http.Header)
	if hv != nil {
		offset = hv.Offset
		rangeHeader.Add("Range", fmt.Sprintf("bytes=%d-%d", offset, size))
	}
	resp, err := rs.httpGetWithAuth(path, rangeHeader)
	if err != nil {
		return err
	}
	defer func() {
		log.ErrorIf(resp.Body.Close(), "[remote client sync] [write] close the resp body error")
	}()

	destFile, err := fs.OpenRWFile(dest)
	if err != nil {
		return err
	}
	defer func() {
		log.ErrorIf(destFile.Close(), "[remote client sync] [write] close the dest file error")
	}()

	if _, err = destFile.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	reader := bufio.NewReader(resp.Body)
	writer := bufio.NewWriter(destFile)

	// truncate first before write to file
	err = destFile.Truncate(offset)
	if err != nil {
		return err
	}

	n, err := reader.WriteTo(writer)
	if err != nil {
		return err
	}

	err = writer.Flush()

	if err == nil {
		log.Info("[remote-client] [write] [success] size[%d => %d] [%s] => [%s]", size, n, path, dest)
		rs.chtimes(dest, aTime, mTime)
	}
	return err
}

// chtimes change file times
func (rs *remoteClientSync) chtimes(dest string, aTime, mTime time.Time) {
	if err := os.Chtimes(dest, aTime, mTime); err != nil {
		log.Warn("[remote client sync] change file times error => %s =>[%s]", err.Error(), dest)
	}
}

func (rs *remoteClientSync) Remove(path string) error {
	return rs.remove(path, false)
}

func (rs *remoteClientSync) remove(path string, forceDelete bool) error {
	dest, err := rs.buildDestAbsFile(path)
	if err != nil {
		return err
	}
	if !forceDelete && rs.enableLogicallyDelete {
		err = rs.logicallyDelete(dest)
	} else {
		err = os.RemoveAll(dest)
	}
	if err == nil {
		log.Info("remove file success [%s] -> [%s]", path, dest)
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

func (rs *remoteClientSync) fileInfo(path string) (size int64, hash string, hvs hashutil.HashValues, cTime, aTime, mTime time.Time, err error) {
	remoteUrl, err := url.Parse(path)
	if err != nil {
		return
	}
	isDir := contract.FsNotDir.Not(remoteUrl.Query().Get(contract.FsDir))
	if isDir {
		return
	}
	size, err = stringutil.Int64(remoteUrl.Query().Get(contract.FsSize))
	if err != nil {
		return
	}
	hash = remoteUrl.Query().Get(contract.FsHash)
	hashValues := remoteUrl.Query().Get(contract.FsHashValues)
	if len(hashValues) > 0 {
		if err = jsonutil.Unmarshal([]byte(hashValues), &hvs); err != nil {
			return
		}
	}

	cTime = time.Now()
	aTime = time.Now()
	mTime = time.Now()
	cTimeStr := remoteUrl.Query().Get(contract.FsCtime)
	aTimeStr := remoteUrl.Query().Get(contract.FsAtime)
	mTimeStr := remoteUrl.Query().Get(contract.FsMtime)
	cTimeL, timeErr := stringutil.Int64(cTimeStr)
	if timeErr == nil {
		cTime = time.Unix(cTimeL, 0)
	}
	aTimeL, timeErr := stringutil.Int64(aTimeStr)
	if timeErr == nil {
		aTime = time.Unix(aTimeL, 0)
	}
	mTimeL, timeErr := stringutil.Int64(mTimeStr)
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
	resp, err := rs.httpGetWithAuth(queryUrl, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var apiResult server.ApiResult
	err = jsonutil.Unmarshal(data, &apiResult)
	if err != nil {
		return err
	}
	if apiResult.Code == contract.NotFound {
		// cancel retry to write when the file does not exist
		return os.ErrNotExist
	} else if apiResult.Code != contract.Success {
		return fmt.Errorf("query error:%s", apiResult.Message)
	}
	if apiResult.Data == nil {
		return nil
	}
	dataBytes, err := jsonutil.Marshal(apiResult.Data)
	if err != nil {
		return err
	}
	var files []contract.FileInfo
	err = jsonutil.Unmarshal(dataBytes, &files)
	if err != nil {
		return err
	}
	rs.syncFiles(files, serverAddr, path)
	return nil
}

func (rs *remoteClientSync) syncFiles(files []contract.FileInfo, serverAddr, path string) {
	for _, file := range files {
		if ignore.MatchPath(file.Path, "remote client sync", "sync once") {
			continue
		}
		currentPath := path + "/" + file.Path
		values := url.Values{}
		values.Add(contract.FsDir, file.IsDir.String())
		values.Add(contract.FsSize, stringutil.String(file.Size))
		values.Add(contract.FsHash, file.Hash)
		values.Add(contract.FsCtime, stringutil.String(file.CTime))
		values.Add(contract.FsAtime, stringutil.String(file.ATime))
		values.Add(contract.FsMtime, stringutil.String(file.MTime))
		syncPath := fmt.Sprintf("%s/%s?%s", serverAddr, currentPath, values.Encode())

		// create directory or file
		log.ErrorIf(rs.Create(syncPath), "sync create directory or file error => [syncPath=%s]", syncPath)

		if file.IsDir.Bool() {
			// sync current directory content
			log.ErrorIf(rs.sync(serverAddr, currentPath), "sync current directory content error => [serverAddr=%s] [currentPath=%s]", serverAddr, currentPath)
		} else {
			// sync remote file to local disk
			log.ErrorIf(rs.Write(syncPath), "sync remote file to local disk error => [syncPath=%s]", syncPath)
		}
	}
}

func (rs *remoteClientSync) buildDestAbsFile(sourceFileAbs string) (string, error) {
	remoteUrl, err := url.Parse(sourceFileAbs)
	if err != nil {
		log.Error(err, "parse url error, sourceFileAbs=%s", sourceFileAbs)
		return "", err
	}
	return filepath.Join(rs.destAbsPath, strings.TrimPrefix(remoteUrl.Path, server.SourceRoutePrefix)), nil
}

func (rs *remoteClientSync) httpGetWithAuth(rawURL string, header http.Header) (resp *http.Response, err error) {
	resp, err = httputil.HttpGetWithCookie(rawURL, header, rs.cookies...)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, os.ErrNotExist
	}
	if resp.StatusCode == http.StatusUnauthorized && rs.currentUser != nil {
		// auto login
		parseUrl, err := url.Parse(rawURL)
		if err != nil {
			return nil, err
		}
		user := rs.currentUser
		cookies, err := client.SignIn(parseUrl.Scheme, parseUrl.Host, user.UserName(), user.Password())
		if err != nil {
			return nil, err
		}
		if len(cookies) > 0 {
			rs.cookies = cookies
			log.Debug("try to auto login file server success maybe, retry to get resource => %s", rawURL)
			return httputil.HttpGetWithCookie(rawURL, header, rs.cookies...)
		}
		return nil, errors.New("file server is unauthorized")
	}
	return resp, err
}
