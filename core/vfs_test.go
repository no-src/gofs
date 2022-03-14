package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/no-src/gofs/util"
	"testing"
)

const (
	testVFSServerPath                     = "rs://127.0.0.1:8105?mode=server&local_sync_disabled=true&path=./source&fs_server=https://127.0.0.1"
	testVFSServerPathWithNoPort           = "rs://127.0.0.1?mode=server&local_sync_disabled=true&path=./source&fs_server=https://127.0.0.1"
	testVFSServerPathWithNoSchemeFsServer = "rs://127.0.0.1:8105?mode=server&local_sync_disabled=true&path=./source&fs_server=127.0.0.1"
)

func TestVFSMarshalText(t *testing.T) {
	path := testVFSServerPath
	vfs := NewVFS(path)
	data, err := util.Marshal(vfs)
	if err != nil {
		t.Errorf("test duration marshal error =>%s", err)
		return
	}
	var buf bytes.Buffer
	json.HTMLEscape(&buf, []byte(path))
	expect := fmt.Sprintf("\"%s\"", buf.String())
	actual := string(data)
	if actual != expect {
		t.Errorf("test vfs marshal error, expect:%s, actual:%s", expect, actual)
	}
}

func TestVFSUnmarshalText(t *testing.T) {
	path := testVFSServerPath
	var actual VFS
	data := []byte(fmt.Sprintf("\"%s\"", path))
	err := util.Unmarshal(data, &actual)
	if err != nil {
		t.Errorf("test vfs unmarshal error =>%s", err)
		return
	}
	expect := NewVFS(path)
	compareVFS(t, expect, actual)
}

func TestNewVFSWithDefaultPort(t *testing.T) {
	path := testVFSServerPathWithNoPort
	actual := NewVFS(path)
	if remoteServerDefaultPort != actual.Port() {
		t.Errorf("test new vfs with default port error, expect:%d, actual:%d", remoteServerDefaultPort, actual.Port())
	}
}

func TestNewVFSWithNoSchemeFsServer(t *testing.T) {
	expect := "https://127.0.0.1"
	path := testVFSServerPathWithNoSchemeFsServer
	actual := NewVFS(path)
	if expect != actual.FsServer() {
		t.Errorf("test new vfs with no scheme fs server error, expect:%s, actual:%s", expect, actual.FsServer())
	}
}

func TestNewVFSError(t *testing.T) {
	path := testVFSServerPath + string([]byte{127}) // 0x7F DEL
	actual := NewVFS(path)
	expect := NewEmptyVFS()
	compareVFS(t, expect, actual)
}

func TestVFSVarDefaultValue(t *testing.T) {
	defaultValue := NewEmptyVFS()
	expect := defaultValue
	var actual VFS
	VFSVar(&actual, "core_test_vfs_var_default", defaultValue, "test vfs var")
	parseFlag()
	compareVFS(t, expect, actual)
}

func TestVFSVar(t *testing.T) {
	defaultValue := NewEmptyVFS()
	expect := NewVFS(testVFSServerPath)
	var actual VFS
	VFSVar(&actual, "core_test_vfs_var", defaultValue, "test vfs var")
	parseFlag("-core_test_vfs_var=" + testVFSServerPath)
	compareVFS(t, expect, actual)
}

func TestVFSFlagDefaultValue(t *testing.T) {
	defaultValue := NewEmptyVFS()
	expect := defaultValue
	var actual *VFS
	actual = VFSFlag("core_test_vfs_flag_default", defaultValue, "test vfs flag")
	parseFlag()
	compareVFS(t, expect, *actual)
}

func TestVFSFlag(t *testing.T) {
	defaultValue := NewEmptyVFS()
	expect := NewVFS(testVFSServerPath)
	var actual *VFS
	actual = VFSFlag("core_test_vfs_flag", defaultValue, "test vfs flag")
	parseFlag("-core_test_vfs_flag=" + testVFSServerPath)
	compareVFS(t, expect, *actual)
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
