package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/no-src/gofs/util/jsonutil"
)

const (
	testVFSServerPath                     = "rs://127.0.0.1:8105?mode=server&local_sync_disabled=true&path=./source&fs_server=https://127.0.0.1"
	testVFSServerPathWithNoPort           = "rs://127.0.0.1?mode=server&local_sync_disabled=true&path=./source&fs_server=https://127.0.0.1"
	testVFSServerPathWithNoSchemeFsServer = "rs://127.0.0.1:8105?mode=server&local_sync_disabled=true&path=./source&fs_server=127.0.0.1"
)

func TestVFS_MarshalText(t *testing.T) {
	testCases := []struct {
		path string
	}{
		{""},
		{testVFSServerPath},
		{testVFSServerPathWithNoPort},
		{testVFSServerPathWithNoSchemeFsServer},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			vfs := NewVFS(tc.path)
			data, err := jsonutil.Marshal(vfs)
			if err != nil {
				t.Errorf("test duration marshal error =>%s", err)
				return
			}
			var buf bytes.Buffer
			json.HTMLEscape(&buf, []byte(tc.path))
			expect := fmt.Sprintf("\"%s\"", buf.String())
			actual := string(data)
			if actual != expect {
				t.Errorf("test vfs marshal error, expect:%s, actual:%s", expect, actual)
			}
		})
	}
}

func TestVFS_UnmarshalText(t *testing.T) {
	testCases := []struct {
		path string
	}{
		{""},
		{testVFSServerPath},
		{testVFSServerPathWithNoPort},
		{testVFSServerPathWithNoSchemeFsServer},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			var actual VFS
			data := []byte(fmt.Sprintf("\"%s\"", tc.path))
			err := jsonutil.Unmarshal(data, &actual)
			if err != nil {
				t.Errorf("test vfs unmarshal error =>%s", err)
				return
			}
			compareVFS(t, NewVFS(tc.path), actual)
		})
	}
}

func TestNewVFS_WithDefaultPort(t *testing.T) {
	testCases := []struct {
		path string
	}{
		{testVFSServerPathWithNoPort},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			actual := NewVFS(tc.path)
			if remoteServerDefaultPort != actual.Port() {
				t.Errorf("test new vfs with default port error, expect:%d, actual:%d", remoteServerDefaultPort, actual.Port())
			}
		})
	}
}

func TestNewVFS_WithNoSchemeFsServer(t *testing.T) {
	testCases := []struct {
		path   string
		expect string
	}{
		{testVFSServerPathWithNoSchemeFsServer, "https://127.0.0.1"},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			actual := NewVFS(tc.path)
			if tc.expect != actual.FsServer() {
				t.Errorf("test new vfs with no scheme fs server error, expect:%s, actual:%s", tc.expect, actual.FsServer())
			}
		})
	}
}

func TestNewVFS_ReturnError(t *testing.T) {
	testCases := []struct {
		path   string
		expect VFS
	}{
		{testVFSServerPath + string([]byte{127}), NewEmptyVFS()}, // 0x7F DEL
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			compareVFS(t, tc.expect, NewVFS(tc.path))
		})
	}
}

func TestVFSVar_DefaultValue(t *testing.T) {
	testCases := []struct {
		name         string
		defaultValue VFS
	}{
		{"default_empty_vfs", NewEmptyVFS()},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var actual VFS
			VFSVar(&actual, "core_test_vfs_var_default"+tc.name, tc.defaultValue, "test vfs var")
			parseFlag()
			compareVFS(t, tc.defaultValue, actual)
		})
	}
}

func TestVFSVar(t *testing.T) {
	testCases := []struct {
		name         string
		path         string
		defaultValue VFS
	}{
		{"testVFSServerPath", testVFSServerPath, NewEmptyVFS()},
		{"testVFSServerPathWithNoPort", testVFSServerPathWithNoPort, NewEmptyVFS()},
		{"testVFSServerPathWithNoSchemeFsServer", testVFSServerPathWithNoSchemeFsServer, NewEmptyVFS()},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var actual VFS
			expect := NewVFS(tc.path)
			flagName := "core_test_vfs_var" + tc.name
			VFSVar(&actual, flagName, tc.defaultValue, "test vfs var")
			parseFlag(fmt.Sprintf("-%s=%s", flagName, tc.path))
			compareVFS(t, expect, actual)
		})
	}
}

