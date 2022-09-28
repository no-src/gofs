package fs

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/no-src/gofs/util/osutil"
)

const (
	testNotFoundFilePath = "./fs_test_not_found.go"
	testExistFilePath    = "./fs_test.go"
)

func TestGetFileTime(t *testing.T) {
	testCases := []struct {
		path string
	}{
		{testExistFilePath},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			if _, _, _, err := GetFileTime(tc.path); err != nil {
				t.Errorf("get file time error %s => %v", tc.path, err)
			}
		})
	}
}

func TestGetFileTime_ReturnError(t *testing.T) {
	testCases := []struct {
		path string
	}{
		{testNotFoundFilePath},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			if _, _, _, err := GetFileTime(tc.path); err == nil {
				t.Errorf("get file time error, expect to get an error but get nil => %s", tc.path)
			}
		})
	}
}

func TestGetFileTimeBySys_ReturnError(t *testing.T) {
	testCases := []struct {
		name string
		sys  any
	}{
		{"nil sys", nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, _, _, err := GetFileTimeBySys(nil); err == nil {
				t.Errorf("test GetFileTimeBySys expect to get an error but get nil")
			}
		})
	}
}

func TestFileExist(t *testing.T) {
	testCases := []struct {
		path   string
		expect bool
	}{
		{testExistFilePath, true},
		{testNotFoundFilePath, false},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			exist, err := FileExist(tc.path)
			if err != nil {
				t.Errorf("check file exist error %s => %v", tc.path, err)
				return
			}
			if exist != tc.expect {
				t.Errorf("check file exist error, exist expect:%v,actual:%v => %s", tc.expect, exist, tc.path)
			}
		})
	}
}

func TestFileExist_ReturnError(t *testing.T) {
	if !osutil.IsWindows() {
		isNotExist = isNotExistAlwaysFalseMock
		defer func() {
			isNotExist = os.IsNotExist
		}()
	}

	testCases := []struct {
		name   string
		path   string
		expect bool
	}{
		{"invalid path", "|/", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := FileExist(tc.path); err == nil {
				t.Errorf("test file exist error, expect to get an error but get nil => %s", tc.path)
			}
		})
	}
}

func TestCreateFile(t *testing.T) {
	testCases := []struct {
		path string
	}{
		{testExistFilePath},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			if _, err := CreateFile(tc.path); err != nil {
				t.Errorf("create file error %s => %v", tc.path, err)
			}
		})
	}
}

func TestOpenRWFile(t *testing.T) {
	testCases := []struct {
		path string
	}{
		{testExistFilePath},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			if _, err := OpenRWFile(tc.path); err != nil {
				t.Errorf("open read write file error %s => %v", tc.path, err)
			}
		})
	}
}

func TestIsDir(t *testing.T) {
	testCases := []struct {
		path string
	}{
		{testExistFilePath},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			if _, err := IsDir(tc.path); err != nil {
				t.Errorf("check path is dir error %s => %v", tc.path, err)
			}
		})
	}
}

func TestIsDir_ReturnError(t *testing.T) {
	testCases := []struct {
		path string
	}{
		{testNotFoundFilePath},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			if _, err := IsDir(tc.path); err == nil {
				t.Errorf("check path is dir error, expect to get an error but get nil => %s", tc.path)
			}
		})
	}
}

func TestIsEOF(t *testing.T) {
	testCases := []struct {
		path string
	}{
		{testExistFilePath},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			f, err := os.Open(tc.path)
			if err != nil {
				t.Errorf("test IsEOF error, open file error [%s] => %s", tc.path, err)
				return
			}
			// move to end
			_, err = f.Seek(0, io.SeekEnd)
			if err != nil {
				t.Errorf("test IsEOF error, seek file error [%s] => %s", tc.path, err)
				return
			}
			data := make([]byte, 1024)
			_, err = f.Read(data)
			if !IsEOF(err) {
				t.Errorf("test IsEOF error, read file error [%s] => %s", tc.path, err)
			}
		})
	}
}

func TestIsNonEOF(t *testing.T) {
	testCases := []struct {
		path string
	}{
		{testNotFoundFilePath},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			_, err := os.Stat(tc.path)
			if !IsNonEOF(err) {
				t.Errorf("test IsNonEOF error, get actual err:%s", err)
			}
		})
	}
}

func TestIsSub(t *testing.T) {
	testCases := []struct {
		parent string
		child  string
		isSub  bool
	}{
		{"", "", true},
		{".", ".", true},
		{"..", "..", true},

		{"/a", "/a/b", true},
		{"/a/b1", "/a/b2", false},
		{"/a1", "/a2", false},
		{"/a/b", "/a", false},
		{"/", "/", true},
		{"/a0", "/a0", true},

		{"./a", "./a/b", true},
		{"./a/b1", "./a/b2", false},
		{"./a1", "./a2", false},
		{"./a/b", "./a", false},
		{"./", "./", true},
		{"./a0", "./a0", true},

		{"../a", "../a/b", true},
		{"../a/b1", "../a/b2", false},
		{"../a1", "../a2", false},
		{"../a/b", "../a", false},
		{"../", "../", true},
		{"../a0", "../a0", true},

		{"/a/b/c", "/a/../a/b/c/./d", true},
	}

	for _, tc := range testCases {
		t.Run(tc.parent, func(t *testing.T) {
			actual, err := IsSub(tc.parent, tc.child)
			if err != nil {
				t.Errorf("test IsSub error, get actual err:%s", err)
			}
			if actual != tc.isSub {
				t.Errorf("test IsSub error, expect get %v but get %v", tc.isSub, actual)
			}
		})
	}
}

func TestIsSub_ReturnError(t *testing.T) {
	abs = absErrorMock
	rel = relErrorMock
	defer func() {
		abs = filepath.Abs
		rel = filepath.Rel
	}()

	testCases := []struct {
		parent string
		child  string
	}{
		{"", "/a"},
		{"/a", ""},
		{"/a", "/a"},
	}

	for _, tc := range testCases {
		t.Run(tc.parent, func(t *testing.T) {
			_, err := IsSub(tc.parent, tc.child)
			if err == nil {
				t.Errorf("test IsSub error, expect to get an error but get nil")
			}
		})
	}
}

func isNotExistAlwaysFalseMock(err error) bool {
	return false
}

// absErrorMock if the path is empty, then returns error
func absErrorMock(path string) (string, error) {
	if len(path) == 0 {
		return "", errors.New("abs error mock")
	}
	return filepath.Abs(path)
}

func relErrorMock(basepath, targpath string) (string, error) {
	return "", errors.New("rel error mock")
}
