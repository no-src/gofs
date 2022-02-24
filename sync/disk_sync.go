package sync

import (
	"bufio"
	"errors"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/fs"
	"github.com/no-src/gofs/ignore"
	"github.com/no-src/gofs/util"
	"github.com/no-src/log"
	iofs "io/fs"
	"os"
	"path/filepath"
)

type diskSync struct {
	baseSync
	sourceAbsPath string
	destAbsPath   string
}

// NewDiskSync create a diskSync instance
// source is source path to read
// dest is dest path to write
func NewDiskSync(source, dest core.VFS, enableLogicallyDelete bool) (s Sync, err error) {
	return newDiskSync(source, dest, enableLogicallyDelete)
}

func newDiskSync(source, dest core.VFS, enableLogicallyDelete bool) (s *diskSync, err error) {
	if source.IsEmpty() {
		return nil, errors.New("source is not found")
	}
	if dest.IsEmpty() {
		return nil, errors.New("dest is not found")
	}

	sourceAbsPath, err := source.Abs()
	if err != nil {
		return nil, err
	}

	destAbsPath, err := dest.Abs()
	if err != nil {
		return nil, err
	}

	s = &diskSync{
		sourceAbsPath: sourceAbsPath,
		destAbsPath:   destAbsPath,
		baseSync:      newBaseSync(source, dest, enableLogicallyDelete),
	}
	return s, nil
}

// Create creates the source file or dir to the dest
func (s *diskSync) Create(path string) error {
	dest, err := s.buildDestAbsFile(path)
	if err != nil {
		return err
	}
	exist, err := fs.FileExist(dest)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}
	isDir, err := s.IsDir(path)
	if err != nil {
		return err
	}
	if isDir {
		err = os.MkdirAll(dest, os.ModePerm)
		if err != nil {
			return err
		}
	} else {
		dir := filepath.Dir(dest)
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
		f, err := fs.CreateFile(dest)
		defer func() {
			if err = f.Close(); err != nil {
				log.Error(err, "Create:close file error")
			}
		}()
		if err != nil {
			return err
		}
	}
	_, aTime, mTime, err := fs.GetFileTime(path)
	if err != nil {
		return err
	}
	err = os.Chtimes(dest, aTime, mTime)
	if err != nil {
		return err
	}
	log.Info("create the dest file success [%s] -> [%s]", path, dest)
	return nil
}

// Write sync the source file to the dest
func (s *diskSync) Write(path string) error {
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
	sourceFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		if err = sourceFile.Close(); err != nil {
			log.Error(err, "Write:close the source file error")
		}
	}()
	sourceStat, err := sourceFile.Stat()
	if err != nil {
		return err
	}

	destFile, err := fs.OpenRWFile(dest)
	if err != nil {
		return err
	}
	defer func() {
		if err = destFile.Close(); err != nil {
			log.Error(err, "Write:close the dest file error")
		}
	}()
	destStat, err := destFile.Stat()
	if err != nil {
		return err
	}

	if sourceStat.Size() == 0 {
		log.Info("write to the dest file success [size=%d] [%s] -> [%s]", sourceStat.Size(), path, dest)
		return nil
	}

	reader := bufio.NewReader(sourceFile)
	writer := bufio.NewWriter(destFile)

	// if source and dest is the same file, ignore the following steps and return directly
	if sourceStat.Size() > 0 && sourceStat.Size() == destStat.Size() {
		isSame, err := s.same(sourceFile, destFile)
		if err == nil && isSame {
			log.Debug("Write:ignored, the file is unmodified => %s", path)
			return nil
		}

		// reset the offset
		if _, err = destFile.Seek(0, 0); err != nil {
			return err
		}
	}

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
		log.Info("write to the dest file success, size[%d => %d] [%s] => [%s]", sourceStat.Size(), n, path, dest)

		// change file times
		if _, aTime, mTime, err := fs.GetFileTime(path); err == nil {
			if err = os.Chtimes(dest, aTime, mTime); err != nil {
				log.Warn("Write:change file times error => %s =>[%s]", err.Error(), dest)
			}
		} else {
			log.Warn("Write:get file times error => %s =>[%s]", err.Error(), path)
		}
	}
	return err
}

func (s *diskSync) same(sourceFile *os.File, destFile *os.File) (bool, error) {
	sourceHash, err := util.MD5FromFile(sourceFile)
	if err != nil {
		log.Error(err, "calculate md5 hash of the source file error [%s]", sourceFile.Name())
		return false, err
	}

	destHash, err := util.MD5FromFile(destFile)
	if err != nil {
		log.Error(err, "calculate md5 hash of the dest file error [%s]", destFile.Name())
		return false, err
	}

	if len(sourceHash) > 0 && sourceHash == destHash {
		return true, nil
	}

	return false, nil
}

// Remove removes the source file or dir in dest
func (s *diskSync) Remove(path string) error {
	return s.remove(path, false)
}

func (s *diskSync) remove(path string, forceDelete bool) error {
	dest, err := s.buildDestAbsFile(path)
	if err != nil {
		return err
	}
	if !forceDelete && s.enableLogicallyDelete {
		err = s.logicallyDelete(dest)
	} else {
		err = os.RemoveAll(dest)
	}
	if err == nil {
		log.Info("remove file success [%s] -> [%s]", path, dest)
	}
	return err
}

// Rename removes the source file or dir in dest, the same as Remove
func (s *diskSync) Rename(path string) error {
	// delete old file, then trigger Create
	return s.remove(path, true)
}

func (s *diskSync) Chmod(path string) error {
	log.Debug("Chmod is unimplemented [%s]", path)
	return nil
}

// buildDestAbsFile build dest abs file path
// sourceFileAbs: source abs file path
func (s *diskSync) buildDestAbsFile(sourceFileAbs string) (string, error) {
	sourceFileRel, err := filepath.Rel(s.sourceAbsPath, sourceFileAbs)
	if err != nil {
		log.Error(err, "parse rel path error, basePath=%s destPath=%s", s.sourceAbsPath, sourceFileRel)
		return "", err
	}
	return filepath.Join(s.destAbsPath, sourceFileRel), nil
}

func (s *diskSync) IsDir(path string) (bool, error) {
	return fs.IsDir(path)
}

// SyncOnce auto sync source directory to dest directory once.
func (s *diskSync) SyncOnce(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	return filepath.WalkDir(absPath, func(currentPath string, d iofs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if ignore.MatchPath(currentPath, "disk sync", "sync once") {
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
