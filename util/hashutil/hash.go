package hashutil

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"os"
)

var (
	errNilFile   = errors.New("file is nil")
	errEmptyPath = errors.New("file path can't be empty")
)

// MD5FromFile calculate the hash value of the file
// If you reuse the file reader, please set its offset to start position first, like os.File.Seek
func MD5FromFile(file io.Reader) (hash string, err error) {
	if file == nil {
		return hash, errNilFile
	}
	h := md5.New()
	reader := bufio.NewReader(file)
	_, err = reader.WriteTo(h)
	if err != nil {
		return hash, err
	}
	sum := h.Sum(nil)
	hash = hex.EncodeToString(sum)
	return hash, nil
}

// MD5FromFileName calculate the hash value of the file
func MD5FromFileName(path string) (hash string, err error) {
	if len(path) == 0 {
		return hash, errEmptyPath
	}
	f, err := os.Open(path)
	if err != nil {
		return hash, err
	}
	defer f.Close()
	return MD5FromFile(f)
}

// MD5FromFileChunk calculate the hash value of the file chunk
func MD5FromFileChunk(path string, offset int64, chunkSize int64) (hash string, err error) {
	if len(path) == 0 {
		return hash, errEmptyPath
	}
	f, err := os.Open(path)
	if err != nil {
		return hash, err
	}
	defer f.Close()
	block := make([]byte, chunkSize)
	n, err := f.ReadAt(block, offset)
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		return hash, err
	}
	return MD5(block[:n]), nil
}

// MD5 calculate the hash value of the bytes
func MD5(bytes []byte) (hash string) {
	h := md5.New()
	h.Write(bytes)
	sum := h.Sum(nil)
	hash = hex.EncodeToString(sum)
	return hash
}

// MD5FromString calculate the hash value of the string
func MD5FromString(s string) (hash string) {
	return MD5([]byte(s))
}

// CheckpointsMD5FromFileName calculate the hash value of the full file and first chunk and some checkpoints
// first chunk hash is optional
// checkpoint hash is optional
// full file hash is required
func CheckpointsMD5FromFileName(path string, chunkSize int64, checkpointCount int) (hvs []*HashValue, err error) {
	if len(path) == 0 {
		return nil, errEmptyPath
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return checkpointsMD5FromFile(f, chunkSize, checkpointCount)
}

func checkpointsMD5FromFile(f *os.File, chunkSize int64, checkpointCount int) (hvs []*HashValue, err error) {
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return checkpointsMD5FromFileWithFileSize(f, stat.Size(), chunkSize, checkpointCount)
}

func checkpointsMD5FromFileWithFileSize(f *os.File, fileSize int64, chunkSize int64, checkpointCount int) (hvs []*HashValue, err error) {
	// add first chunk hash
	if fileSize > chunkSize {
		hvs = append(hvs, &HashValue{
			Offset: chunkSize,
		})
	}

	checkpointSize := fileSize / int64(checkpointCount)
	// checkpoint size equals one times or more the chunk size
	if checkpointSize/chunkSize == 0 {
		checkpointSize = chunkSize
	} else {
		checkpointSize = checkpointSize / chunkSize * chunkSize
	}

	// reset the checkpoint count
	checkpointCount = int(fileSize / checkpointSize)

	// add checkpoint hash
	for i := 1; i <= checkpointCount; i++ {
		hvs = append(hvs, &HashValue{
			Offset: checkpointSize * int64(i),
		})
	}

	// add full file hash
	if (len(hvs) > 0 && hvs[len(hvs)-1].Offset < fileSize) || len(hvs) == 0 {
		hvs = append(hvs, &HashValue{
			Offset: fileSize,
		})
	}

	block := make([]byte, chunkSize)
	h := md5.New()

	var writeLen int64
	hvi := 0
	hv := hvs[0]
	isEOF := false
	// calculate hash
	for {
		n, err := f.Read(block)
		if err == io.EOF {
			isEOF = true
			err = nil
		}
		if err != nil {
			return nil, err
		}

		writeLen += int64(n)
		h.Write(block[:n])
		if writeLen >= hv.Offset {
			hv.Offset = writeLen
			hv.Hash = hex.EncodeToString(h.Sum(nil))
			hvi++
			if hvi < len(hvs) {
				hv = hvs[hvi]
			}
		}
		// read to end or all tasks finished
		if isEOF || hvi >= len(hvs) {
			break
		}
	}
	return hvs, nil
}
