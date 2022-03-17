package hashutil

import (
	"errors"
	"io"
	"os"
	"testing"
)

const testFilePath = "./hash_test.go"

func TestMD5FromFileError(t *testing.T) {
	_, err := MD5FromFile(nil)
	if err == nil {
		t.Errorf("test MD5FromFile error, should get an error")
	}

	f, err := os.Open(testFilePath)
	if err != nil {
		t.Errorf("test MD5FromFile error => %s", err)
		return
	}
	_, err = MD5FromFile(readwrite{f})
	if err == nil {
		t.Errorf("test MD5FromFile error, should get an error")
	}
}

func TestMD5FromFileName(t *testing.T) {
	_, err := MD5FromFileName(testFilePath)
	if err != nil {
		t.Errorf("test MD5FromFileName error => %s", err)
	}
}

func TestMD5FromFileNameError(t *testing.T) {
	_, err := MD5FromFileName("./not_exist.txt")
	if err == nil {
		t.Errorf("test MD5FromFileName error, should get an error")
	}

	_, err = MD5FromFileName("")
	if err == nil {
		t.Errorf("test MD5FromFileName error, should get an error")
	}
}

func TestMD5(t *testing.T) {
	testMD5(t, "", "d41d8cd98f00b204e9800998ecf8427e")
	testMD5(t, "golang", "21cc28409729565fc1a4d2dd92db269f")
	testMD5(t, "hello test", "7a6d667ea5ed4467c017b2ed6ea07e78")
}

func testMD5(t *testing.T, str, expect string) {
	actual := MD5(str)
	if actual != expect {
		t.Errorf("test MD5 error, expect:%s, actual:%s", expect, actual)
	}
}

type readwrite struct {
	*os.File
}

func (rw readwrite) WriteTo(w io.Writer) (n int64, err error) {
	return 0, errors.New("write error test")
}
