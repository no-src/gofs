package rate

import (
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestNewHTTPDir(t *testing.T) {
	bytesPerSecond := KB
	dir := NewHTTPDir("./", bytesPerSecond)
	f, err := dir.Open("fs_test.go")
	if err != nil {
		t.Errorf("open file error, %v", err)
		return
	}
	stat, err := f.Stat()
	if err != nil {
		t.Errorf("get file state error, %v", err)
		return
	}
	start := time.Now()
	_, err = io.ReadAll(f)
	end := time.Now()
	if err != nil {
		t.Errorf("read file error")
		return
	}

	dataSize := stat.Size()
	actualCost := end.Sub(start)
	t.Logf("dataSize=%d bytesPerSecond=%d actualCost=%s", dataSize, bytesPerSecond, actualCost)

	_, err = dir.Open("not_exist_file.go")
	if !os.IsNotExist(err) {
		t.Errorf("expect to get ErrNotExist error, actual get %v", err)
	}
}

func TestNewHTTPDir_DisableOrEnableRate(t *testing.T) {
	testCases := []struct {
		name           string
		bytesPerSecond int64
		expectHTTPDir  bool
	}{
		{"disable rate by zero rate", 0, true},
		{"disable rate by negative rate", -1, true},
		{"enable rate", 1, false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := NewHTTPDir("./", tc.bytesPerSecond)
			switch d.(type) {
			case http.Dir:
				if !tc.expectHTTPDir {
					t.Errorf("expect to get *dir type, actual get http.Dir")
				}
			case *dir:
				if tc.expectHTTPDir {
					t.Errorf("expect to get http.Dir type, actual get *dir")
				}
			default:
				t.Errorf("unexpected type")
			}
		})
	}
}
