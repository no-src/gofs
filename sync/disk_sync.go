package sync

import (
	"bufio"
	"errors"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/fs"
	"github.com/no-src/gofs/ignore"
	"github.com/no-src/gofs/util/hashutil"
	"github.com/no-src/log"
	"io"
	iofs "io/fs"
	"os"
	"path/filepath"
)

type diskSync struct {
	baseSync
	sourceAbsPath   string
	destAbsPath     string
	chunkSize       int64
	checkpointCount int
}

// NewDiskSync create a diskSync instance
// source is source path to read
// dest is dest path to write
func NewDiskSync(source, dest core.VFS, enableLogicallyDelete bool, chunkSize int64, checkpointCount int) (s Sync, err error) {
	return newDiskSync(source, dest, enableLogicallyDelete, chunkSize, checkpointCount)
}

func newDiskSync(source, dest core.VFS, enableLogicallyDelete bool, chunkSize int64, checkpointCount int) (s *diskSync, err error) {
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
		sourceAbsPath:   sourceAbsPath,
		destAbsPath:     destAbsPath,
		baseSync:        newBaseSync(source, dest, enableLogicallyDelete),
		chunkSize:       chunkSize,
		checkpointCount: checkpointCount,
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
		if err != nil {
			return err
		}
		defer func() {
			log.ErrorIf(f.Close(), "Create:close file error")
		}()
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
	return s.write(path, dest)
}

// write try to write a file to the destination
func (s *diskSync) write(path, dest string) error {
	sourceFile, err := os.Open(path)
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

	var offset int64
	if destSize > 0 && s.compare(sourceFile, sourceSize, dest, &offset) {
		log.Debug("Write:ignored, the file is unmodified => %s", path)
		return nil
	}

	destFile, err := fs.OpenRWFile(dest)
	if err != nil {
		return err
	}
	defer func() {
		log.ErrorIf(destFile.Close(), "Write:close the dest file error")
	}()

	if _, err = sourceFile.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	if _, err = destFile.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	reader := bufio.NewReader(sourceFile)
	writer := bufio.NewWriter(destFile)

	// truncate first before write to file
	err = destFile.Truncate(offset)
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

func (s *diskSync) compare(sourceFile *os.File, sourceSize int64, dest string, offset *int64) (equal bool) {
	hvs, _ := hashutil.CheckpointsMD5FromFile(sourceFile, s.chunkSize, s.checkpointCount)
	if len(hvs) > 0 && hvs.Last().Offset == sourceSize {
		// if source and dest is the same file, ignore the following steps and return directly
		equal, hv := s.compareHashValues(dest, sourceSize, hvs.Last().Hash, s.chunkSize, hvs)
		if equal {
			return equal
		}

		if hv != nil {
			*offset = hv.Offset
		}
	}
	return false
}

// chtimes change file times
func (s *diskSync) chtimes(source, dest string) {
	if _, aTime, mTime, err := fs.GetFileTime(source); err == nil {
		if err = os.Chtimes(dest, aTime, mTime); err != nil {
			log.Warn("Write:change file times error => %s =>[%s]", err.Error(), dest)
		}
	} else {
		log.Warn("Write:get file times error => %s =>[%s]", err.Error(), source)
	}
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

func (s *diskSync) getFileSizeAndHashCheckpoints(path string, chunkSize int64, checkpointCount int) (size int64, hash string, hvs hashutil.HashValues, err error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return size, hash, hvs, err
	}
	if fileInfo.IsDir() {
		return size, hash, hvs, nil
	}
	size = fileInfo.Size()
	if size > 0 {
		hvs, err = hashutil.CheckpointsMD5FromFileName(path, chunkSize, checkpointCount)
		if err != nil {
			return size, hash, hvs, err
		}
		if len(hvs) > 0 {
			hash = hvs.Last().Hash
		}
	}
	return size, hash, hvs, nil
}
