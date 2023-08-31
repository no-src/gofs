package sync

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/no-src/gofs/encrypt"
	nsfs "github.com/no-src/gofs/fs"
	"github.com/no-src/gofs/ignore"
	"github.com/no-src/gofs/internal/rate"
	"github.com/no-src/gofs/progress"
	"github.com/no-src/gofs/util/hashutil"
	"github.com/no-src/log"
)

type diskSync struct {
	baseSync

	sourceAbsPath         string
	destAbsPath           string
	chunkSize             int64
	checkpointCount       int
	enableLogicallyDelete bool
	forceChecksum         bool
	progress              bool
	maxTranRate           int64
	enc                   *encrypt.Encrypt
	hash                  hashutil.Hash
	pi                    ignore.PathIgnore
	copyLink              bool
	copyUnsafeLink        bool

	isDirFn       nsfs.IsDirFunc
	statFn        nsfs.StatFunc
	getFileTimeFn nsfs.GetFileTimeFunc
}

// NewDiskSync create a diskSync instance
// source is source path to read
// dest is dest path to write
func NewDiskSync(opt Option) (s Sync, err error) {
	return newDiskSync(opt)
}

func newDiskSync(opt Option) (s *diskSync, err error) {
	// the fields of option
	source := opt.Source
	dest := opt.Dest
	pi := opt.PathIgnore
	encOpt := opt.EncOpt
	chunkSize := opt.ChunkSize
	checkpointCount := opt.CheckpointCount
	forceChecksum := opt.ForceChecksum
	checksumAlgorithm := opt.ChecksumAlgorithm
	enableLogicallyDelete := opt.EnableLogicallyDelete
	progress := opt.Progress
	maxTranRate := opt.MaxTranRate
	copyLink := opt.CopyLink
	copyUnsafeLink := opt.CopyUnsafeLink

	if source.IsEmpty() {
		return nil, errSourceNotFound
	}
	if dest.IsEmpty() {
		return nil, errDestNotFound
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

	hash, err := hashutil.NewHash(checksumAlgorithm)
	if err != nil {
		return nil, err
	}

	s = &diskSync{
		sourceAbsPath:         sourceAbsPath,
		destAbsPath:           destAbsPath,
		baseSync:              newBaseSync(source, dest),
		chunkSize:             chunkSize,
		checkpointCount:       checkpointCount,
		enableLogicallyDelete: enableLogicallyDelete,
		forceChecksum:         forceChecksum,
		progress:              progress,
		maxTranRate:           maxTranRate,
		enc:                   enc,
		hash:                  hash,
		pi:                    pi,
		copyLink:              copyLink,
		copyUnsafeLink:        copyUnsafeLink,
		isDirFn:               nsfs.IsDir,
		statFn:                os.Stat,
		getFileTimeFn:         nsfs.GetFileTime,
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
		err = s.createDir(path, dest)
	} else {
		err = s.createFile(dest)
	}
	if err == nil {
		log.Info("create the dest file success [%s] -> [%s]", path, dest)
	}
	return err
}

func (s *diskSync) createDir(source, dest string) error {
	err := os.MkdirAll(dest, fs.ModePerm)
	if err != nil {
		return err
	}

	// only changes the times when the destination path is a directory
	// because the file's modified time will be used to compare whether file changed
	return s.chtimes(source, dest)
}

func (s *diskSync) createFile(dest string) error {
	dir := filepath.Dir(dest)
	err := os.MkdirAll(dir, fs.ModePerm)
	if err != nil {
		return err
	}
	f, err := nsfs.CreateFile(dest)
	if err != nil {
		return err
	}
	log.ErrorIf(f.Close(), "[create] close the dest file error")
	return nil
}

// Symlink create a symbolic link
func (s *diskSync) Symlink(oldname, newname string) error {
	dest, err := s.buildDestAbsFile(newname)
	if err != nil {
		return err
	}
	return s.symlink(oldname, dest)
}

func (s *diskSync) symlink(oldname, newname string) error {
	if err := os.RemoveAll(newname); err != nil {
		return err
	}
	return nsfs.Symlink(oldname, newname)
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
		if s.hash.QuickCompare(s.forceChecksum, 0, 0, sourceStat.ModTime(), destStat.ModTime()) {
			log.Debug("[write] [ignored], the file modification time is unmodified => %s", path)
			return nil
		}
	} else {
		if s.hash.QuickCompare(s.forceChecksum, sourceSize, destSize, sourceStat.ModTime(), destStat.ModTime()) {
			log.Debug("[write] [ignored], the file size and file modification time are both unmodified => %s", path)
			return nil
		}

		if s.hash.Compare(s.chunkSize, s.checkpointCount, sourceFile, sourceSize, dest, destSize, &offset) {
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

	reader := bufio.NewReader(rate.NewReader(sourceFile, s.maxTranRate))
	writer, err := s.enc.NewWriter(destFile, path, destStat.Name())
	if err != nil {
		return err
	}

	// truncate first before write to file
	err = destFile.Truncate(offset)
	if err != nil {
		return err
	}

	n, err := reader.WriteTo(progress.NewWriterWithEnable(writer, sourceSize-offset, fmt.Sprintf("[sync] => %s", destStat.Name()), s.progress))
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
		err = nsfs.LogicallyDelete(dest)
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
		if s.pi.MatchPath(currentPath, "disk sync", "sync once") {
			return nil
		}
		return s.syncWalk(currentPath, d, s, nsfs.Readlink)
	})
}

