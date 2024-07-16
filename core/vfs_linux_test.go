package core

import (
	"testing"
)

func TestVFS_SSHConfig(t *testing.T) {
	testCases := []struct {
		name string
		path string
	}{
		{"testVFSSFTPSSHConfigDestPath", testVFSSFTPSSHConfigDestPath},
		{"testVFSSFTPSSHConfigDestPathWithDefaultIdentity", testVFSSFTPSSHConfigDestPathWithDefaultIdentity},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			vfs := NewVFS(tc.path)
			assert(t, vfs.Host() == "127.0.0.1", "compare Host error, expect:%s, actual:%s", "127.0.0.1", vfs.SSHConfig().Username)
			assert(t, vfs.SSHConfig().Username == "hello", "compare Username error, expect:%s, actual:%s", "hello", vfs.SSHConfig().Username)
			assert(t, vfs.SSHConfig().Password == "sftp_pwd", "compare Password error, expect:%s, actual:%s", "sftp_pwd", vfs.SSHConfig().Password)
			assert(t, vfs.Port() == 8818, "compare Port error, expect:%d, actual:%d", 8818, vfs.Port())
		})
	}
}

func TestVFS_SSHConfigWithCover(t *testing.T) {
	vfs := NewVFS(testVFSSFTPSSHConfigDestPathWithCover)
	assert(t, vfs.Host() == "127.0.0.1", "compare Host error, expect:%s, actual:%s", "127.0.0.1", vfs.Host())
	assert(t, vfs.SSHConfig().Username == "sftp_user", "compare Username error, expect:%s, actual:%s", "sftp_user", vfs.SSHConfig().Username)
	assert(t, vfs.SSHConfig().Password == "sftp_pwd", "compare Password error, expect:%s, actual:%s", "sftp_pwd", vfs.SSHConfig().Password)
	assert(t, vfs.Port() == 8818, "compare Port error, expect:%d, actual:%d", 8818, vfs.Port())
}
