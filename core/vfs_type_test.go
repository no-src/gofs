package core

import "testing"

func TestVFSType(t *testing.T) {
	testVFSType(t, -1, "Unknown")
	testVFSType(t, Unknown, "Unknown")
	testVFSType(t, Disk, "Disk")
	testVFSType(t, RemoteDisk, "RemoteDisk")
	testVFSType(t, FTP, "FTP")
	testVFSType(t, SFTP, "SFTP")
	testVFSType(t, NetworkDisk, "NetworkDisk")
	testVFSType(t, SharedFolder, "SharedFolder")
}

func testVFSType(t *testing.T, vfsType VFSType, desc string) {
	if vfsType.String() != desc {
		t.Errorf("test VFSType string error, expect:%s, actual:%s", vfsType.String(), desc)
	}
}
