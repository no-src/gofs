package sync

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/no-src/gofs/contract"
	"github.com/no-src/gofs/driver"
	nsfs "github.com/no-src/gofs/fs"
)

type driverPushClientSync struct {
	diskSync

	basePath string
	driver   driver.Driver
	files    sync.Map
}

func newDriverPushClientSync(ds diskSync, basePath string) driverPushClientSync {
	return driverPushClientSync{
		diskSync: ds,
		basePath: basePath,
	}
}

func (s *driverPushClientSync) start() error {
	if err := s.initFileInfo(); err != nil {
		return err
	}
	return s.driver.Connect()
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
		err = s.driver.MkdirAll(destPath)
	} else {
		err = s.driver.Create(destPath)
	}
	return err
}

func (s *driverPushClientSync) Symlink(oldname, newname string) error {
	if !s.dest.LocalSyncDisabled() {
		if err := s.diskSync.Symlink(oldname, newname); err != nil {
			return err
		}
	}
	destPath, err := s.buildDestAbsFile(newname)
	if err != nil {
		return err
	}
	return s.driver.Symlink(oldname, destPath)
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
		s.logger.Debug("[push] [ignored], the file modification time is unmodified => %s", path)
		return nil
	}

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

	err = s.driver.Write(encryptPath, destPath)
	if err == nil {
		s.logger.Info("[%s-driver-push] [write] [success] => %s", s.driver.DriverName(), path)
		if _, aTime, mTime, err := nsfs.GetFileTime(path); err == nil {
			s.logger.ErrorIf(s.driver.Chtimes(destPath, aTime, mTime), "[%s push client sync] [write] change file times error", s.driver.DriverName())
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
		err = s.driver.Remove(destPath)
	}
	s.removeFileInfo(path)
	return err
}

// logicallyDelete delete the path logically
func (s *driverPushClientSync) logicallyDelete(path string) error {
	if nsfs.IsDeleted(path) {
		return nil
	}
	deletedFile := nsfs.ToDeletedPath(path)
	err := s.driver.Rename(path, deletedFile)
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
	s.logger.Debug("Chmod is unimplemented [%s]", path)
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
	return filepath.WalkDir(absPath, func(currentPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if s.pi.MatchPath(currentPath, s.driver.DriverName()+" push client sync", "sync once") {
			return nil
		}
		return s.syncWalk(currentPath, d, s, nsfs.Readlink)
	})
}

func (s *driverPushClientSync) buildDestAbsFile(sourceFileAbs string) (string, error) {
	sourceFileRel, err := filepath.Rel(s.sourceAbsPath, sourceFileAbs)
	if err != nil {
		s.logger.Error(err, "parse rel path error, basePath=%s destPath=%s", s.sourceAbsPath, sourceFileRel)
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
			s.logger.Error(err, "get source file stat error => %s", sourcePath)
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
		s.logger.Error(err, "get source file stat error => %s", sourcePath)
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
	err := filepath.WalkDir(s.sourceAbsPath, func(path string, d fs.DirEntry, err error) error {
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

func (s *driverPushClientSync) removeFileInfo(sourcePath string) {
	if s.forceChecksum {
		return
	}
	s.files.Delete(sourcePath)
}
