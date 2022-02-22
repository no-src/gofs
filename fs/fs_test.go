package fs

import (
	"github.com/no-src/gofs/util"
	"testing"
)

func TestGetFileTime(t *testing.T) {
	file := "./fs_test.go"
	_, _, _, err := GetFileTime(file)
	if err != nil {
		t.Errorf("get file time error %s => %v", file, err)
		return
	}

	file = "./fs_test_not_found.go"
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
	file := "./fs_test.go"
	exist, err := FileExist(file)
	if err != nil {
		t.Errorf("check file exist error %s => %v", file, err)
		return
	}
	if !exist {
		t.Errorf("check file exist error, file should be exist => %s", file)
		return
	}

	file = "./fs_test_not_found.go"
	exist, err = FileExist(file)
	if err != nil {
		t.Errorf("check file exist error %s => %v", file, err)
		return
	}

	if exist {
		t.Errorf("check file exist error, file should be not exist => %s", file)
		return
	}

	if util.IsWindows() {
		file = "|/"
		exist, err = FileExist(file)
		if err == nil {
			t.Errorf("check an invalid file should be return error %s => %v", file, err)
			return
		}
	}
}

func TestCreateFile(t *testing.T) {
	file := "./fs_test.go"
	_, err := CreateFile(file)
	if err != nil {
		t.Errorf("create file error %s => %v", file, err)
		return
	}
}

func TestOpenRWFile(t *testing.T) {
	file := "./fs_test.go"
	_, err := OpenRWFile(file)
	if err != nil {
		t.Errorf("create file error %s => %v", file, err)
		return
	}
}

func TestIsDir(t *testing.T) {
	file := "./fs_test.go"
	_, err := IsDir(file)
	if err != nil {
		t.Errorf("check path is dir error %s => %v", file, err)
		return
	}

	file = "./fs_test_not_found.go"
	_, err = IsDir(file)
	if err == nil {
		t.Errorf("check path is dir from a not exist file should be return error %s => %v", file, err)
		return
	}
}
