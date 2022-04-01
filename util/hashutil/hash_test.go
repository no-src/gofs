package hashutil

import (
	"errors"
	"io"
	"os"
	"testing"
)

const testFilePath = "./hash_test.go"
const notExistFilePath = "./not_exist.txt"

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
	_, err := MD5FromFileName(notExistFilePath)
	if err == nil {
		t.Errorf("test MD5FromFileName with not exist file error, should get an error")
	}

	_, err = MD5FromFileName("")
	if err == nil {
		t.Errorf("test MD5FromFileName with empty path error, should get an error")
	}
}

func TestMD5FromString(t *testing.T) {
	testMD5FromString(t, "", "d41d8cd98f00b204e9800998ecf8427e")
	testMD5FromString(t, "golang", "21cc28409729565fc1a4d2dd92db269f")
	testMD5FromString(t, "hello test", "7a6d667ea5ed4467c017b2ed6ea07e78")
}

func testMD5FromString(t *testing.T, str, expect string) {
	actual := MD5FromString(str)
	if actual != expect {
		t.Errorf("test MD5 error, expect:%s, actual:%s", expect, actual)
	}
}

func TestMD5FromFileChunk(t *testing.T) {
	_, err := MD5FromFileChunk(testFilePath, 10, 100)
	if err != nil {
		t.Errorf("test MD5FromFileChunk error => %s", err)
	}

	_, err = MD5FromFileChunk(testFilePath, 1024*1024*10, 1000)
	if err != nil {
		t.Errorf("test MD5FromFileChunk with read to end error => %s", err)
	}
}

func TestMD5FromFileChunkError(t *testing.T) {
	_, err := MD5FromFileChunk(notExistFilePath, 10, 100)
	if err == nil {
		t.Errorf("test MD5FromFileChunk with not exist file error, should get an error")
	}

	_, err = MD5FromFileChunk("", 10, 100)
	if err == nil {
		t.Errorf("test MD5FromFileChunk with empty path error, should get an error")
	}

	_, err = MD5FromFileChunk(testFilePath, -1, 100)
	if err == nil {
		t.Errorf("test MD5FromFileChunk with invalid offset error, should get an error")
	}
}

func TestCheckpointsMD5FromFileNameError(t *testing.T) {
	var chunkSize int64 = 20
	checkpointCount := 10

	_, err := CheckpointsMD5FromFileName("", chunkSize, checkpointCount)
	if err == nil {
		t.Errorf("test TestCheckpointsMD5FromFileName with empty path error, expect get an error")
	}

	_, err = CheckpointsMD5FromFileName(notExistFilePath, chunkSize, checkpointCount)
	if err == nil {
		t.Errorf("test TestCheckpointsMD5FromFileName with not exist file path error, expect get an error")
	}

	_, err = checkpointsMD5FromFile(nil, chunkSize, checkpointCount)
	if err == nil {
		t.Errorf("test checkpointsMD5FromFile with nil *os.File error, expect get an error")
	}

	_, err = checkpointsMD5FromFileWithFileSize(nil, 0, chunkSize, checkpointCount)
	if err == nil {
		t.Errorf("test checkpointsMD5FromFileWithFileSize with nil *os.File error, expect get an error")
	}
}

func TestCheckpointsMD5FromFileName(t *testing.T) {
	var chunkSize int64 = 20
	checkpointCount := 10
	path := testFilePath
	hash, err := MD5FromFileName(path)
	if err != nil {
		t.Errorf("test MD5FromFileName error => %s", err)
		return
	}

	testCheckpointsMD5FromFileName(t, path, chunkSize, checkpointCount, hash)
	testCheckpointsMD5FromFileName(t, path, chunkSize, 0, hash)

	chunkSize = 1024
	testCheckpointsMD5FromFileName(t, path, chunkSize, checkpointCount, hash)
	testCheckpointsMD5FromFileName(t, path, chunkSize, 0, hash)
	testCheckpointsMD5FromFileName(t, path, 0, checkpointCount, hash)
	testCheckpointsMD5FromFileName(t, path, 0, 0, hash)
}

