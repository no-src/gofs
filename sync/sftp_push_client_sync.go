package sync

import (
	"errors"
	iofs "io/fs"
	"os"
	"path/filepath"

	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/driver"
	"github.com/no-src/gofs/driver/sftp"
	"github.com/no-src/gofs/fs"
	"github.com/no-src/gofs/ignore"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/log"
)

type sftpPushClientSync struct {
	diskSync

	remoteAddr  string
	remotePath  string
	currentUser *auth.User
	client      driver.Driver
}

// NewSftpPushClientSync create an instance of the sftpPushClientSync
func NewSftpPushClientSync(source, dest core.VFS, users []*auth.User, enableLogicallyDelete bool, chunkSize int64, checkpointCount int, forceChecksum bool, r retry.Retry) (Sync, error) {
	if chunkSize <= 0 {
		return nil, errors.New("chunk size must greater than zero")
	}

	if len(users) == 0 {
		return nil, errors.New("user account is required")
	}

	ds, err := newDiskSync(source, dest, enableLogicallyDelete, chunkSize, checkpointCount, forceChecksum)
	if err != nil {
		return nil, err
	}

	s := &sftpPushClientSync{
		diskSync:    *ds,
		remoteAddr:  dest.Addr(),
		remotePath:  dest.RemotePath(),
		currentUser: users[0],
	}

	s.client = sftp.NewSFTPClient(s.remoteAddr, s.currentUser.UserName(), s.currentUser.Password(), true, r)

	err = s.start()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *sftpPushClientSync) start() error {
	return s.client.Connect()
}

func (s *sftpPushClientSync) Create(path string) error {
	if !s.dest.LocalSyncDisabled() {
		if err := s.diskSync.Create(path); err != nil {
			return err
		}
	}

	destPath, err := s.buildDestAbsFile(path)
	if err != nil {
		return err
	}
	isDir, err := s.IsDir(path)
	if err != nil {
		return err
	}
	if isDir {
		err = s.client.MkdirAll(destPath)
	} else {
		err = s.client.Create(destPath)
	}
	return err
}

func (s *sftpPushClientSync) Write(path string) error {
	if !s.dest.LocalSyncDisabled() {
		if err := s.diskSync.Write(path); err != nil {
			return err
		}
	}
	isDir, err := s.IsDir(path)
	if err != nil {
		return err
	}
	if isDir {
		return s.SyncOnce(path)
	}

	// write to sftp server
	destPath, err := s.buildDestAbsFile(path)
	if err != nil {
		return err
	}

	err = s.client.Write(path, destPath)
	if err == nil {
		if _, aTime, mTime, err := fs.GetFileTime(path); err == nil {
			log.ErrorIf(s.client.Chtimes(destPath, aTime, mTime), "[sftp push client sync] [write] change file times error")
		}
	}
	return err
}

func (s *sftpPushClientSync) Remove(path string) error {
	if !s.dest.LocalSyncDisabled() {
		if err := s.diskSync.Remove(path); err != nil {
			return err
		}
	}

	return s.remove(path, false)
}

func (s *sftpPushClientSync) remove(path string, forceDelete bool) (err error) {
	destPath, err := s.buildDestAbsFile(path)
	if err != nil {
		return err
	}
	if !forceDelete && s.enableLogicallyDelete {
		err = s.logicallyDelete(destPath)
	} else {
		err = s.client.Remove(destPath)
	}
	return err
}

// logicallyDelete delete the path logically
func (s *sftpPushClientSync) logicallyDelete(path string) error {
	if fs.IsDeleted(path) {
		return nil
	}
	deletedFile := fs.ToDeletedPath(path)
	err := s.client.Rename(path, deletedFile)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (s *sftpPushClientSync) Rename(path string) error {
	if !s.dest.LocalSyncDisabled() {
		if err := s.diskSync.Remove(path); err != nil {
			return err
		}
	}

	return s.remove(path, true)
}

func (s *sftpPushClientSync) Chmod(path string) error {
	log.Debug("Chmod is unimplemented [%s]", path)
	return nil
}

func (s *sftpPushClientSync) IsDir(path string) (bool, error) {
	return s.diskSync.IsDir(path)
}

func (s *sftpPushClientSync) SyncOnce(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	return filepath.WalkDir(absPath, func(currentPath string, d iofs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if ignore.MatchPath(currentPath, "sftp push client sync", "sync once") {
			return nil
		}
		if d.IsDir() {
			err = s.Create(currentPath)
		} else {
			err = s.Create(currentPath)
			if err == nil {
				err = s.Write(currentPath)
			}
		}
		return err
	})
}

func (s *sftpPushClientSync) buildDestAbsFile(sourceFileAbs string) (string, error) {
	sourceFileRel, err := filepath.Rel(s.sourceAbsPath, sourceFileAbs)
	if err != nil {
		log.Error(err, "parse rel path error, basePath=%s destPath=%s", s.sourceAbsPath, sourceFileRel)
		return "", err
	}

	return filepath.ToSlash(filepath.Join(s.remotePath, sourceFileRel)), nil
}
