package fs

import (
	"errors"
	"os"
	"testing"

	"github.com/no-src/gofs/util/osutil"
)

func TestIsDeleted(t *testing.T) {
	testCases := []struct {
		path   string
		expect bool
	}{
		{"/test/dir", false},
		{"/test/README.MD", false},
		{"./test/dir", false},
		{"./test/README.MD", false},

		{".1643351810.deleted", false},
		{"/test/README.MD.164335181.deleted", false},
		{"./test/dir.164335181.deleted", false},
		{"./test/README.MD.164335181.deleted", false},

		{"/test/dir.1643351810.deleted", true},
		{"/test/README.MD.1643351810.deleted", true},
		{"./test/dir.1643351810.deleted", true},
		{"./test/README.MD.1643351810.deleted", true},

		{"/test/dir.16433518101.deleted", true},
		{"/test/README.MD.16433518101.deleted", true},
		{"./test/dir.16433518101.deleted", true},
		{"./test/README.MD.16433518101.deleted", true},

		{".1643351810.DELETED", false},
		{"/test/dir.1643351810.DELETED", true},

		{"C:\\test\\dir", false},
		{"C:\\test\\README.MD", false},
		{".\\test\\dir", false},
		{".\\test\\README.MD", false},

		{"C:\\test\\README.MD.164335181.deleted", false},
		{".\\test\\dir.164335181.deleted", false},
		{".\\test\\README.MD.164335181.deleted", false},

		{"C:\\test\\dir.1643351810.deleted", true},
		{"C:\\test\\README.MD.1643351810.deleted", true},
		{".\\test\\dir.1643351810.deleted", true},
		{".\\test\\README.MD.1643351810.deleted", true},

		{"C:\\test\\dir.16433518101.deleted", true},
		{"C:\\test\\README.MD.16433518101.deleted", true},
		{".\\test\\dir.16433518101.deleted", true},
		{".\\test\\README.MD.16433518101.deleted", true},

		{"C:\\test\\README.MD.164335181.DELETED", false},
		{"C:\\test\\dir.1643351810.DELETED", true},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			actual := IsDeleted(tc.path)
			if actual != tc.expect {
				t.Errorf("[%s] => expect: %v, but actual: %v \n", tc.path, tc.expect, actual)
			}
		})
	}
}

func TestClearDeletedFile(t *testing.T) {
	testCases := []struct {
		path string
	}{
		{"./"},
		{"./delete_test.notfoud"},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			if err := ClearDeletedFile(tc.path); err != nil {
				t.Errorf("clear deleted file error %s => %v", tc.path, err)
			}
		})
	}
}

func TestClearDeletedFile_ReturnError_RemoveAllAlwaysReturnError(t *testing.T) {
	removeAll = removeAllErrorMock
	isDeleted = isDeleteMock
	defer func() {
		removeAll = os.RemoveAll
		isDeleted = isDeletedCore
	}()

	testCases := []struct {
		path string
	}{
		{"./"},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			if err := ClearDeletedFile(tc.path); err == nil {
				t.Errorf("clear deleted file expect to get an error but get nil => %s", tc.path)
			}
		})
	}
}

func TestClearDeletedFile_ReturnError_NotExistAlwaysFalse(t *testing.T) {
	if !osutil.IsWindows() {
		isNotExist = isNotExistAlwaysFalseMock
		defer func() {
			isNotExist = os.IsNotExist
		}()
	}

	testCases := []struct {
		path string
	}{
		{"|/"},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			if err := ClearDeletedFile(tc.path); err == nil {
				t.Errorf("clear deleted file expect to get an error but get nil => %s", tc.path)
			}
		})
	}
}

func TestClearDeletedFile_RemoveAllAlwaysReturnNil(t *testing.T) {
	removeAll = removeAllSuccessMock
	isDeleted = isDeleteMock
	defer func() {
		removeAll = os.RemoveAll
		isDeleted = isDeletedCore
	}()

	testCases := []struct {
		path string
	}{
		{"./"},
		{"./delete_test.notfoud"},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			if err := ClearDeletedFile(tc.path); err != nil {
				t.Errorf("clear deleted file error %s => %v", tc.path, err)
			}
		})
	}
}

func TestToDeletedPath(t *testing.T) {
	testCases := []struct {
		path string
	}{
		{"./delete_test.go"},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			deletedPath := ToDeletedPath(tc.path)
			if len(deletedPath) == 0 {
				t.Errorf("convert to deleted path error %s => %s", tc.path, deletedPath)
			}
		})
	}
}

func TestLogicallyDelete(t *testing.T) {
	testCases := []struct {
		path string
	}{
		{"./delete_test_notfound.go"},
		{"./delete_test_notfound.go.1643351810.deleted"},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			if err := LogicallyDelete(tc.path); err != nil {
				t.Errorf("logical delete error %s => %v", tc.path, err)
			}
		})
	}
}

func TestLogicallyDelete_MockRename(t *testing.T) {
	rename = renameMock
	defer func() {
		rename = os.Rename
	}()

	testCases := []struct {
		path string
	}{
		{"./delete_test.go"},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			if err := LogicallyDelete(tc.path); err != nil {
				t.Errorf("logical delete error %s => %v", tc.path, err)
			}
		})
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
