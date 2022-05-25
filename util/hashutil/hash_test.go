package hashutil

import (
	"errors"
	"io"
	"os"
	"testing"
)

const (
	testFilePath     = "./hash_test.go"
	testDirPath      = "./"
	notExistFilePath = "./not_exist.txt"
)

func TestHashFromFile_ReturnError(t *testing.T) {
	testCases := []struct {
		name   string
		reader io.Reader
	}{
		{"nil reader", nil},
		{"always write to error", readwrite{nil}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := HashFromFile(tc.reader); err == nil {
				t.Errorf("test HashFromFile error, expect to get an error but get nil")
			}
		})
	}
}

func TestHashFromFileName(t *testing.T) {
	testCases := []struct {
		path string
	}{
		{testFilePath},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			if _, err := HashFromFileName(tc.path); err != nil {
				t.Errorf("test HashFromFileName error => %s", err)
			}
		})
	}
}

func TestHashFromFileName_ReturnError(t *testing.T) {
	testCases := []struct {
		name string
		path string
	}{
		{"not exist file path", notExistFilePath},
		{"empty path", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := HashFromFileName(tc.path); err == nil {
				t.Errorf("test HashFromFileName error, expect to get an error but get nil")
			}
		})
	}
}

func TestHashFromString(t *testing.T) {
	testCases := []struct {
		str    string
		expect string
	}{
		{"", "d41d8cd98f00b204e9800998ecf8427e"},
		{"golang", "21cc28409729565fc1a4d2dd92db269f"},
		{"hello test", "7a6d667ea5ed4467c017b2ed6ea07e78"},
	}

	for _, tc := range testCases {
		t.Run("["+tc.str+"]", func(t *testing.T) {
			actual := HashFromString(tc.str)
			if actual != tc.expect {
				t.Errorf("test HashFromString error, expect:%s, actual:%s", tc.expect, actual)
			}
		})
	}
}

func TestHashFromFileChunk(t *testing.T) {
	testCases := []struct {
		name      string
		path      string
		offset    int64
		chunkSize int64
	}{
		{"normal", testFilePath, 10, 100},
		{"with read to end", testFilePath, 1024 * 1024 * 10, 1000},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := HashFromFileChunk(tc.path, tc.offset, tc.chunkSize); err != nil {
				t.Errorf("test HashFromFileChunk error => %s", err)
			}
		})
	}
}

func TestHashFromFileChunk_ReturnError(t *testing.T) {
	testCases := []struct {
		name      string
		path      string
		offset    int64
		chunkSize int64
	}{
		{"with not exist file", notExistFilePath, 10, 100},
		{"with empty path", "", 10, 100},
		{"with invalid offset", testFilePath, -1, 100},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := HashFromFileChunk(tc.path, tc.offset, tc.chunkSize); err == nil {
				t.Errorf("test HashFromFileChunk error, expect to get an error but get nil")
			}
		})
	}
}

func TestCheckpointsHashFromFileName_ReturnError(t *testing.T) {
	var chunkSize int64 = 20
	checkpointCount := 10

	testCases := []struct {
		name            string
		path            string
		chunkSize       int64
		checkpointCount int
	}{
		{"with empty path", "", chunkSize, checkpointCount},
		{"with not exist file path", notExistFilePath, chunkSize, checkpointCount},
		{"with dir path", testDirPath, chunkSize, checkpointCount},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := CheckpointsHashFromFileName(tc.path, tc.chunkSize, tc.checkpointCount); err == nil {
				t.Errorf("test CheckpointsHashFromFileName error, expect to get an error but get nil")
			}
		})
	}
}

func TestCheckpointsHashFromFile_ReturnError(t *testing.T) {
	var chunkSize int64 = 20
	checkpointCount := 10

	testCases := []struct {
		name            string
		f               *os.File
		chunkSize       int64
		checkpointCount int
	}{
		{"with nil *os.File", nil, chunkSize, checkpointCount},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := CheckpointsHashFromFile(tc.f, tc.chunkSize, tc.checkpointCount); err == nil {
				t.Errorf("test checkpointsHashFromFile error, expect to get an error but get nil")
			}
		})
	}
}

func TestCheckpointsHashFromFileName(t *testing.T) {
	checkpointCount := 10
	path := testFilePath
	hash, err := HashFromFileName(path)
	if err != nil {
		t.Errorf("test HashFromFileName error => %s", err)
		return
	}

	testCases := []struct {
		name            string
		path            string
		chunkSize       int64
		checkpointCount int
		expect          string
	}{
		{"", path, 20, checkpointCount, hash},
		{"", path, 20, 0, hash},
		{"", path, 1024, checkpointCount, hash},
		{"", path, 1024, 0, hash},
		{"", path, 0, checkpointCount, hash},
		{"", path, 0, 0, hash},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCheckpointsHashFromFileName(t, tc.path, tc.chunkSize, tc.checkpointCount, tc.expect)
		})
	}
}

