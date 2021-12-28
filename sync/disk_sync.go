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
	src           core.VFS
	target        core.VFS
	srcAbsPath    string
	targetAbsPath string
	bufSize       int
}

// NewDiskSync create a diskSync instance
// src is source path to read
// target is target path to write
// bufSize is read and write buffer byte size
func NewDiskSync(src, target core.VFS, bufSize int) (s Sync, err error) {
	if len(src.Path()) == 0 {
		err = errors.New("src is not found")
		return nil, err
	}
	if len(target.Path()) == 0 {
		err = errors.New("target is not found")
		return nil, err
	}
	if bufSize <= 0 {
		err = errors.New("bufSize must greater than zero")
		return nil, err
	}

	srcAbsPath, err := filepath.Abs(src.Path())
	if err != nil {
		return nil, err
	}

	targetAbsPath, err := filepath.Abs(target.Path())
	if err != nil {
		return nil, err
	}

	s = &diskSync{
		srcAbsPath:    srcAbsPath,
		targetAbsPath: targetAbsPath,
		bufSize:       bufSize,
		src:           src,
		target:        target,
	}
	return s, nil
}

// Create creates the source file or dir to the target
func (s *diskSync) Create(path string) error {
	target, err := s.buildTargetAbsFile(path)
	if err != nil {
		return err
	}
	isDir, err := s.IsDir(path)
	if err != nil {
		return err
	}
	if isDir {
		err = os.MkdirAll(target, os.ModePerm)
		if err != nil {
			return err
		}
	} else {
		dir := filepath.Dir(target)
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
		f, err := util.CreateFile(target)
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
	err = os.Chtimes(target, aTime, mTime)
	if err != nil {
		return err
	}
	log.Info("create the target file success [%s] -> [%s]", path, target)
	return nil
}

// Write sync the src file to the target
func (s *diskSync) Write(path string) error {
	target, err := s.buildTargetAbsFile(path)
	if err != nil {
		return err
	}

	isDir, err := s.IsDir(path)
	if err != nil {
		return err
	}

	if isDir {
		s.SyncOnce(path)
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

		targetFile, err := util.OpenRWFile(target)
		if err != nil {
			return err
		}
		defer func() {
			if err = targetFile.Close(); err != nil {
				log.Error(err, "Write:close the target file error")
			}
		}()
		targetStat, err := targetFile.Stat()
		if err != nil {
			return err
		}

		if srcStat.Size() == 0 {
			log.Info("write to the target file success [size=%d] [%s] -> [%s]", srcStat.Size(), path, target)
			return nil
		}

		reader := bufio.NewReader(srcFile)
		writer := bufio.NewWriter(targetFile)

		// if src and target is the same file, ignore the following steps and return directly
		if srcStat.Size() > 0 && srcStat.Size() == targetStat.Size() {
			isSame, err := s.same(srcFile, targetFile)
			if err == nil && isSame {
				log.Debug("Write:ignored, the file is unmodified")
				return nil
			}

			// reset the offset
			if _, err = targetFile.Seek(0, 0); err != nil {
				return err
			}
		}

		// truncate first before write to file
		err = targetFile.Truncate(0)
		if err != nil {
			return err
		}

		n, err := reader.WriteTo(writer)
		if err != nil {
			return err
		}

		err = writer.Flush()

		if err == nil {
			log.Info("write to the target file success, size[%d => %d] [%s] => [%s]", srcStat.Size(), n, path, target)

			// change file times
			if _, aTime, mTime, err := util.GetFileTime(path); err == nil {
				if err = os.Chtimes(target, aTime, mTime); err != nil {
					log.Warn("Write:change file times error => %s =>[%s]", err.Error(), target)
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

func (s *diskSync) same(srcFile *os.File, targetFile *os.File) (bool, error) {
	srcHash, err := util.MD5FromFile(srcFile, s.bufSize)
	if err != nil {
		log.Error(err, "calculate md5 hash of the src file error [%s]", srcFile.Name())
		return false, err
	}

	targetHash, err := util.MD5FromFile(targetFile, s.bufSize)
	if err != nil {
		log.Error(err, "calculate md5 hash of the target file error [%s]", targetFile.Name())
		return false, err
	}

	if len(srcHash) > 0 && srcHash == targetHash {
		return true, nil
	}

	return false, nil
}

// Remove removes the source file or dir in target
func (s *diskSync) Remove(path string) error {
	target, err := s.buildTargetAbsFile(path)
	if err != nil {
		return err
	}
	err = os.RemoveAll(target)
	if err == nil {
		log.Info("remove file success [%s] -> [%s]", path, target)
	}
	return err
}

// Rename removes the source file or dir in target, the same as Remove
func (s *diskSync) Rename(path string) error {
	// delete old file, then trigger Create
	return s.Remove(path)
}

func (s *diskSync) Chmod(path string) error {
	log.Debug("Chmod is unimplemented [%s]", path)
	return nil
}

// buildTargetAbsFile build target abs file path
// srcFileAbs: src abs file path
func (s *diskSync) buildTargetAbsFile(srcFileAbs string) (string, error) {
	srcFileRel, err := filepath.Rel(s.srcAbsPath, srcFileAbs)
	if err != nil {
		log.Error(err, "parse rel path error, basePath=%s targetPath=%s", s.srcAbsPath, srcFileRel)
		return "", err
	}
	target := filepath.Join(s.targetAbsPath, srcFileRel)
	return target, nil
}

func (s *diskSync) IsDir(path string) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		log.Error(err, "check path stat error")
		return false, err
	}
	return f.IsDir(), nil
}

// SyncOnce auto sync src directory to target directory once.
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

func (s *diskSync) Target() core.VFS {
	return s.target
}
