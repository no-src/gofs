package checksum

import (
	"testing"

	"github.com/no-src/gofs/logger"
	"github.com/no-src/gofs/util/hashutil"
)

func TestPrintChecksum(t *testing.T) {
	logger := logger.NewTestLogger()
	defer logger.Close()
	path := "./checksum_test.go"
	err := PrintChecksum(path, 1024*1024, 10, hashutil.DefaultHash, logger)
	if err != nil {
		t.Errorf("test PrintChecksum error => %v", err)
	}
}

func TestPrintChecksum_ReturnError(t *testing.T) {
	logger := logger.NewTestLogger()
	defer logger.Close()
	path := "./"
	err := PrintChecksum(path, 1024*1024, 10, hashutil.DefaultHash, logger)
	if err == nil {
		t.Errorf("test PrintChecksum expect to get an error but get nil")
	}
}

func TestPrintChecksum_InvalidAlgorithm(t *testing.T) {
	logger := logger.NewTestLogger()
	defer logger.Close()
	path := "./checksum_test.go"
	err := PrintChecksum(path, 1024*1024, 10, "", logger)
	if err == nil {
		t.Errorf("test PrintChecksum expect to get an error but get nil")
	}
}