func testCheckpointsMD5FromFileName(t *testing.T, path string, chunkSize int64, checkpointCount int, expect string) {
	hvs, err := CheckpointsMD5FromFileName(path, chunkSize, checkpointCount)
	if err != nil {
		t.Errorf("test TestCheckpointsMD5FromFileName error chunkSize=%d checkpointCount=%d => %s", chunkSize, checkpointCount, err)
	}

	if len(hvs) == 0 {
		t.Errorf("test TestCheckpointsMD5FromFileName error chunkSize=%d checkpointCount=%d expect:%s, actual:nothing", chunkSize, checkpointCount, expect)
	} else if hvs.Last().Hash != expect {
		t.Errorf("test TestCheckpointsMD5FromFileName error chunkSize=%d checkpointCount=%d expect:%s, actual:%s", chunkSize, checkpointCount, expect, hvs.Last().Hash)
	}
}

func TestCalcHashValuesWithFile(t *testing.T) {
	var hvs HashValues
	err := calcHashValuesWithFile(nil, 0, hvs)
	if err == nil {
		t.Errorf("test calcHashValuesWithFile with zero chunk size error, expect get an error")
	}

	err = calcHashValuesWithFile(nil, defaultChunkSize, hvs)
	if err != nil {
		t.Errorf("test calcHashValuesWithFile with empty HashValues error, expect get an nil, actual:%v", err)
	}
}

func TestCompareHashValuesWithFileNameError(t *testing.T) {
	var hvs HashValues
	_, err := CompareHashValuesWithFileName("", defaultChunkSize, hvs)
	if err == nil {
		t.Errorf("test CompareHashValuesWithFileName with empty path error, expect get an error")
	}

	_, err = CompareHashValuesWithFileName(notExistFilePath, defaultChunkSize, hvs)
	if err == nil {
		t.Errorf("test CompareHashValuesWithFileName with not exist file path error, expect get an error")
	}

	_, err = CompareHashValuesWithFileName(testFilePath, 0, hvs)
	if err == nil {
		t.Errorf("test CompareHashValuesWithFileName with zero chunk size error, expect get an error")
	}

	hvs = append(hvs, NewHashValue(2, "e529a9cea4a728eb9c5828b13b22844c"))
	_, err = compareHashValuesWithFile(nil, defaultChunkSize, hvs)
	if err == nil {
		t.Errorf("test CompareHashValuesWithFileName with nil *os.File error, expect get an error")
	}
}

func TestCompareHashValuesWithFileName(t *testing.T) {
	path := testFilePath
	var chunkSize int64 = 1
	var hvs HashValues
	invalidHash := "815417267f76f6f460a4a61f9db75fdb"

	testCompareHashValuesWithFileName(t, path, chunkSize, hvs, nil)

	hvs = append(hvs, NewHashValue(1, invalidHash))
	testCompareHashValuesWithFileName(t, path, chunkSize, hvs, nil)

	hvs = make(HashValues, 0)
	hvs = append(hvs, NewHashValue(2, "e529a9cea4a728eb9c5828b13b22844c"))
	testCompareHashValuesWithFileName(t, path, chunkSize, hvs, NewHashValue(2, "e529a9cea4a728eb9c5828b13b22844c"))

	hvs = append(hvs, NewHashValue(7, "efe90a8e604a7c840e88d03a67f6b7d8"))
	testCompareHashValuesWithFileName(t, path, chunkSize, hvs, NewHashValue(7, "efe90a8e604a7c840e88d03a67f6b7d8"))

	hvs = make(HashValues, 0)
	hvs = append(hvs, NewHashValue(1024*1024, invalidHash))
	testCompareHashValuesWithFileName(t, path, chunkSize, hvs, nil)
}

func testCompareHashValuesWithFileName(t *testing.T, path string, chunkSize int64, hvs HashValues, expect *HashValue) {
	hv, err := CompareHashValuesWithFileName(path, chunkSize, hvs)
	if err != nil {
		t.Errorf("test CompareHashValuesWithFileName error %v", err)
		return
	}

	if expect == nil {
		if hv != nil {
			t.Errorf("test CompareHashValuesWithFileName error, expect get an nil HashValue")
		}
		return
	}

	if hv == nil {
		t.Errorf("test CompareHashValuesWithFileName error, get an nil HashValue, expact:[%d => %s]", expect.Offset, expect.Hash)
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
