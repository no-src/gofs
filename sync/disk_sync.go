package sync

import (
	"bufio"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/encrypt"
	nsfs "github.com/no-src/gofs/fs"
	"github.com/no-src/gofs/ignore"
	"github.com/no-src/gofs/util/hashutil"
	"github.com/no-src/log"
)

type diskSync struct {
	baseSync

	sourceAbsPath   string
	destAbsPath     string
	chunkSize       int64
	checkpointCount int
	enc             *encrypt.Encrypt

	isDirFn       nsfs.IsDirFunc
	statFn        nsfs.StatFunc
	getFileTimeFn nsfs.GetFileTimeFunc
}

// NewDiskSync create a diskSync instance
// source is source path to read
// dest is dest path to write
func NewDiskSync(source, dest core.VFS, enableLogicallyDelete bool, chunkSize int64, checkpointCount int, forceChecksum bool, encOpt encrypt.Option) (s Sync, err error) {
	return newDiskSync(source, dest, enableLogicallyDelete, chunkSize, checkpointCount, forceChecksum, encOpt)
}

func newDiskSync(source, dest core.VFS, enableLogicallyDelete bool, chunkSize int64, checkpointCount int, forceChecksum bool, encOpt encrypt.Option) (s *diskSync, err error) {
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

	enc, err := encrypt.NewEncrypt(encOpt, sourceAbsPath)
	if err != nil {
		return nil, err
	}

	s = &diskSync{
		sourceAbsPath:   sourceAbsPath,
		destAbsPath:     destAbsPath,
		baseSync:        newBaseSync(source, dest, enableLogicallyDelete, forceChecksum),
		chunkSize:       chunkSize,
		checkpointCount: checkpointCount,
		enc:             enc,
		isDirFn:         nsfs.IsDir,
		statFn:          os.Stat,
		getFileTimeFn:   nsfs.GetFileTime,
	}
	return s, nil
}

// Create creates the source file or dir to the dest
func (s *diskSync) Create(path string) error {
	dest, err := s.buildDestAbsFile(path)
	if err != nil {
		return err
	}

	exist, err := nsfs.FileExist(dest)
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

		// only changes the times when the destination path is a directory
		// because the file's modified time will be used to compare whether file changed
		if err = s.chtimes(path, dest); err != nil {
			return err
		}
	} else {
		dir := filepath.Dir(dest)
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
		f, err := nsfs.CreateFile(dest)
		if err != nil {
			return err
		}
		defer func() {
			log.ErrorIf(f.Close(), "[create] close the dest file error")
		}()
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
		log.ErrorIf(sourceFile.Close(), "[write] close the source file error")
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
	if s.enc.NeedEncrypt(path) {
		// ignore the size compare from encryption file because the size of encryption file may not equal to the source file
		if s.quickCompare(0, 0, sourceStat.ModTime(), destStat.ModTime()) {
			log.Debug("[write] [ignored], the file modification time is unmodified => %s", path)
			return nil
		}
	} else {
		if s.quickCompare(sourceSize, destSize, sourceStat.ModTime(), destStat.ModTime()) {
			log.Debug("[write] [ignored], the file size and file modification time are both unmodified => %s", path)
			return nil
		}

		if s.compare(sourceFile, sourceSize, dest, destSize, &offset) {
			log.Debug("[write] [ignored], the file is unmodified => %s", path)
			return nil
		}
	}

	destFile, err := nsfs.OpenRWFile(dest)
	if err != nil {
		return err
	}
	defer func() {
		log.ErrorIf(destFile.Close(), "[write] close the dest file error")
	}()

	if _, err = sourceFile.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	if _, err = destFile.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	reader := bufio.NewReader(sourceFile)
	writer, err := s.enc.NewWriter(destFile, path, destStat.Name())
	if err != nil {
		return err
	}

	// truncate first before write to file
	err = destFile.Truncate(offset)
	if err != nil {
		return err
	}

	n, err := reader.WriteTo(writer)
	if err != nil {
		return err
	}

	err = writer.Close()

	if err == nil {
		log.Info("[disk] [write] [success] size[%d => %d] [%s] => [%s]", sourceSize, n, path, dest)
		s.chtimes(path, dest)
	}
	return err
}

func (s *diskSync) compare(sourceFile *os.File, sourceSize int64, dest string, destSize int64, offset *int64) (equal bool) {
	if destSize <= 0 {
		return false
	}
	hvs, _ := hashutil.CheckpointsHashFromFile(sourceFile, s.chunkSize, s.checkpointCount)
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
func (s *diskSync) chtimes(source, dest string) error {
	_, aTime, mTime, err := s.getFileTimeFn(source)
	if err == nil {
		if err = os.Chtimes(dest, aTime, mTime); err != nil {
			log.Warn("[write] change file times error => %s =>[%s]", err.Error(), dest)
		}
	} else {
		log.Warn("[write] get file times error => %s =>[%s]", err.Error(), source)
	}
	return err
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
	return s.isDirFn(path)
}

// SyncOnce auto sync source directory to dest directory once.
func (s *diskSync) SyncOnce(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	return filepath.WalkDir(absPath, func(currentPath string, d fs.DirEntry, err error) error {
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
		hvs, err = hashutil.CheckpointsHashFromFileName(path, chunkSize, checkpointCount)
		if err != nil {
			return size, hash, hvs, err
		}
		if len(hvs) > 0 {
			hash = hvs.Last().Hash
		}
	}
	return size, hash, hvs, nil
}
