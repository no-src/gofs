package checksum

import "testing"

func TestPrintChecksum(t *testing.T) {
	path := "./checksum_test.go"
	err := PrintChecksum(path, 1024*1024, 10)
	if err != nil {
		t.Errorf("test PrintChecksum error => %v", err)
	}
}

func TestPrintChecksum_ReturnError(t *testing.T) {
	path := "./"
	err := PrintChecksum(path, 1024*1024, 10)
	if err == nil {
		t.Errorf("test PrintChecksum expect to get an error but get nil")
	}
}
