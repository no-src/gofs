package fs

import (
	"github.com/no-src/gofs/util/osutil"
	"os"
	"testing"
)

const (
	testNotFoundFilePath = "./fs_test_not_found.go"
	testExistFilePath    = "./fs_test.go"
)

func TestGetFileTime(t *testing.T) {
	file := testExistFilePath
	_, _, _, err := GetFileTime(file)
	if err != nil {
		t.Errorf("get file time error %s => %v", file, err)
		return
	}

	file = testNotFoundFilePath
	_, _, _, err = GetFileTime(file)
	if err == nil {
		t.Errorf("get file time from a not exist file should be return error %s => %v", file, err)
		return
	}
}

func TestGetFileTimeBySys(t *testing.T) {
	_, _, _, err := GetFileTimeBySys(nil)
	if err == nil {
		t.Errorf("GetFileTimeBySys with a nil value should be return error")
		return
	}
}

func TestFileExist(t *testing.T) {
	file := testExistFilePath
	exist, err := FileExist(file)
	if err != nil {
		t.Errorf("check file exist error %s => %v", file, err)
		return
	}
	if !exist {
		t.Errorf("check file exist error, file should be exist => %s", file)
		return
	}

	file = testNotFoundFilePath
	exist, err = FileExist(file)
	if err != nil {
		t.Errorf("check file exist error %s => %v", file, err)
		return
	}

	if exist {
		t.Errorf("check file exist error, file should be not exist => %s", file)
		return
	}

	if !osutil.IsWindows() {
		isNotExist = isNotExistAlwaysFalseMock
		defer func() {
			isNotExist = os.IsNotExist
		}()
	}
	file = "|/"
	_, err = FileExist(file)
	if err == nil {
		t.Errorf("check an invalid file should be return error %s => %v", file, err)
		return
	}
}

func TestCreateFile(t *testing.T) {
	file := testExistFilePath
	_, err := CreateFile(file)
	if err != nil {
		t.Errorf("create file error %s => %v", file, err)
		return
	}
}

func TestOpenRWFile(t *testing.T) {
	file := testExistFilePath
	_, err := OpenRWFile(file)
	if err != nil {
		t.Errorf("create file error %s => %v", file, err)
		return
	}
}

func TestIsDir(t *testing.T) {
	file := testExistFilePath
	_, err := IsDir(file)
	if err != nil {
		t.Errorf("check path is dir error %s => %v", file, err)
		return
	}

	file = testNotFoundFilePath
	_, err = IsDir(file)
	if err == nil {
		t.Errorf("check path is dir from a not exist file should be return error %s => %v", file, err)
		return
	}
}

func TestIsEOF(t *testing.T) {
	file := testExistFilePath
	f, err := os.Open(file)
	if err != nil {
		t.Errorf("test IsEOF error, open file error [%s] => %s", file, err)
		return
	}
	// move to end
	_, err = f.Seek(1, 2)
	if err != nil {
		t.Errorf("test IsEOF error, seek file error [%s] => %s", file, err)
		return
	}
	data := make([]byte, 1024)
	_, err = f.Read(data)
	if !IsEOF(err) {
		t.Errorf("test IsEOF error, read file error [%s] => %s", file, err)
	}
}

func TestIsNonEOF(t *testing.T) {
	file := testNotFoundFilePath
	_, err := os.Stat(file)
	if !IsNonEOF(err) {
		t.Errorf("test IsNonEOF error, get actual err:%s", err)
	}
}

func isNotExistAlwaysFalseMock(err error) bool {
	return false
}
