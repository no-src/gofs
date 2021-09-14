package sync

import (
	"bufio"
	"errors"
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

		targetFile, err := os.Create(target)
		if err != nil {
			log.Error(err, "Write:create the target file failed")
			return err
		}
		defer targetFile.Close()

		block := make([]byte, s.bufSize)
		reader := bufio.NewReader(srcFile)
		writer := bufio.NewWriter(targetFile)
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
			log.Debug("Write:write to the target file [%d] bytes [%s]", nn, target)
		}
		err = writer.Flush()
		if err == nil {
			log.Log("Write:write to the target file success [%s]", target)
		} else {
			log.Error(err, "Write:flush to the target file failed [%s]", target)
			return err
		}
	}
	return nil
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
	}
	return err
}

// Rename removes the source file or dir in target, the same as Remove
func (s *diskSync) Rename(path string) error {
	// delete old file, then trigger Create
	return s.Remove(path)
}

func (s *diskSync) Chmod(path string) error {
	panic("Chmod not implemented")
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