func TestVFSFlag_DefaultValue(t *testing.T) {
	testCases := []struct {
		name         string
		defaultValue VFS
	}{
		{"default_empty_vfs", NewEmptyVFS()},
		{"with_normal_vfs", NewVFS(testVFSServerPath)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var actual *VFS
			flagName := "core_test_vfs_flag_default" + tc.name
			actual = VFSFlag(flagName, tc.defaultValue, "test vfs flag")
			parseFlag()
			compareVFS(t, tc.defaultValue, *actual)
		})
	}
}

func TestVFSFlag(t *testing.T) {
	testCases := []struct {
		name         string
		path         string
		defaultValue VFS
	}{
		{"testVFSServerPath", testVFSServerPath, NewEmptyVFS()},
		{"testVFSServerPathWithNoPort", testVFSServerPathWithNoPort, NewEmptyVFS()},
		{"testVFSServerPathWithNoSchemeFsServer", testVFSServerPathWithNoSchemeFsServer, NewEmptyVFS()},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expect := NewVFS(tc.path)
			flagName := "core_test_vfs_flag" + tc.name
			actual := VFSFlag(flagName, tc.defaultValue, "test vfs flag")
			parseFlag(fmt.Sprintf("-%s=%s", flagName, tc.path))
			compareVFS(t, expect, *actual)
		})
	}
}

func compareVFS(t *testing.T, expect, actual VFS) {
	if expect.original != actual.original {
		t.Errorf("compare vfs original error, expect:%s, actual:%s", expect.original, actual.original)
	}

	if expect.Path() != actual.Path() {
		t.Errorf("compare vfs Path error, expect:%s, actual:%s", expect.Path(), actual.Path())
	}

	expectAbs, err := expect.Abs()
	if err != nil {
		t.Errorf("compare vfs Abs error, parse expect abs error =>%s", err)
		return
	}
	actualAbs, err := actual.Abs()
	if err != nil {
		t.Errorf("compare vfs Abs error, parse actual abs error =>%s", err)
		return
	}
	if expectAbs != actualAbs {
		t.Errorf("compare vfs Abs error, expect:%s, actual:%s", expectAbs, actualAbs)
	}

	if expect.IsEmpty() != actual.IsEmpty() {
		t.Errorf("compare vfs IsEmpty error, expect:%v, actual:%v", expect.IsEmpty(), actual.IsEmpty())
	}

	if expect.Type() != actual.Type() {
		t.Errorf("compare vfs Type error, expect:%v, actual:%v", expect.Type(), actual.Type())
	}

	if expect.Host() != actual.Host() {
		t.Errorf("compare vfs Host error, expect:%s, actual:%s", expect.Host(), actual.Host())
	}

	if expect.Port() != actual.Port() {
		t.Errorf("compare vfs Port error, expect:%d, actual:%d", expect.Port(), actual.Port())
	}

	if expect.IsDisk() != actual.IsDisk() {
		t.Errorf("compare vfs IsDisk error, expect:%v, actual:%v", expect.IsDisk(), actual.IsDisk())
	}

	if expect.Server() != actual.Server() {
		t.Errorf("compare vfs Server error, expect:%v, actual:%v", expect.Server(), actual.Server())
	}

	if expect.FsServer() != actual.FsServer() {
		t.Errorf("compare vfs FsServer error, expect:%s, actual:%s", expect.FsServer(), actual.FsServer())
	}

	if expect.LocalSyncDisabled() != actual.LocalSyncDisabled() {
		t.Errorf("compare vfs LocalSyncDisabled error, expect:%v, actual:%v", expect.LocalSyncDisabled(), actual.LocalSyncDisabled())
	}
}
