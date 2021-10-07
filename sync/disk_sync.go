package sync

import (
	"bufio"
	"errors"
	"github.com/no-src/gofs/util"
	"github.com/no-src/log"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type diskSync struct {
	srcAbsPath    string
	targetAbsPath string
	bufSize       int
}

// NewDiskSync create a diskSync instance
// src is source path to read
// target is target path to write
// bufSize is read and write buffer byte size
func NewDiskSync(src, target string, bufSize int) (s Sync, err error) {
	if len(src) == 0 {
		err = errors.New("src is not found")
		return nil, err
	}
	if len(target) == 0 {
		err = errors.New("target is not found")
		return nil, err
	}
	if bufSize <= 0 {
		err = errors.New("bufSize must greater than zero")
		return nil, err
	}

	srcAbsPath, err := filepath.Abs(src)
	if err != nil {
		return nil, err
	}

	targetAbsPath, err := filepath.Abs(target)
	if err != nil {
		return nil, err
	}

	s = &diskSync{
		srcAbsPath:    srcAbsPath,
		targetAbsPath: targetAbsPath,
		bufSize:       bufSize,
	}
	return s, nil
}

// Create creates the source file or dir to the target
func (s *diskSync) Create(path string) error {
	target, err := s.buildTargetAbsFile(path)
	if err != nil {
		log.Error(err, "Create:build to target abs file error [%s]", path)
		return err
	}
	isDir, err := s.IsDir(path)
	if err != nil {
		log.Error(err, "Create:check if the path is dir error")
		return err
	}
	if isDir {
		err = os.MkdirAll(target, os.ModePerm)
		if err != nil {
			log.Error(err, "Create:create dir error")
			return err
		}
	} else {
		dir := filepath.Dir(target)
		err = os.MkdirAll(dir, fs.ModePerm)
		if err != nil {
			log.Error(err, "Create:create dir error")
			return err
		}
		f, err := os.Create(target)
		defer f.Close()
		if err != nil {
			log.Error(err, "Create:create file error")
			return err
		}
	}
	log.Debug("create target file success [%s] -> [%s]", path, target)
	return nil
}

// Write sync the src file to the target
func (s *diskSync) Write(path string) error {
	target, err := s.buildTargetAbsFile(path)
	if err != nil {
		log.Error(err, "Write:build to target abs file error [%s]", path)
		return err
	}

	isDir, err := s.IsDir(path)
	if err != nil {
		log.Error(err, "Write:check if the path is dir error")
		return err
	}

	if isDir {
		// ignored
	} else {
		srcFile, err := os.Open(path)
		if err != nil {
			log.Error(err, "Write:open the src file failed")
			return err
		}
		defer srcFile.Close()
		srcStat, err := srcFile.Stat()
		if err != nil {
			log.Error(err, "Write:get the src file stat failed")
			return err
		}

		targetFile, err := os.OpenFile(target, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Error(err, "Write:create the target file failed")
			return err
		}
		defer targetFile.Close()
		targetStat, err := targetFile.Stat()
		if err != nil {
			log.Error(err, "Write:get the target file stat failed")
			return err
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
		}

		block := make([]byte, s.bufSize)
		var wc int64 = 0
		for {
			n, err := reader.Read(block)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Error(err, "Write:read from the src file bytes failed [%s]", path)
				return err
			}
			log.Debug("Write:read from the src file [%d] bytes [%s]", n, path)
			nn, err := writer.Write(block[:n])
			if err != nil {
				log.Error(err, "Write:write to the target file bytes failed [%s]", target)
				return err
			}
			wc += int64(nn)
			progress := float64(wc) / float64(srcStat.Size()) * 100
			log.Debug("Write:write to the target file [%d] bytes, current progress [%d/%d][%.2f%%] [%s]", nn, wc, srcStat.Size(), progress, target)
		}
		err = writer.Flush()
		if err == nil {
			log.Info("write to the target file success [%s] -> [%s]", path, target)
		} else {
			log.Error(err, "Write:flush to the target file failed [%s]", target)
			return err
		}
	}
	return nil
}

func (s *diskSync) same(srcFile *os.File, targetFile *os.File) (bool, error) {
	srcHash, err := util.MD5FromFile(srcFile, s.bufSize)
	if err != nil {
		log.Error(err, "calc src file md5 error [%s]", srcFile.Name())
		return false, err
	}

	targetHash, err := util.MD5FromFile(targetFile, s.bufSize)
	if err != nil {
		log.Error(err, "calc target file md5 error [%s]", targetFile.Name())
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
		log.Error(err, "Remove:build to target abs file error [%s]", path)
		return err
	}
	err = os.RemoveAll(target)
	if err != nil {
		log.Error(err, "Remove:remove the target file error")
	} else {
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
	log.Debug("Chmod not implemented [%s]", path)
	return nil
}

// buildTargetAbsFile build target abs file path
// srcFileAbs: src abs file path
func (s *diskSync) buildTargetAbsFile(srcFileAbs string) (string, error) {
	srcFileRel, err := filepath.Rel(s.srcAbsPath, srcFileAbs)
	if err != nil {
		log.Error(err, "parse rel path error, basepath=%s targpath=%s", s.srcAbsPath, srcFileRel)
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
func (s *diskSync) SyncOnce() error {
	return filepath.WalkDir(s.srcAbsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			err = s.Create(path)
			if err == nil {
				err = s.Write(path)
			}
		}
		return err
	})
}
