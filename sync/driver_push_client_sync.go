package sync

import (
	"errors"
	iofs "io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/no-src/gofs/contract"
	"github.com/no-src/gofs/driver"
	"github.com/no-src/gofs/fs"
	"github.com/no-src/gofs/ignore"
	"github.com/no-src/log"
)

type driverPushClientSync struct {
	diskSync

	basePath string
	client   driver.Driver
	files    sync.Map
}

func (s *driverPushClientSync) start() error {
	if err := s.initFileInfo(); err != nil {
		return err
	}
	return s.client.Connect()
}

func (s *driverPushClientSync) Create(path string) error {
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

func (s *driverPushClientSync) Write(path string) error {
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

	// compare whether the file is changed or not
	if s.fileInfoCompare(path) {
		log.Debug("[push] [ignored], the file modification time is unmodified => %s", path)
		return nil
	}

	// write to sftp server
	destPath, err := s.buildDestAbsFile(path)
	if err != nil {
		return err
	}

	encryptPath, removeTemp, err := s.enc.CreateEncryptTemp(path)
	if err != nil {
		return err
	}
	// remove the temporary file
	defer removeTemp()

	err = s.client.Write(encryptPath, destPath)
	if err == nil {
		log.Debug("[push] [success] => %s", path)
		if _, aTime, mTime, err := fs.GetFileTime(path); err == nil {
			log.ErrorIf(s.client.Chtimes(destPath, aTime, mTime), "[%s push client sync] [write] change file times error", s.client.DriverName())
		}
		s.storeFileInfo(path)
	}
	return err
}

func (s *driverPushClientSync) Remove(path string) error {
	if !s.dest.LocalSyncDisabled() {
		if err := s.diskSync.Remove(path); err != nil {
			return err
		}
	}

	return s.remove(path, false)
}

func (s *driverPushClientSync) remove(path string, forceDelete bool) (err error) {
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
func (s *driverPushClientSync) logicallyDelete(path string) error {
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

func (s *driverPushClientSync) Rename(path string) error {
	if !s.dest.LocalSyncDisabled() {
		if err := s.diskSync.Remove(path); err != nil {
			return err
		}
	}

	return s.remove(path, true)
}

func (s *driverPushClientSync) Chmod(path string) error {
	log.Debug("Chmod is unimplemented [%s]", path)
	return nil
}

func (s *driverPushClientSync) IsDir(path string) (bool, error) {
	return s.diskSync.IsDir(path)
}

func (s *driverPushClientSync) SyncOnce(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	return filepath.WalkDir(absPath, func(currentPath string, d iofs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if ignore.MatchPath(currentPath, s.client.DriverName()+" push client sync", "sync once") {
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

func (s *driverPushClientSync) buildDestAbsFile(sourceFileAbs string) (string, error) {
	sourceFileRel, err := filepath.Rel(s.sourceAbsPath, sourceFileAbs)
	if err != nil {
		log.Error(err, "parse rel path error, basePath=%s destPath=%s", s.sourceAbsPath, sourceFileRel)
		return "", err
	}

	return filepath.ToSlash(filepath.Join(s.basePath, sourceFileRel)), nil
}

func (s *driverPushClientSync) fileInfoCompare(sourcePath string) (equal bool) {
	if s.forceChecksum {
		return false
	}
	fiv, ok := s.files.Load(sourcePath)
	if ok {
		fi := fiv.(contract.FileInfo)
		sourceStat, err := os.Stat(sourcePath)
		if err != nil {
			log.Error(err, "get source file stat error => %s", sourcePath)
			return false
		}
		if sourceStat.Size() == fi.Size && sourceStat.ModTime().Unix() == fi.MTime {
			return true
		}
	}
	return false
}

// storeFileInfo store the source file info to compare file whether it is changed or not
func (s *driverPushClientSync) storeFileInfo(sourcePath string) {
	if s.forceChecksum {
		return
	}
	sourceStat, err := os.Stat(sourcePath)
	if err != nil {
		log.Error(err, "get source file stat error => %s", sourcePath)
		return
	}
	s.files.Store(sourcePath, contract.FileInfo{
		Path:  sourcePath,
		IsDir: contract.FsNotDir,
		Size:  sourceStat.Size(),
		MTime: sourceStat.ModTime().Unix(),
	})
}

func (s *driverPushClientSync) initFileInfo() error {
	initMax := 5000
	count := 0
	errWalkDirStop := errors.New("walk dir stop")
	err := filepath.WalkDir(s.sourceAbsPath, func(path string, d iofs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		s.storeFileInfo(path)
		count++
		if count >= initMax {
			return errWalkDirStop
		}
		return nil
	})
	if errors.Is(err, errWalkDirStop) {
		err = nil
	}
	return err
}
