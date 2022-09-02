package core

import (
	"testing"
)

func TestVFSType(t *testing.T) {
	testCases := []struct {
		vfsType VFSType
		desc    string
	}{
		{-1, "Unknown"},
		{Unknown, "Unknown"},
		{Disk, "Disk"},
		{RemoteDisk, "RemoteDisk"},
		{FTP, "FTP"},
		{SFTP, "SFTP"},
		{MinIO, "MinIO"},
		{NetworkDisk, "NetworkDisk"},
		{SharedFolder, "SharedFolder"},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			if tc.vfsType.String() != tc.desc {
				t.Errorf("test VFSType string error, expect:%s, actual:%s", tc.vfsType.String(), tc.desc)
			}
		})
	}
}
