package sync

import (
	"bufio"
	"bytes"
	"errors"
	"path/filepath"
	"strings"
	"time"

	"github.com/no-src/gofs/action"
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/contract"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/fs"
	"github.com/no-src/gofs/server"
	"github.com/no-src/gofs/tran"
	"github.com/no-src/gofs/util/hashutil"
	"github.com/no-src/gofs/util/jsonutil"
	"github.com/no-src/log"
)

type remoteServerSync struct {
	diskSync

	server          tran.Server
	serverAddr      string
	chunkSize       int64
	checkpointCount int
}

// NewRemoteServerSync create an instance of remoteServerSync execute send file change message
func NewRemoteServerSync(source, dest core.VFS, enableTLS bool, certFile string, keyFile string, users []*auth.User, enableLogicallyDelete bool, chunkSize int64, checkpointCount int) (Sync, error) {
	ds, err := newDiskSync(source, dest, enableLogicallyDelete, chunkSize, checkpointCount)
	if err != nil {
		return nil, err
	}

	rs := &remoteServerSync{
		diskSync:        *ds,
		chunkSize:       chunkSize,
		checkpointCount: checkpointCount,
	}
	rs.server = tran.NewServer(source.Host(), source.Port(), enableTLS, certFile, keyFile, users)

	if len(source.FsServer()) == 0 {
		scheme := server.SchemeHttps
		if !server.EnableTLS() {
			scheme = server.SchemeHttp
		}
		rs.serverAddr = server.GenerateAddr(scheme, rs.server.Host(), server.Port())
	} else {
		rs.serverAddr = source.FsServer()
	}
	rs.serverAddr = strings.TrimRight(rs.serverAddr, "/")
	if server.Port() <= 0 {
		log.Warn("create remote server sync warning, you should enable the file server with -server flag")
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
		size, hash, hvs, err = rs.getFileSizeAndHashCheckpoints(path, rs.chunkSize, rs.checkpointCount)
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
	path = filepath.ToSlash(path)
	req := Message{
		Status:  contract.SuccessStatus(contract.SyncMessageApi),
		Action:  act,
		BaseUrl: rs.serverAddr + server.SourceRoutePrefix,
		FileInfo: contract.FileInfo{
			Path:       path,
			IsDir:      isDirValue,
			Size:       size,
			Hash:       hash,
			HashValues: hvs,
			CTime:      cTime.Unix(),
			ATime:      aTime.Unix(),
			MTime:      mTime.Unix(),
		},
	}

	data, err := jsonutil.Marshal(req)
	if err != nil {
		return err
	}
	return rs.server.Send(data)
}

func (rs *remoteServerSync) IsDir(path string) (bool, error) {
	return rs.diskSync.IsDir(path)
}

func (rs *remoteServerSync) SyncOnce(path string) error {
	return rs.diskSync.SyncOnce(path)
}

func (rs *remoteServerSync) start() error {
	if rs.server == nil {
		return errors.New("remote sync server is nil")
	}

	err := rs.server.Listen()
	if err != nil {
		return err
	}
	go rs.server.Accept(func(client *tran.Conn, data []byte) {
		if bytes.HasSuffix(data, tran.EndIdentity) {
			data = data[:len(data)-len(tran.EndIdentity)]
		}
		if client == nil {
			log.Warn("client conn is nil, data => %s", string(data))
			return
		}
		log.Debug("receive message [%s] => %s", client.RemoteAddr().String(), string(data))
		writer := bufio.NewWriter(client)
		var result []byte
		var cmd contract.Command
		if bytes.Equal(data, contract.InfoCommand) {
			cmd, result, err = rs.infoCommand(client)
		} else if bytes.HasPrefix(data, contract.AuthCommand) {
			cmd, result, err = rs.authCommand(client, data)
		} else {
			cmd, result, err = rs.unknownCommand()
		}

		// write to response
		if err != nil {
			result = append(result, []byte(err.Error())...)
			result = append(result, tran.ErrorEndIdentity...)
		} else {
			result = append(result, tran.EndIdentity...)
		}
		result = append(result, tran.LFBytes...)
		_, err = writer.Write(result)
		if err != nil {
			log.Error(err, "[%s]=>[%s] write message error", client.RemoteAddr().String(), string(cmd))
		}
		err = writer.Flush()
		if err != nil {
			log.Error(err, "[%s]=>[%s] flush message error", client.RemoteAddr().String(), string(cmd))
		}
	})
	return nil
}

func (rs *remoteServerSync) infoCommand(client *tran.Conn) (cmd contract.Command, result []byte, err error) {
	cmd = contract.InfoCommand
	var info contract.FileServerInfo
	if client.Authorized() {
		if client.CheckPerm(auth.ReadPerm) {
			info = contract.FileServerInfo{
				Status:     contract.SuccessStatus(contract.InfoApi),
				ServerAddr: rs.serverAddr,
				SourcePath: server.SourceRoutePrefix,
				DestPath:   server.DestRoutePrefix,
				QueryAddr:  server.QueryRoute,
				PushAddr:   server.PushFullRoute,
			}
		} else {
			info = contract.FileServerInfo{
				Status: contract.NoPermissionStatus(contract.InfoApi),
			}
		}
	} else {
		info = contract.FileServerInfo{
			Status: contract.UnauthorizedStatus(contract.InfoApi),
		}
	}
	result, err = jsonutil.Marshal(info)
	return
}

func (rs *remoteServerSync) authCommand(client *tran.Conn, data []byte) (cmd contract.Command, result []byte, err error) {
	cmd = contract.AuthCommand
	authData := contract.FailStatus(contract.AuthApi)
	hashUser, err := auth.ParseAuthCommandData(data)
	if err == nil && client != nil {
		authed, perm := rs.server.Auth(hashUser)
		if authed {
			hashUser.Perm = perm
			client.MarkAuthorized(hashUser)
			if auth.ToPerm(auth.ReadPerm).CheckTo(hashUser.Perm) {
				authData = contract.SuccessStatus(contract.AuthApi)
			} else {
				authData = contract.NewStatus(contract.Success, "warning: you are authorized but have no permission to read", contract.AuthApi)
			}
		}
	} else if err != nil {
		log.Error(err, "parse auth command data error")
	}
	result, err = jsonutil.Marshal(authData)
	return
}

func (rs *remoteServerSync) unknownCommand() (cmd contract.Command, result []byte, err error) {
	cmd = contract.UnknownCommand
	respData := contract.FailStatus(contract.UnknownApi)
	respData.Message = "unknown command"
	result, err = jsonutil.Marshal(respData)
	return
}