func testCheckpointsHashFromFileName(t *testing.T, path string, chunkSize int64, checkpointCount int, expect string) {
	hvs, err := CheckpointsHashFromFileName(path, chunkSize, checkpointCount)
	if err != nil {
		t.Errorf("test TestCheckpointsHashFromFileName error chunkSize=%d checkpointCount=%d => %s", chunkSize, checkpointCount, err)
	}

	if len(hvs) == 0 {
		t.Errorf("test TestCheckpointsHashFromFileName error chunkSize=%d checkpointCount=%d expect:%s, actual:nothing", chunkSize, checkpointCount, expect)
	} else if hvs.Last().Hash != expect {
		t.Errorf("test TestCheckpointsHashFromFileName error chunkSize=%d checkpointCount=%d expect:%s, actual:%s", chunkSize, checkpointCount, expect, hvs.Last().Hash)
	}
}

func TestCalcHashValuesWithFile(t *testing.T) {
	var hvs HashValues
	testCases := []struct {
		name      string
		f         *os.File
		chunkSize int64
		hvs       HashValues
	}{
		{"with empty HashValues", nil, defaultChunkSize, hvs},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := calcHashValuesWithFile(tc.f, tc.chunkSize, tc.hvs); err != nil {
				t.Errorf("test calcHashValuesWithFile error => %v", err)
			}
		})
	}
}

func TestCalcHashValuesWithFile_ReturnError(t *testing.T) {
	var hvs HashValues
	testCases := []struct {
		name      string
		f         *os.File
		chunkSize int64
		hvs       HashValues
	}{
		{"with zero chunk size", nil, 0, hvs},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := calcHashValuesWithFile(tc.f, tc.chunkSize, tc.hvs); err == nil {
				t.Errorf("test calcHashValuesWithFile error, expect to get an error but get nil")
			}
		})
	}
}

func TestCompareHashValuesWithFileName_ReturnError(t *testing.T) {
	var hvs HashValues
	testCases := []struct {
		name      string
		path      string
		chunkSize int64
		hvs       HashValues
	}{
		{"with empty path", "", defaultChunkSize, hvs},
		{"with not exist file path", notExistFilePath, defaultChunkSize, hvs},
		{"with zero chunk size", testFilePath, 0, hvs},
		{"with dir path", testDirPath, defaultChunkSize, append(HashValues{}, NewHashValue(2, "e529a9cea4a728eb9c5828b13b22844c"))},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := CompareHashValuesWithFileName(tc.path, tc.chunkSize, tc.hvs); err == nil {
				t.Errorf("test CompareHashValuesWithFileName error, expect to get an error but get nil")
			}
		})
	}
}

func TestCompareHashValuesWithFileName(t *testing.T) {
	path := testFilePath
	var chunkSize int64 = 1
	var hvs HashValues
	invalidHash := "815417267f76f6f460a4a61f9db75fdb"

	testCases := []struct {
		name      string
		path      string
		chunkSize int64
		hvs       HashValues
		expect    *HashValue
	}{
		{"", path, chunkSize, hvs, nil},
		{"", path, chunkSize, append(HashValues{}, NewHashValue(1, invalidHash)), nil},
		{"", path, chunkSize, append(HashValues{}, NewHashValue(2, "e529a9cea4a728eb9c5828b13b22844c")), NewHashValue(2, "e529a9cea4a728eb9c5828b13b22844c")},
		{"", path, chunkSize, append(HashValues{}, NewHashValue(2, "e529a9cea4a728eb9c5828b13b22844c"), NewHashValue(7, "efe90a8e604a7c840e88d03a67f6b7d8")), NewHashValue(7, "efe90a8e604a7c840e88d03a67f6b7d8")},
		{"", path, chunkSize, append(HashValues{}, NewHashValue(1024*1024, invalidHash)), nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCompareHashValuesWithFileName(t, tc.path, tc.chunkSize, tc.hvs, tc.expect)
		})
	}
}

func testCompareHashValuesWithFileName(t *testing.T, path string, chunkSize int64, hvs HashValues, expect *HashValue) {
	hv, err := CompareHashValuesWithFileName(path, chunkSize, hvs)
	if err != nil {
		t.Errorf("test CompareHashValuesWithFileName error %v", err)
		return
	}

	if expect == nil {
		if hv != nil {
			t.Errorf("test CompareHashValuesWithFileName error, expect to get an nil HashValue")
		}
		return
	}

	if hv == nil {
		t.Errorf("test CompareHashValuesWithFileName error, get an nil HashValue, expect:[%d => %s]", expect.Offset, expect.Hash)
		return
	}

	if hv.Offset != expect.Offset || hv.Hash != expect.Hash {
		t.Errorf("test CompareHashValuesWithFileName expect:[%d => %s] actual:[%d => %s]", expect.Offset, expect.Hash, hv.Offset, hv.Hash)
	}
}

type readwrite struct {
	*os.File
}

func (rw readwrite) WriteTo(w io.Writer) (n int64, err error) {
	return 0, errors.New("write error test")
}
