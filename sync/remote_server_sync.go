package sync

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/server"
	"github.com/no-src/gofs/tran"
	"github.com/no-src/gofs/util"
	"github.com/no-src/log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

type remoteServerSync struct {
	diskSync
	server tran.Server
}

func NewRemoteServerSync(src, target core.VFS, bufSize int) (Sync, error) {
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
	}

	rs := &remoteServerSync{
		diskSync: ds,
	}
	rs.server = tran.NewServer(src.Host(), src.Port())
	return rs, rs.start()
}

func (rs *remoteServerSync) Create(path string) error {
	if err := rs.diskSync.Create(path); err != nil {
		return err
	}
	return rs.send(CreateAction, path)
}

func (rs *remoteServerSync) Write(path string) error {
	if err := rs.diskSync.Write(path); err != nil {
		return err
	}
	return rs.send(WriteAction, path)
}

func (rs *remoteServerSync) Remove(path string) error {
	if err := rs.diskSync.Remove(path); err != nil {
		return err
	}
	return rs.send(RemoveAction, path)
}

func (rs *remoteServerSync) Rename(path string) error {
	if err := rs.diskSync.Rename(path); err != nil {
		return err
	}
	return rs.send(RenameAction, path)
}

func (rs *remoteServerSync) Chmod(path string) error {
	if err := rs.diskSync.Chmod(path); err != nil {
		return err
	}
	return rs.send(ChmodAction, path)
}

func (rs *remoteServerSync) send(action Action, path string) error {
	isDirValue := -1
	if action != RemoveAction && action != RenameAction {
		isDir, err := rs.IsDir(path)
		if err != nil {
			return err
		}
		if isDir {
			isDirValue = 1
		} else {
			isDirValue = 0
		}
	}

	var size int64
	hash := ""
	if isDirValue == 0 && action == WriteAction {
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
	path = strings.TrimPrefix(path, rs.srcAbsPath)
	path = filepath.ToSlash(path)
	req := Request{
		Action:  action,
		Path:    path,
		BaseUrl: rs.src.FsServer(),
		IsDir:   isDirValue,
		Size:    size,
		Hash:    hash,
	}
	
	if len(rs.src.FsServer()) == 0 {
		req.BaseUrl = fmt.Sprintf("http://%s:%d", rs.server.Host(), server.ServerPort())
	}

	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return rs.server.Send(data)
}

func (rs *remoteServerSync) IsDir(path string) (bool, error) {
	return rs.diskSync.IsDir(path)
}

func (rs *remoteServerSync) SyncOnce() error {
	return rs.diskSync.SyncOnce()
}

func (rs *remoteServerSync) start() error {
	if rs.server == nil {
		return errors.New("remote sync server is nil")
	}

	err := rs.server.Listen()
	if err != nil {
		return err
	}
	go rs.server.Accept(func(client net.Conn, data []byte) {
		log.Debug("receive message [%s] => %s", client.RemoteAddr().String(), string(data))
		//writer := bufio.NewWriter(client)
		//result := append(data, tran.EndIdentity...)
		//result = append(result, tran.LFBytes...)
		//_, err = writer.Write(result)
		//writer.Flush()
	})
	return nil
}
