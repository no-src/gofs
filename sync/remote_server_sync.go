package sync

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/contract"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/server"
	"github.com/no-src/gofs/tran"
	"github.com/no-src/gofs/util"
	"github.com/no-src/log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type remoteServerSync struct {
	diskSync
	server     tran.Server
	serverAddr string
}

// NewRemoteServerSync create an instance of remoteServerSync execute send file change message
func NewRemoteServerSync(src, target core.VFS, bufSize int, enableTLS bool, certFile string, keyFile string, users []*auth.User, enableLogicallyDelete bool) (Sync, error) {
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

	ds := diskSync{
		srcAbsPath:    srcAbsPath,
		targetAbsPath: targetAbsPath,
		bufSize:       bufSize,
		src:           src,
		target:        target,
		baseSync:      newBaseSync(enableLogicallyDelete),
	}

	rs := &remoteServerSync{
		diskSync: ds,
	}
	rs.server = tran.NewServer(src.Host(), src.Port(), enableTLS, certFile, keyFile, users)

	if len(src.FsServer()) == 0 {
		scheme := server.SchemeHttps
		if !server.EnableTLS() {
			scheme = server.SchemeHttp
		}
		rs.serverAddr = server.GenerateAddr(scheme, rs.server.Host(), server.ServerPort())
	} else {
		rs.serverAddr = src.FsServer()
	}
	rs.serverAddr = strings.TrimRight(rs.serverAddr, "/")
	if server.ServerPort() <= 0 {
		log.Warn("create remote server sync warning, you should enable the file server with -server flag")
	}
	return rs, rs.start()
}

func (rs *remoteServerSync) Create(path string) error {
	if !rs.src.LocalSyncDisabled() {
		if err := rs.diskSync.Create(path); err != nil {
			return err
		}
	}
	return rs.send(CreateAction, path)
}

func (rs *remoteServerSync) Write(path string) error {
	if !rs.src.LocalSyncDisabled() {
		if err := rs.diskSync.Write(path); err != nil {
			return err
		}
	}
	return rs.send(WriteAction, path)
}

func (rs *remoteServerSync) Remove(path string) error {
	if !rs.src.LocalSyncDisabled() {
		if err := rs.diskSync.Remove(path); err != nil {
			return err
		}
	}
	return rs.send(RemoveAction, path)
}

func (rs *remoteServerSync) Rename(path string) error {
	if !rs.src.LocalSyncDisabled() {
		if err := rs.diskSync.Rename(path); err != nil {
			return err
		}
	}
	return rs.send(RenameAction, path)
}

func (rs *remoteServerSync) Chmod(path string) error {
	if !rs.src.LocalSyncDisabled() {
		if err := rs.diskSync.Chmod(path); err != nil {
			return err
		}
	}
	return rs.send(ChmodAction, path)
}

func (rs *remoteServerSync) send(action Action, path string) (err error) {
	isDir := false
	if action != RemoveAction && action != RenameAction {
		isDir, err = rs.IsDir(path)
		if err != nil {
			return err
		}
	}

	var size int64
	hash := ""
	cTime := time.Now()
	aTime := time.Now()
	mTime := time.Now()
	if !isDir && action == WriteAction {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		fileInfo, err := file.Stat()
		if err != nil {
			return err
		}
		size = fileInfo.Size()
		if size > 0 {
			hash, err = util.MD5FromFile(file, rs.bufSize)
			if err != nil {
				return err
			}
		}
	}

	if action == WriteAction || action == CreateAction {
		var timeErr error
		cTime, aTime, mTime, timeErr = util.GetFileTime(path)
		if timeErr != nil {
			return timeErr
		}
	}

	isDirValue := contract.FsUnknown
	if isDir {
		isDirValue = contract.FsIsDir
	} else {
		isDirValue = contract.FsNotDir
	}

	path, err = filepath.Rel(rs.srcAbsPath, path)
	path = filepath.ToSlash(path)
	req := Message{
		Status:  contract.SuccessStatus(contract.SyncMessageApi),
		Action:  action,
		BaseUrl: rs.serverAddr + server.SrcRoutePrefix,
		FileInfo: contract.FileInfo{
			Path:  path,
			IsDir: isDirValue,
			Size:  size,
			Hash:  hash,
			CTime: cTime.Unix(),
			ATime: aTime.Unix(),
			MTime: mTime.Unix(),
		},
	}

	data, err := util.Marshal(req)
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
		info = contract.FileServerInfo{
			Status:     contract.SuccessStatus(contract.InfoApi),
			ServerAddr: rs.serverAddr,
			SrcPath:    server.SrcRoutePrefix,
			TargetPath: server.TargetRoutePrefix,
			QueryAddr:  server.QueryRoute,
		}
	} else {
		info = contract.FileServerInfo{
			Status: contract.UnauthorizedStatus(contract.InfoApi),
		}
	}
	result, err = util.Marshal(info)
	return
}

func (rs *remoteServerSync) authCommand(client *tran.Conn, data []byte) (cmd contract.Command, result []byte, err error) {
	cmd = contract.AuthCommand
	authData := contract.FailStatus(contract.AuthApi)
	hashUser, err := auth.ParseAuthCommandData(data)
	if err == nil && client != nil {
		if rs.server.Auth(hashUser) {
			client.MarkAuthorized(hashUser.UserNameHash, hashUser.PasswordHash)
			authData = contract.SuccessStatus(contract.AuthApi)
		}
	} else if err != nil {
		log.Error(err, "parse auth command data error")
	}
	result, err = util.Marshal(authData)
	return
}

func (rs *remoteServerSync) unknownCommand() (cmd contract.Command, result []byte, err error) {
	cmd = contract.UnknownCommand
	respData := contract.FailStatus(contract.UnknownApi)
	respData.Message = "unknown command"
	result, err = util.Marshal(respData)
	return
}
