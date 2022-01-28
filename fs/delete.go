package fs

import (
	"fmt"
	"os"
	"time"
)

// LogicallyDelete delete the path logically
func LogicallyDelete(path string) error {
	deletedFile := toDeletedPath(path)
	err := os.Rename(path, deletedFile)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// toDeletedPath convert to the logically deleted file name
func toDeletedPath(path string) string {
	return fmt.Sprintf("%s.%d.deleted", path, time.Now().Unix())
}
