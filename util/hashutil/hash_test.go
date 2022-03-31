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

func TestHashValues(t *testing.T) {
	var hvs HashValues
	expect := "461d19e03559ff8a1284951bab8327e1"
	if hvs.Last() != nil {
		t.Errorf("test TestHashValues.Last error, expect get a nil")
	}
	hvs = append(hvs, &HashValue{Offset: 1, Hash: "21cc28409729565fc1a4d2dd92db269f"})
	hvs = append(hvs, &HashValue{Offset: 2, Hash: expect})

	if hvs.Last() == nil {
		t.Errorf("test TestHashValues.Last error, expect:%s, actual get a nil", expect)
		return
	}

	actual := hvs.Last().Hash
	if actual != expect {
		t.Errorf("test TestHashValues.Last error, expect:%s, actual:%s", expect, actual)
	}
}

type readwrite struct {
	*os.File
}

func (rw readwrite) WriteTo(w io.Writer) (n int64, err error) {
	return 0, errors.New("write error test")
}
