package sync

import (
	"bufio"
	"io/fs"
	"os"

	"github.com/no-src/gofs/driver"
	nsfs "github.com/no-src/gofs/fs"
)

type driverPullClientSync struct {
	diskSync

	driver driver.Driver
}

func (s *driverPullClientSync) start() error {
	return s.driver.Connect()
}

func (s *driverPullClientSync) Create(path string) error {
	return s.diskSync.Create(path)
}

func (s *driverPullClientSync) Symlink(oldname, newname string) error {
	return s.diskSync.Symlink(oldname, newname)
}

func (s *driverPullClientSync) Write(path string) error {
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
func (s *driverPullClientSync) write(path, dest string) error {
	sourceFile, err := s.driver.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		s.logger.ErrorIf(sourceFile.Close(), "[%s pull client sync] [write] close the source file error", s.driver.DriverName())
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

	if s.hash.QuickCompare(s.forceChecksum, sourceSize, destSize, sourceStat.ModTime(), destStat.ModTime()) {
		s.logger.Debug("[%s pull client sync] [write] [ignored], the file size and file modification time are both unmodified => %s", s.driver.DriverName(), path)
		return nil
	}

	destFile, err := nsfs.OpenRWFile(dest)
	if err != nil {
		return err
	}
	defer func() {
		s.logger.ErrorIf(destFile.Close(), "[%s pull client sync] [write] close the dest file error", s.driver.DriverName())
	}()

	reader := bufio.NewReader(sourceFile)
	writer := bufio.NewWriter(destFile)

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
		s.logger.Info("[driver-pull] [write] [success] size[%d => %d] [%s] => [%s]", sourceSize, n, path, dest)
		s.chtimes(path, dest)
	}
	return err
}

func (s *driverPullClientSync) Remove(path string) error {
	return s.diskSync.Remove(path)
}

func (s *driverPullClientSync) Rename(path string) error {
	return s.diskSync.Rename(path)
}

func (s *driverPullClientSync) Chmod(path string) error {
	s.logger.Debug("Chmod is unimplemented [%s]", path)
	return nil
}

func (s *driverPullClientSync) IsDir(path string) (bool, error) {
	fi, err := s.driver.Stat(path)
	if err != nil {
		return false, err
	}
	return fi.IsDir(), nil
}

func (s *driverPullClientSync) SyncOnce(path string) error {
	return s.driver.WalkDir(path, func(currentPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if s.pi.MatchPath(currentPath, s.driver.DriverName()+" pull client sync", "sync once") {
			return nil
		}
		return s.syncWalk(currentPath, d, s, s.driver.ReadLink)
	})
}
