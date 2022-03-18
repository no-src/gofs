package fs

import (
	"errors"
	"github.com/no-src/gofs/util/osutil"
	"os"
	"testing"
)

func TestIsDeleted(t *testing.T) {
	testIsDeleted(t, false, "/test/dir")
	testIsDeleted(t, false, "/test/README.MD")
	testIsDeleted(t, false, "./test/dir")
	testIsDeleted(t, false, "./test/README.MD")

	testIsDeleted(t, false, ".1643351810.deleted")
	testIsDeleted(t, false, "/test/README.MD.164335181.deleted")
	testIsDeleted(t, false, "./test/dir.164335181.deleted")
	testIsDeleted(t, false, "./test/README.MD.164335181.deleted")

	testIsDeleted(t, true, "/test/dir.1643351810.deleted")
	testIsDeleted(t, true, "/test/README.MD.1643351810.deleted")
	testIsDeleted(t, true, "./test/dir.1643351810.deleted")
	testIsDeleted(t, true, "./test/README.MD.1643351810.deleted")

	testIsDeleted(t, true, "/test/dir.16433518101.deleted")
	testIsDeleted(t, true, "/test/README.MD.16433518101.deleted")
	testIsDeleted(t, true, "./test/dir.16433518101.deleted")
	testIsDeleted(t, true, "./test/README.MD.16433518101.deleted")

	testIsDeleted(t, false, ".1643351810.DELETED")
	testIsDeleted(t, true, "/test/dir.1643351810.DELETED")

	testIsDeleted(t, false, "C:\\test\\dir")
	testIsDeleted(t, false, "C:\\test\\README.MD")
	testIsDeleted(t, false, ".\\test\\dir")
	testIsDeleted(t, false, ".\\test\\README.MD")

	testIsDeleted(t, false, "C:\\test\\README.MD.164335181.deleted")
	testIsDeleted(t, false, ".\\test\\dir.164335181.deleted")
	testIsDeleted(t, false, ".\\test\\README.MD.164335181.deleted")

	testIsDeleted(t, true, "C:\\test\\dir.1643351810.deleted")
	testIsDeleted(t, true, "C:\\test\\README.MD.1643351810.deleted")
	testIsDeleted(t, true, ".\\test\\dir.1643351810.deleted")
	testIsDeleted(t, true, ".\\test\\README.MD.1643351810.deleted")

	testIsDeleted(t, true, "C:\\test\\dir.16433518101.deleted")
	testIsDeleted(t, true, "C:\\test\\README.MD.16433518101.deleted")
	testIsDeleted(t, true, ".\\test\\dir.16433518101.deleted")
	testIsDeleted(t, true, ".\\test\\README.MD.16433518101.deleted")

	testIsDeleted(t, false, "C:\\test\\README.MD.164335181.DELETED")
	testIsDeleted(t, true, "C:\\test\\dir.1643351810.DELETED")

}

func testIsDeleted(t *testing.T, expect bool, path string) {
	actual := IsDeleted(path)
	if actual != expect {
		t.Logf("[%s] => expect: %v, but actual: %v \n", path, expect, actual)
		t.Fail()
	}
}

func TestClearDeletedFile(t *testing.T) {
	path := "./"
	err := ClearDeletedFile(path)
	if err != nil {
		t.Errorf("clear deleted file error %s => %v", path, err)
		return
	}

	path = "./delete_test.notfoud"
	err = ClearDeletedFile(path)
	if err != nil {
		t.Errorf("clear deleted file error %s => %v", path, err)
		return
	}
}

func TestClearDeletedFileRemoveAllError(t *testing.T) {
	removeAll = removeAllErrorMock
	isDeleted = isDeleteMock
	defer func() {
		removeAll = os.RemoveAll
		isDeleted = isDeletedCore
	}()
	path := "./"
	err := ClearDeletedFile(path)
	if err == nil {
		t.Errorf("clear deleted file should get an error => %s", path)
		return
	}
}

func TestClearDeletedFilePathError(t *testing.T) {
	if !osutil.IsWindows() {
		isNotExist = isNotExistAlwaysFalseMock
		defer func() {
			isNotExist = os.IsNotExist
		}()
	}
	path := "|/"
	err := ClearDeletedFile(path)
	if err == nil {
		t.Errorf("clear deleted file should get an error => %s", path)
		return
	}
}

func TestClearDeletedFileSuccess(t *testing.T) {
	removeAll = removeAllSuccessMock
	isDeleted = isDeleteMock
	defer func() {
		removeAll = os.RemoveAll
		isDeleted = isDeletedCore
	}()
	path := "./"
	err := ClearDeletedFile(path)
	if err != nil {
		t.Errorf("clear deleted file error %s => %v", path, err)
		return
	}
}

func TestToDeletedPath(t *testing.T) {
	path := "./delete_test.go"
	deletedPath := toDeletedPath(path)
	if len(deletedPath) == 0 {
		t.Errorf("convert to deleted path error %s => %s", path, deletedPath)
		return
	}
}

func TestLogicallyDelete(t *testing.T) {
	path := "./delete_test_notfound.go"
	err := LogicallyDelete(path)
	if err != nil {
		t.Errorf("logical delete error %s => %v", path, err)
		return
	}

	path = "./delete_test_notfound.go.1643351810.deleted"
	err = LogicallyDelete(path)
	if err != nil {
		t.Errorf("logical delete error %s => %v", path, err)
		return
	}
}

func TestLogicallyDeleteSuccess(t *testing.T) {
	rename = renameMock
	defer func() {
		rename = os.Rename
	}()
	path := "./delete_test.go"
	err := LogicallyDelete(path)
	if err != nil {
		t.Errorf("logical delete error %s => %v", path, err)
		return
	}
}

func renameMock(oldpath, newpath string) error {
	return nil
}

func removeAllSuccessMock(path string) error {
	return nil
}

func removeAllErrorMock(path string) error {
	return errors.New("remove all error test")
}

func isDeleteMock(path string) bool {
	return true
}