func (s *diskSync) syncWalk(currentPath string, d fs.DirEntry, sync Sync, readLink func(path string) (string, error)) (err error) {
	if d.IsDir() {
		err = sync.Create(currentPath)
	} else if nsfs.IsSymlinkMode(d.Type()) {
		err = s.syncSymlink(currentPath, sync, readLink)
	} else {
		err = sync.Create(currentPath)
		if err == nil {
			err = sync.Write(currentPath)
		}
	}
	return err
}

func (s *diskSync) syncSymlink(currentPath string, sync Sync, readLink func(path string) (string, error)) (err error) {
	ok := false
	var realPath string
	if s.copyLink {
		if s.copyUnsafeLink {
			ok = true
		} else {
			realPath, err = readLink(currentPath)
			if err != nil {
				return err
			}
			absRealPath := realPath
			if !filepath.IsAbs(absRealPath) {
				absRealPath = filepath.Join(filepath.Dir(currentPath), realPath)
			}
			// ignore unsafe file
			if isSub, err := nsfs.IsSub(s.sourceAbsPath, realPath); err == nil && isSub {
				ok = true
			}
		}
	} else {
		realPath, err = readLink(currentPath)
		if err != nil {
			return err
		}
	}

	if ok {
		// only for local disk
		var dest string
		dest, err = s.buildDestAbsFile(currentPath)
		if err == nil {
			err = s.deepCopy(currentPath, dest)
		}
	} else {
		err = sync.Symlink(realPath, currentPath)
	}
	return err
}

// deepCopy copy the source file to dest recursively, if source is a symlink, copy the real file to dest,
// only use for local disk mode.
func (s *diskSync) deepCopy(source, dest string) error {
	sourceStat, err := os.Lstat(source)
	if err != nil {
		return err
	}
	// sync dir
	if sourceStat.IsDir() {
		if err = s.createDir(source, dest); err != nil {
			return err
		}
		files, err := os.ReadDir(source)
		if err != nil {
			return err
		}
		for _, f := range files {
			if err = s.deepCopy(filepath.Join(source, f.Name()), filepath.Join(dest, f.Name())); err != nil {
				return err
			}
		}
		return nil
	}

	// sync symlink
	if nsfs.IsSymlinkMode(sourceStat.Mode()) {
		if s.pi.MatchPath(source, "local disk deep copy", "the symlink") {
			return nil
		}
		originalSource := source
		realPath, err := nsfs.Readlink(source)
		if err != nil {
			return err
		}
		if filepath.IsAbs(realPath) {
			source = realPath
		} else {
			source = filepath.Join(filepath.Dir(source), realPath)
		}
		sourceStat, err = os.Stat(source)
		if os.IsNotExist(err) {
			return s.symlink(realPath, dest)
		}
		if err != nil {
			return err
		}
		if s.pi.MatchPath(source, "local disk deep copy", "the symlink's real path") {
			return nil
		}
		// check unsafe link
		if !s.copyUnsafeLink {
			// ignore unsafe file
			isSub, err := nsfs.IsSub(s.sourceAbsPath, source)
			if err != nil {
				return err
			}
			if !isSub {
				log.Info("ignore the symlink because it point outside the source tree, symlink => %s, real file => %s", originalSource, source)
				return nil
			}
		}
		return s.deepCopy(source, dest)
	}

	// sync file
	if err = s.createFile(dest); err != nil {
		return err
	}
	return s.write(source, dest)
}
