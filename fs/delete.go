package fs

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/no-src/gofs/logger"
)

var deletedPathRegexp *regexp.Regexp

// LogicallyDelete delete the path logically
func LogicallyDelete(path string) error {
	if IsDeleted(path) {
		return nil
	}
	deletedFile := ToDeletedPath(path)
	err := rename(path, deletedFile)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// IsDeleted is deleted path
func IsDeleted(path string) bool {
	return isDeleted(path)
}

func isDeletedCore(path string) bool {
	return deletedPathRegexp.MatchString(path)
}

// ClearDeletedFile remove all the deleted files in the path
func ClearDeletedFile(clearPath string, logger *logger.Logger) error {
	return filepath.WalkDir(clearPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil && isNotExist(err) {
			return nil
		}
		if err != nil {
			return err
		}
		if IsDeleted(path) {
			err = removeAll(path)
			if err != nil {
				logger.Error(err, "remove the deleted files error => [%s]", path)
			} else {
				logger.Debug("remove the deleted files success => [%s]", path)
			}
		}
		return err
	})
}

// ToDeletedPath convert to the logically deleted file name
func ToDeletedPath(path string) string {
	return fmt.Sprintf("%s.%d.deleted", path, time.Now().Unix())
}

var (
	removeAll = os.RemoveAll
	rename    = os.Rename
	isDeleted = isDeletedCore
)

func init() {
	deletedPathRegexp = regexp.MustCompile(`^[\s\S]+\.[0-9]{10,}\.(?i)deleted$`)
}
