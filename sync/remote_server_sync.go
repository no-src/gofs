package sync

import (
	"errors"
	"fmt"
	"net"
	"path/filepath"
	"strings"
	"time"

	"github.com/no-src/gofs/action"
	"github.com/no-src/gofs/api/apiserver"
	"github.com/no-src/gofs/api/monitor"
	"github.com/no-src/gofs/contract"
	"github.com/no-src/gofs/fs"
	"github.com/no-src/gofs/server"
	"github.com/no-src/gofs/util/hashutil"
	"github.com/no-src/log"
)

var (
	errNilRemoteSyncServer = errors.New("remote sync server is nil")
	errInvalidServerPort   = errors.New("invalid server port")
)

type remoteServerSync struct {
	diskSync

	server     apiserver.Server
	serverAddr string
}

// NewRemoteServerSync create an instance of remoteServerSync execute send file change message
func NewRemoteServerSync(opt Option) (Sync, error) {
	// the fields of option
	source := opt.Source
	fileServerAddr := opt.FileServerAddr
	enableTLS := opt.EnableTLS
	certFile := opt.TLSCertFile
	keyFile := opt.TLSKeyFile
	tokenSecret := opt.TokenSecret
	users := opt.Users
	taskConf := opt.TaskConf

	ds, err := newDiskSync(opt)
	if err != nil {
		return nil, err
	}

	rs := &remoteServerSync{
		diskSync: *ds,
	}

	invalidPort := false
	fsAddr, errAddr := net.ResolveTCPAddr("tcp", fileServerAddr)
	if errAddr != nil || fsAddr.Port <= 0 {
		invalidPort = true
	}
	if len(source.FsServer()) == 0 {
		scheme := server.SchemeHttps
		if !enableTLS {
			scheme = server.SchemeHttp
		}
		if errAddr != nil {
			return nil, errAddr
		}
		if invalidPort {
			return nil, fmt.Errorf("%w => %d", errInvalidServerPort, fsAddr.Port)
		}
		rs.serverAddr = server.GenerateAddr(scheme, source.Host(), fsAddr.Port)
	} else {
		rs.serverAddr = source.FsServer()
	}
	rs.serverAddr = strings.TrimRight(rs.serverAddr, "/")
	if invalidPort {
		log.Warn("create remote server sync warning, you should enable the file server with -server and -server_addr flags")
	}

	rs.server, err = apiserver.New(source.Host(), source.Port(), enableTLS, certFile, keyFile, tokenSecret, users, opt.Reporter, rs.serverAddr, log.DefaultLogger(), taskConf)
	if err != nil {
		return nil, err
	}
	return rs, rs.start()
}

func (rs *remoteServerSync) Create(path string) error {
	if !rs.source.LocalSyncDisabled() {
		if err := rs.diskSync.Create(path); err != nil {
			return err
		}
	}
	return rs.send(action.CreateAction, path)
}

func (rs *remoteServerSync) Write(path string) error {
	if !rs.source.LocalSyncDisabled() {
		if err := rs.diskSync.Write(path); err != nil {
			return err
		}
	}
	return rs.send(action.WriteAction, path)
}

func (rs *remoteServerSync) Remove(path string) error {
	if !rs.source.LocalSyncDisabled() {
		if err := rs.diskSync.Remove(path); err != nil {
			return err
		}
	}
	return rs.send(action.RemoveAction, path)
}

func (rs *remoteServerSync) Rename(path string) error {
	if !rs.source.LocalSyncDisabled() {
		if err := rs.diskSync.Rename(path); err != nil {
			return err
		}
	}
	return rs.send(action.RenameAction, path)
}

func (rs *remoteServerSync) Chmod(path string) error {
	if !rs.source.LocalSyncDisabled() {
		if err := rs.diskSync.Chmod(path); err != nil {
			return err
		}
	}
	return rs.send(action.ChmodAction, path)
}

func (rs *remoteServerSync) send(act action.Action, path string) (err error) {
	isDir := false
	if act != action.RemoveAction && act != action.RenameAction {
		isDir, err = rs.IsDir(path)
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
	if !isDir && act == action.WriteAction {
		size, hash, hvs, err = rs.hash.GetFileSizeAndHashCheckpoints(path, rs.chunkSize, rs.checkpointCount)
		if err != nil {
			return err
		}
	}

	if act == action.WriteAction || act == action.CreateAction {
		var timeErr error
		cTime, aTime, mTime, timeErr = fs.GetFileTime(path)
		if timeErr != nil {
			return timeErr
		}
	}

	isDirValue := contract.FsNotDir
	if isDir {
		isDirValue = contract.FsIsDir
	}

	path, err = filepath.Rel(rs.sourceAbsPath, path)
	if err != nil {
		return err
	}
	path = filepath.ToSlash(path)
	req := &monitor.MonitorMessage{
		Action:  int32(act),
		BaseUrl: rs.serverAddr + server.SourceRoutePrefix,
		FileInfo: &monitor.FileInfo{
			Path:       path,
			IsDir:      int32(isDirValue),
			Size:       size,
			Hash:       hash,
			HashValues: monitor.ToHashValueMessageList(hvs),
			CTime:      cTime.Unix(),
			ATime:      aTime.Unix(),
			MTime:      mTime.Unix(),
		},
	}
	rs.server.SendMonitorMessage(req)
	return nil
}

func (rs *remoteServerSync) IsDir(path string) (bool, error) {
	return rs.diskSync.IsDir(path)
}

func (rs *remoteServerSync) SyncOnce(path string) error {
	return rs.diskSync.SyncOnce(path)
}

func (rs *remoteServerSync) start() error {
	if rs.server == nil {
		return errNilRemoteSyncServer
	}
	go log.ErrorIf(rs.server.Start(), "start api server error")
	return nil
}
