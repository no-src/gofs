package sync

import (
	iofs "io/fs"
	"os"
	"path/filepath"

	"github.com/no-src/gofs/driver"
	"github.com/no-src/gofs/fs"
	"github.com/no-src/gofs/ignore"
	"github.com/no-src/log"
)

type driverPushClientSync struct {
	diskSync

	basePath string
	client   driver.Driver
}

func (s *driverPushClientSync) start() error {
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

	// write to sftp server
	destPath, err := s.buildDestAbsFile(path)
	if err != nil {
		return err
	}

	err = s.client.Write(path, destPath)
	if err == nil {
		if _, aTime, mTime, err := fs.GetFileTime(path); err == nil {
			log.ErrorIf(s.client.Chtimes(destPath, aTime, mTime), "[%s push client sync] [write] change file times error", s.client.DriverName())
		}
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
