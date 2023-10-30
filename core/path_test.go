package core

import (
	"testing"
)

func TestNewPath(t *testing.T) {
	testCases := []struct {
		path           string
		fsType         VFSType
		expectBucket   string
		expectBasePath string
	}{
		{"/workspace", Disk, "", "/workspace"},
		{"/workspace", SFTP, "", "/workspace"},
		{"myBucket:/workspace", MinIO, "myBucket", "/workspace"},
		{"myBucket", MinIO, "myBucket", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			p := newPath(tc.path, tc.fsType)
			if p.Bucket() != tc.expectBucket {
				t.Errorf("test new path error, expect bucket:%s, actual:%s", tc.expectBucket, p.Bucket())
				return
			}
			if p.Base() != tc.expectBasePath {
				t.Errorf("test new path error, expect base path:%s, actual:%s", tc.expectBasePath, p.Base())
				return
			}
			if p.String() != tc.path {
				t.Errorf("test new path error, expect path:%s, actual:%s", tc.path, p.String())
			}
		})
	}
}
