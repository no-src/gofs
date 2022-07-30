package sync

import (
	"bufio"
	"errors"
	iofs "io/fs"
	"os"

	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/driver"
	"github.com/no-src/gofs/driver/sftp"
	nsfs "github.com/no-src/gofs/fs"
	"github.com/no-src/gofs/ignore"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/log"
)

type sftpPullClientSync struct {
	diskSync

	remoteAddr  string
	remotePath  string
	currentUser *auth.User
	client      driver.Driver
}

// NewSftpPullClientSync create an instance of the sftpPullClientSync
func NewSftpPullClientSync(source, dest core.VFS, users []*auth.User, enableLogicallyDelete bool, chunkSize int64, checkpointCount int, forceChecksum bool, r retry.Retry) (Sync, error) {
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

	s := &sftpPullClientSync{
		diskSync:    *ds,
		remoteAddr:  source.Addr(),
		remotePath:  source.RemotePath(),
		currentUser: users[0],
	}

	s.client = sftp.NewSFTPClient(s.remoteAddr, s.currentUser.UserName(), s.currentUser.Password(), true, r)

	err = s.start()
	if err != nil {
		return nil, err
	}

	// reset the sourceAbsPath because the source.Path() or source.RemotePath() is absolute representation of path and the source.Path() or source.RemotePath() may be cross-platform
	// source.Path() and source.RemotePath() are equivalent here, and source.RemotePath() has higher priority
	s.diskSync.sourceAbsPath = source.RemotePath()
	if len(s.diskSync.sourceAbsPath) == 0 {
		s.diskSync.sourceAbsPath = source.Path()
	}
	// reset some functions for sftp
	s.diskSync.isDirFn = s.IsDir
	s.diskSync.statFn = s.client.Stat
	s.diskSync.getFileTimeFn = s.client.GetFileTime

	return s, nil
}

func (s *sftpPullClientSync) start() error {
	return s.client.Connect()
}

func (s *sftpPullClientSync) Create(path string) error {
	return s.diskSync.Create(path)
}

func (s *sftpPullClientSync) Write(path string) error {
	dest, err := s.buildDestAbsFile(path)
	if err != nil {
		return err
	}

	isDir, err := s.IsDir(path)
	if err != nil {
		return err
	}

	// process directory
	if isDir {
		return s.SyncOnce(path)
	}

	// process file
	return s.write(path, dest)
}

// write try to write a file to the destination
func (s *sftpPullClientSync) write(path, dest string) error {
	sourceFile, err := s.client.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		log.ErrorIf(sourceFile.Close(), "Write:close the source file error")
	}()

	sourceStat, err := sourceFile.Stat()
	if err != nil {
		return err
	}

	destStat, err := os.Stat(dest)
	if err != nil {
		return err
	}

	sourceSize := sourceStat.Size()
	destSize := destStat.Size()

	if s.quickCompare(sourceSize, destSize, sourceStat.ModTime(), destStat.ModTime()) {
		log.Debug("Write:ignored, the file size and file modification time are both unmodified => %s", path)
		return nil
	}

	destFile, err := nsfs.OpenRWFile(dest)
	if err != nil {
		return err
	}
	defer func() {
		log.ErrorIf(destFile.Close(), "Write:close the dest file error")
	}()

	reader := bufio.NewReader(sourceFile)
	writer := bufio.NewWriter(destFile)

	// truncate first before write to file
	err = destFile.Truncate(0)
	if err != nil {
		return err
	}

	n, err := reader.WriteTo(writer)
	if err != nil {
		return err
	}

	err = writer.Flush()

	if err == nil {
		log.Info("write to the dest file success, size[%d => %d] [%s] => [%s]", sourceSize, n, path, dest)
		s.chtimes(path, dest)
	}
	return err
}

func (s *sftpPullClientSync) Remove(path string) error {
	return s.diskSync.Remove(path)
}

func (s *sftpPullClientSync) Rename(path string) error {
	return s.diskSync.Rename(path)
}

func (s *sftpPullClientSync) Chmod(path string) error {
	log.Debug("Chmod is unimplemented [%s]", path)
	return nil
}

func (s *sftpPullClientSync) IsDir(path string) (bool, error) {
	fi, err := s.client.Stat(path)
	if err != nil {
		return false, err
	}
	return fi.IsDir(), nil
}

func (s *sftpPullClientSync) SyncOnce(path string) error {
	return s.client.WalkDir(path, func(currentPath string, d iofs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if ignore.MatchPath(currentPath, "sftp pull client sync", "sync once") {
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
