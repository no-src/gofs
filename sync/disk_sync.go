package sync

import (
	"bufio"
	"errors"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/util"
	"github.com/no-src/log"
	"io/fs"
	"os"
	"path/filepath"
)

type diskSync struct {
	baseSync
	src         core.VFS
	dest        core.VFS
	srcAbsPath  string
	destAbsPath string
}

// NewDiskSync create a diskSync instance
// src is source path to read
// dest is dest path to write
func NewDiskSync(src, dest core.VFS, enableLogicallyDelete bool) (s Sync, err error) {
	if len(src.Path()) == 0 {
		err = errors.New("src is not found")
		return nil, err
	}
	if len(dest.Path()) == 0 {
		err = errors.New("dest is not found")
		return nil, err
	}

	srcAbsPath, err := filepath.Abs(src.Path())
	if err != nil {
		return nil, err
	}

	destAbsPath, err := filepath.Abs(dest.Path())
	if err != nil {
		return nil, err
	}

	s = &diskSync{
		srcAbsPath:  srcAbsPath,
		destAbsPath: destAbsPath,
		src:         src,
		dest:        dest,
		baseSync:    newBaseSync(enableLogicallyDelete),
	}
	return s, nil
}

// Create creates the source file or dir to the dest
func (s *diskSync) Create(path string) error {
	dest, err := s.buildDestAbsFile(path)
	if err != nil {
		return err
	}
	exist, err := util.FileExist(dest)
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
		f, err := util.CreateFile(dest)
		defer func() {
			if err = f.Close(); err != nil {
				log.Error(err, "Create:close file error")
			}
		}()
		if err != nil {
			return err
		}
	}
	_, aTime, mTime, err := util.GetFileTime(path)
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

// Write sync the src file to the dest
func (s *diskSync) Write(path string) error {
	dest, err := s.buildDestAbsFile(path)
	if err != nil {
		return err
	}

	isDir, err := s.IsDir(path)
	if err != nil {
		return err
	}

	if isDir {
		return s.SyncOnce(path)
	} else {
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func() {
			if err = srcFile.Close(); err != nil {
				log.Error(err, "Write:close the src file error")
			}
		}()
		srcStat, err := srcFile.Stat()
		if err != nil {
			return err
		}

		destFile, err := util.OpenRWFile(dest)
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

		if srcStat.Size() == 0 {
			log.Info("write to the dest file success [size=%d] [%s] -> [%s]", srcStat.Size(), path, dest)
			return nil
		}

		reader := bufio.NewReader(srcFile)
		writer := bufio.NewWriter(destFile)

		// if src and dest is the same file, ignore the following steps and return directly
		if srcStat.Size() > 0 && srcStat.Size() == destStat.Size() {
			isSame, err := s.same(srcFile, destFile)
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
			log.Info("write to the dest file success, size[%d => %d] [%s] => [%s]", srcStat.Size(), n, path, dest)

			// change file times
			if _, aTime, mTime, err := util.GetFileTime(path); err == nil {
				if err = os.Chtimes(dest, aTime, mTime); err != nil {
					log.Warn("Write:change file times error => %s =>[%s]", err.Error(), dest)
				}
			} else {
				log.Warn("Write:get file times error => %s =>[%s]", err.Error(), path)
			}
		} else {
			return err
		}
	}
	return nil
}

func (s *diskSync) same(srcFile *os.File, destFile *os.File) (bool, error) {
	srcHash, err := util.MD5FromFile(srcFile)
	if err != nil {
		log.Error(err, "calculate md5 hash of the src file error [%s]", srcFile.Name())
		return false, err
	}

	destHash, err := util.MD5FromFile(destFile)
	if err != nil {
		log.Error(err, "calculate md5 hash of the dest file error [%s]", destFile.Name())
		return false, err
	}

	if len(srcHash) > 0 && srcHash == destHash {
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
	if s.enableLogicallyDelete {
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
// srcFileAbs: src abs file path
func (s *diskSync) buildDestAbsFile(srcFileAbs string) (string, error) {
	srcFileRel, err := filepath.Rel(s.srcAbsPath, srcFileAbs)
	if err != nil {
		log.Error(err, "parse rel path error, basePath=%s destPath=%s", s.srcAbsPath, srcFileRel)
		return "", err
	}
	return filepath.Join(s.destAbsPath, srcFileRel), nil
}

func (s *diskSync) IsDir(path string) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return f.IsDir(), nil
}

// SyncOnce auto sync src directory to dest directory once.
func (s *diskSync) SyncOnce(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	return filepath.WalkDir(absPath, func(currentPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
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

func (s *diskSync) Source() core.VFS {
	return s.src
}

func (s *diskSync) Dest() core.VFS {
	return s.dest
}
