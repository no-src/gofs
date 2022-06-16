package sync

import (
	"bufio"
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
	"github.com/no-src/log"
)

type sftpClientSync struct {
	diskSync

	remoteAddr  string
	remotePath  string
	currentUser *auth.User
	client      driver.Driver
}

// NewSftpClientSync create an instance of the sftpClientSync
func NewSftpClientSync(source, dest core.VFS, users []*auth.User, enableLogicallyDelete bool, chunkSize int64, checkpointCount int, forceChecksum bool) (Sync, error) {
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

	s := &sftpClientSync{
		diskSync:    *ds,
		remoteAddr:  dest.Addr(),
		remotePath:  dest.RemotePath(),
		currentUser: users[0],
	}

	s.client = sftp.NewSFTPClient(s.remoteAddr, s.currentUser.UserName(), s.currentUser.Password(), true)

	err = s.start()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *sftpClientSync) start() error {
	return s.client.Connect()
}

func (s *sftpClientSync) Create(path string) error {
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
		_, err = s.client.Create(destPath)
	}
	return err
}

func (s *sftpClientSync) Write(path string) error {
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
	sourceFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		log.ErrorIf(sourceFile.Close(), "sftp write: close the source file error")
	}()

	reader := bufio.NewReader(sourceFile)

	destFile, err := s.client.Create(destPath)
	if err != nil {
		return err
	}
	defer func() {
		log.ErrorIf(destFile.Close(), "sftp write: close the dest file error")
	}()
	_, err = reader.WriteTo(destFile)
	if err == nil {
		if _, aTime, mTime, err := fs.GetFileTime(path); err == nil {
			log.ErrorIf(s.client.Chtimes(destPath, aTime, mTime), "sftp change file times error")
		}
	}
	return err
}

func (s *sftpClientSync) Remove(path string) error {
	if !s.dest.LocalSyncDisabled() {
		if err := s.diskSync.Remove(path); err != nil {
			return err
		}
	}

	return s.remove(path, false)
}

func (s *sftpClientSync) remove(path string, forceDelete bool) (err error) {
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
func (s *sftpClientSync) logicallyDelete(path string) error {
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

func (s *sftpClientSync) Rename(path string) error {
	if !s.dest.LocalSyncDisabled() {
		if err := s.diskSync.Remove(path); err != nil {
			return err
		}
	}

	return s.remove(path, true)
}

func (s *sftpClientSync) Chmod(path string) error {
	log.Debug("Chmod is unimplemented [%s]", path)
	return nil
}

func (s *sftpClientSync) IsDir(path string) (bool, error) {
	return s.diskSync.IsDir(path)
}

func (s *sftpClientSync) SyncOnce(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	return filepath.WalkDir(absPath, func(currentPath string, d iofs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if ignore.MatchPath(currentPath, "sftp client sync", "sync once") {
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

func (s *sftpClientSync) buildDestAbsFile(sourceFileAbs string) (string, error) {
	sourceFileRel, err := filepath.Rel(s.sourceAbsPath, sourceFileAbs)
	if err != nil {
		log.Error(err, "parse rel path error, basePath=%s destPath=%s", s.sourceAbsPath, sourceFileRel)
		return "", err
	}

	return filepath.ToSlash(filepath.Join(s.remotePath, sourceFileRel)), nil
}
