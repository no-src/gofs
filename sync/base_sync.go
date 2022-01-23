package sync

import (
	"fmt"
	"os"
	"time"
)

type baseSync struct {
	enableLogicallyDelete bool
}

func newBaseSync(enableLogicallyDelete bool) baseSync {
	return baseSync{
		enableLogicallyDelete: enableLogicallyDelete,
	}
}

// logicallyDelete delete the path logically
func (s *baseSync) logicallyDelete(path string) error {
	deletedFile := s.deletedFileName(path)
	err := os.Rename(path, deletedFile)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// deletedFileName convert to the logically deleted file name
func (s *baseSync) deletedFileName(path string) string {
	return fmt.Sprintf("%s.%d.deleted", path, time.Now().Unix())
}
