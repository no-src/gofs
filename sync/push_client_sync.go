package sync

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/no-src/gofs/action"
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/contract"
	"github.com/no-src/gofs/contract/push"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/encrypt"
	nsfs "github.com/no-src/gofs/fs"
	"github.com/no-src/gofs/ignore"
	"github.com/no-src/gofs/server"
	"github.com/no-src/gofs/server/client"
	"github.com/no-src/gofs/tran"
	"github.com/no-src/gofs/util/hashutil"
	"github.com/no-src/gofs/util/httputil"
	"github.com/no-src/gofs/util/jsonutil"
	"github.com/no-src/log"
)

type pushClientSync struct {
	diskSync

	pushAddr        string
	cookies         []*http.Cookie
	currentUser     *auth.User
	currentHashUser *auth.HashUser
	client          tran.Client
	authChan        chan contract.Status
	infoChan        chan contract.Message
	timeout         time.Duration
}

// NewPushClientSync create an instance of the pushClientSync
func NewPushClientSync(source, dest core.VFS, enableTLS bool, certFile string, insecureSkipVerify bool, users []*auth.User, enableLogicallyDelete bool, chunkSize int64, checkpointCount int, forceChecksum bool) (Sync, error) {
	if chunkSize <= 0 {
		return nil, errors.New("chunk size must greater than zero")
	}

	ds, err := newDiskSync(source, dest, enableLogicallyDelete, chunkSize, checkpointCount, forceChecksum, encrypt.EmptyOption())
	if err != nil {
		return nil, err
	}

	s := &pushClientSync{
		diskSync: *ds,
		client:   tran.NewClient(dest.Host(), dest.Port(), enableTLS, certFile, insecureSkipVerify),
		authChan: make(chan contract.Status, 100),
		infoChan: make(chan contract.Message, 100),
		timeout:  time.Minute * 3,
	}

	if len(users) > 0 {
		user := users[0]
		hashUser := user.ToHashUser()
		s.currentUser = user
		s.currentHashUser = hashUser
	}

	err = s.start()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (pcs *pushClientSync) start() error {
	err := pcs.client.Connect()
	if err != nil {
		return err
	}
	pcs.receive()
	err = pcs.auth()
	if err != nil {
		return err
	}
	err = pcs.info()
	if err == nil {
		pcs.client.Close()
	}
	return err
}

func (pcs *pushClientSync) auth() error {
	// if the currentHashUser is nil, it means to anonymous access
	if pcs.currentHashUser == nil {
		return nil
	}
	go func() {
		pcs.currentHashUser.RefreshExpires()
		authData := auth.GenerateAuthCommandData(pcs.currentHashUser)
		log.ErrorIf(pcs.client.Write(authData), "send auth request error")
	}()

	var status contract.Status
	select {
	case status = <-pcs.authChan:
	case <-time.After(pcs.timeout):
		return fmt.Errorf("auth timeout for %s", pcs.timeout.String())
	}
	if status.Code != contract.Success {
		return errors.New("receive auth command response error => " + status.Message)
	}

	log.Info("auth success, current client is authorized => [%s] ", status.Message)
	return nil
}

func (pcs *pushClientSync) info() error {
	go func() {
		log.ErrorIf(pcs.client.Write(contract.InfoCommand), "write info command error")
	}()
	var info contract.FileServerInfo
	var infoMsg contract.Message
	select {
	case infoMsg = <-pcs.infoChan:
	case <-time.After(pcs.timeout):
		return fmt.Errorf("info timeout for %s", pcs.timeout.String())
	}
	err := jsonutil.Unmarshal(infoMsg.Data, &info)
	if err != nil {
		return err
	}

	if info.Code != contract.Success {
		return errors.New("receive info command response error => " + info.Message)
	}
	pcs.pushAddr = info.ServerAddr + info.PushAddr
	return nil
}

func (pcs *pushClientSync) receive() {
	go func() {
		for {
			if pcs.client.IsClosed() {
				break
			}
			data, err := pcs.client.ReadAll()
			if err != nil {
				if pcs.client.IsClosed() {
					break
				} else {
					log.Error(err, "[push client sync] read data error")
				}
			} else {
				var status contract.Status
				err = jsonutil.Unmarshal(data, &status)
				if err != nil {
					log.Error(err, "[push client sync] unmarshal data error")
					continue
				}
				switch status.ApiType {
				case contract.AuthApi:
					pcs.authChan <- status
				case contract.InfoApi:
					pcs.infoChan <- contract.NewMessage(data)
				default:
					log.Warn("[push client sync] receive and discard data => %s", string(data))
				}
			}
		}
	}()
}

func (pcs *pushClientSync) Create(path string) error {
	if !pcs.dest.LocalSyncDisabled() {
		if err := pcs.diskSync.Create(path); err != nil {
			return err
		}
	}
	return pcs.send(action.CreateAction, path)
}

func (pcs *pushClientSync) Write(path string) error {
	if !pcs.dest.LocalSyncDisabled() {
		if err := pcs.diskSync.Write(path); err != nil {
			return err
		}
	}
	isDir, err := pcs.IsDir(path)
	if err != nil {
		return err
	}
	if isDir {
		return pcs.SyncOnce(path)
	}
	return pcs.send(action.WriteAction, path)
}

func (pcs *pushClientSync) Remove(path string) error {
	if !pcs.dest.LocalSyncDisabled() {
		if err := pcs.diskSync.Remove(path); err != nil {
			return err
		}
	}
	return pcs.send(action.RemoveAction, path)
}

func (pcs *pushClientSync) Rename(path string) error {
	if !pcs.dest.LocalSyncDisabled() {
		if err := pcs.diskSync.Remove(path); err != nil {
			return err
		}
	}
	return pcs.send(action.RenameAction, path)
}

func (pcs *pushClientSync) Chmod(path string) error {
	log.Debug("Chmod is unimplemented [%s]", path)
	return nil
}

func (pcs *pushClientSync) IsDir(path string) (bool, error) {
	return pcs.diskSync.IsDir(path)
}

func (pcs *pushClientSync) SyncOnce(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	return filepath.WalkDir(absPath, func(currentPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if ignore.MatchPath(currentPath, "push client sync", "sync once") {
			return nil
		}
		if d.IsDir() {
			err = pcs.Create(currentPath)
		} else {
			err = pcs.Create(currentPath)
			if err == nil {
				err = pcs.Write(currentPath)
			}
		}
		return err
	})
}

func (pcs *pushClientSync) send(act action.Action, path string) (err error) {
	isDir := false
	if pcs.needCheckDir(act) {
		isDir, err = pcs.IsDir(path)
		if err != nil {
			return err
		}
	}

	var size int64
	hash := ""
	var hvs hashutil.HashValues
	cTime := time.Now()
	aTime := time.Now()
	mTime := time.Now()
	if pcs.needGetFileSizeAndHash(isDir, act) {
		size, hash, hvs, err = pcs.getFileSizeAndHashCheckpoints(path, pcs.chunkSize, pcs.checkpointCount)
		if err != nil {
			return err
		}
	} else if pcs.needIgnoreDirWrite(isDir, act) {
		return nil
	}

	if pcs.needGetFileTime(act) {
		var timeErr error
		cTime, aTime, mTime, timeErr = nsfs.GetFileTime(path)
		if timeErr != nil {
			return timeErr
		}
	}

	isDirValue := contract.FsNotDir
	if isDir {
		isDirValue = contract.FsIsDir
	}

	relPath, err := filepath.Rel(pcs.sourceAbsPath, path)
	if err != nil {
		return err
	}
	relPath = filepath.ToSlash(relPath)
	pd := push.PushData{
		Action: act,
		FileInfo: contract.FileInfo{
			Path:       relPath,
			IsDir:      isDirValue,
			Size:       size,
			Hash:       hash,
			HashValues: hvs,
			CTime:      cTime.Unix(),
			ATime:      aTime.Unix(),
			MTime:      mTime.Unix(),
		},
		ForceChecksum: pcs.forceChecksum,
	}
	return pcs.sendPushData(pd, act, path)
}

func (pcs *pushClientSync) needCheckDir(act action.Action) bool {
	return act != action.RemoveAction && act != action.RenameAction
}

func (pcs *pushClientSync) needGetFileSizeAndHash(isDir bool, act action.Action) bool {
	return !isDir && act == action.WriteAction
}

func (pcs *pushClientSync) needIgnoreDirWrite(isDir bool, act action.Action) bool {
	return isDir && act == action.WriteAction
}

func (pcs *pushClientSync) needGetFileTime(act action.Action) bool {
	return act == action.WriteAction || act == action.CreateAction
}

func (pcs *pushClientSync) sendPushData(pd push.PushData, act action.Action, path string) error {
	if act == action.WriteAction {
		return pcs.sendFileChunk(path, pd)
	}
	resp, err := pcs.httpPostWithAuth(pcs.pushAddr, act, push.ParamUpFile, path, pd, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, _, err = pcs.checkApiResult(resp)
	return err
}

func (pcs *pushClientSync) sendFileChunk(path string, pd push.PushData) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	var offset int64
	chunk := make([]byte, pcs.chunkSize)
	isEnd := false

	// if loopCount == 0 means is a request that check file size and hash value only, send it
	// if loopCount == 1 means read an empty file maybe, send it
	loopCount := -1
	checkChunkHash := false
	for {
		loopCount++
		n, err := f.ReadAt(chunk, offset)
		if nsfs.IsNonEOF(err) {
			return err
		}
		if nsfs.IsEOF(err) {
			isEnd = true
		}
		chunkSize := n
		pd.PushAction = push.WritePushAction
		if pcs.needCheckHash(loopCount, chunkSize) {
			// if file size is greater than zero in first loop, try to compare file and chunk hash
			isEnd = false
			n = 0
			pd.PushAction = push.CompareFileAndChunkPushAction
			checkChunkHash = true
		} else if pcs.isReadEmptyFile(loopCount, chunkSize) {
			// if read data nothing, send the empty file and cancel file compare request
			loopCount++
		} else if checkChunkHash {
			pd.PushAction = push.CompareChunkPushAction
			n = 0
		}

		if pcs.needSendChunkRequest(loopCount, chunkSize) {
			broken, err := pcs.sendChunkRequest(path, &pd, &offset, chunkSize, &checkChunkHash, chunk, n, &isEnd)
			if broken {
				return err
			}
		}
		if isEnd {
			// read to end, send a truncate request finally
			return pcs.sendTruncate(path, pd, offset)
		}
	}
}

func (pcs *pushClientSync) sendChunkRequest(path string, pd *push.PushData, offset *int64, chunkSize int, checkChunkHash *bool, chunk []byte, n int, isEnd *bool) (broken bool, err error) {
	defer func() {
		// only send HashValues once
		if len(pd.FileInfo.HashValues) > 0 {
			pd.FileInfo.HashValues = nil
		}
	}()
	pd.Chunk.Offset = *offset
	pd.Chunk.Size = int64(chunkSize)
	if *checkChunkHash {
		pd.Chunk.Hash = hashutil.Hash(chunk[:chunkSize])
	}

	resp, err := pcs.httpPostWithAuth(pcs.pushAddr, action.WriteAction, push.ParamUpFile, path, *pd, chunk[:n])
	if err != nil {
		return true, err
	}
	code, hv, err := pcs.checkApiResult(resp)
	resp.Body.Close()
	if err != nil {
		return true, err
	} else if code == contract.NotModified {
		log.Debug("upload a file that not modified, ignore and abort next request => %s", path)
		return true, nil
	} else if code == contract.ChunkNotModified {
		// current chunk is not modified, continue to compare next chunk in the next loop
		log.Debug("upload a file chunk that not modified, continue to compare next chunk [%d]=> %s", *offset, path)
		*checkChunkHash = true
		*offset += int64(chunkSize)
		// if the checkpoint compare result offset is greater than the next offset, then replace it
		if hv != nil && hv.Offset > *offset {
			*offset = hv.Offset
		}
	} else if code == contract.Modified {
		// file is modified and the first chunk is modified too, upload the file in the next loop
		*offset = 0
		*checkChunkHash = false
		*isEnd = false
	} else if code == contract.ChunkModified {
		// write current chunk in the next loop
		*checkChunkHash = false
		*isEnd = false
	} else {
		// get success code, continue to write next chunk in the next loop or send a truncate request in the end
		*offset += int64(chunkSize)
	}
	return false, nil
}

func (pcs *pushClientSync) sendTruncate(path string, pd push.PushData, offset int64) error {
	pd.PushAction = push.TruncatePushAction
	pd.Chunk.Offset = offset
	resp, err := pcs.httpPostWithAuth(pcs.pushAddr, action.WriteAction, push.ParamUpFile, path, pd, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, _, err = pcs.checkApiResult(resp)
	return err
}

func (pcs *pushClientSync) needCheckHash(loopCount, dataLen int) bool {
	return loopCount == 0 && dataLen > 0
}

func (pcs *pushClientSync) isReadEmptyFile(loopCount, dataLen int) bool {
	return loopCount == 0 && dataLen <= 0
}

func (pcs *pushClientSync) needSendChunkRequest(loopCount, dataLen int) bool {
	return dataLen > 0 || loopCount <= 1
}

func (pcs *pushClientSync) checkApiResult(resp *http.Response) (code contract.Code, hv *hashutil.HashValue, err error) {
	var apiResult server.ApiResult
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return code, nil, err
	}
	err = jsonutil.Unmarshal(respData, &apiResult)
	if err != nil {
		return code, nil, err
	}

	if apiResult.Data != nil {
		dataBytes, err := jsonutil.Marshal(apiResult.Data)
		if err != nil {
			return code, nil, err
		}
		err = jsonutil.Unmarshal(dataBytes, &hv)
		if err != nil {
			return code, nil, err
		}
	}

	code = apiResult.Code
	switch code {
	case contract.NotModified, contract.ChunkNotModified, contract.Modified, contract.ChunkModified:
		return code, hv, nil
	}

	if code != contract.Success {
		err = fmt.Errorf("send a request to the push server error => %s", apiResult.Message)
	}
	return code, hv, err
}

func (pcs *pushClientSync) httpPostWithAuth(rawURL string, act action.Action, fieldName string, fileName string, pd push.PushData, chunk []byte) (resp *http.Response, err error) {
	pdData, err := jsonutil.Marshal(pd)
	if err != nil {
		return nil, err
	}
	data := url.Values{}
	data.Set(push.ParamPushData, string(pdData))

	sendFile := false
	if act == action.WriteAction {
		sendFile = true
	}
	if sendFile {
		resp, err = httputil.HttpPostFileChunkWithCookie(rawURL, fieldName, fileName, data, chunk, pcs.cookies...)
	} else {
		resp, err = httputil.HttpPostWithCookie(rawURL, data, pcs.cookies...)
	}

	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusUnauthorized && pcs.currentUser != nil {
		// auto login
		parseUrl, err := url.Parse(rawURL)
		if err != nil {
			return nil, err
		}
		user := pcs.currentUser
		cookies, err := client.SignIn(parseUrl.Scheme, parseUrl.Host, user.UserName(), user.Password())
		if err != nil {
			return nil, err
		}
		if len(cookies) > 0 {
			pcs.cookies = cookies
			log.Debug("try to auto login file server success maybe, retry to get resource => %s", rawURL)
			if sendFile {
				return httputil.HttpPostFileChunkWithCookie(rawURL, fieldName, fileName, data, chunk, pcs.cookies...)
			}
			return httputil.HttpPostWithCookie(rawURL, data, pcs.cookies...)
		}
		return nil, errors.New("file server is unauthorized")
	} else if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("the push server is unsupported => %s", rawURL)
	}
	return resp, err
}
