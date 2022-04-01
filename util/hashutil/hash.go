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
	errNilFile          = errors.New("file is nil")
	errEmptyPath        = errors.New("file path can't be empty")
	errChunkSizeInvalid = errors.New("chunk size must be greater than zero")
)

const (
	defaultChunkSize = 4096
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
	chunk := make([]byte, chunkSize)
	n, err := f.ReadAt(chunk, offset)
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		return hash, err
	}
	return MD5(chunk[:n]), nil
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

// CheckpointsMD5FromFileName calculate the hash value of the entire file and first chunk and some checkpoints
// the first chunk hash is optional
// the checkpoint hash is optional
// the entire file hash is required
func CheckpointsMD5FromFileName(path string, chunkSize int64, checkpointCount int) (hvs HashValues, err error) {
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

func checkpointsMD5FromFile(f *os.File, chunkSize int64, checkpointCount int) (hvs HashValues, err error) {
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return checkpointsMD5FromFileWithFileSize(f, stat.Size(), chunkSize, checkpointCount)
}

func checkpointsMD5FromFileWithFileSize(f *os.File, fileSize int64, chunkSize int64, checkpointCount int) (hvs HashValues, err error) {
	// add first chunk hash
	if chunkSize > 0 && fileSize > chunkSize {
		hvs = append(hvs, NewHashValue(chunkSize, ""))
	}

	// init default chunk size
	if chunkSize == 0 {
		chunkSize = defaultChunkSize
	}

	// add checkpoint hash
	hvs = append(hvs, buildCheckpointHashValues(fileSize, chunkSize, checkpointCount)...)

	// add entire file hash
	if (len(hvs) > 0 && hvs.Last().Offset < fileSize) || len(hvs) == 0 {
		hvs = append(hvs, NewHashValue(fileSize, ""))
	}

	err = calcHashValuesWithFile(f, chunkSize, hvs)
	return hvs, err
}

func buildCheckpointHashValues(fileSize int64, chunkSize int64, checkpointCount int) (hvs HashValues) {
	if checkpointCount > 0 {
		checkpointSize := fileSize / int64(checkpointCount)

		// use chunk size to reset checkpoint size and count
		if chunkSize > 0 {
			// checkpoint size equals one times or more the chunk size
			if checkpointSize/chunkSize == 0 {
				checkpointSize = chunkSize
			} else {
				checkpointSize = checkpointSize / chunkSize * chunkSize
			}
			// reset the checkpoint count
			checkpointCount = int(fileSize / checkpointSize)
		}

		// add checkpoint hash
		for i := 1; i <= checkpointCount; i++ {
			hvs = append(hvs, NewHashValue(checkpointSize*int64(i), ""))
		}
	}
	return hvs
}

func calcHashValuesWithFile(f *os.File, chunkSize int64, hvs HashValues) error {
	if chunkSize <= 0 {
		return errChunkSizeInvalid
	}
	if len(hvs) == 0 {
		return nil
	}
	h := md5.New()
	var writeLen int64
	hvi := 0
	hv := hvs[0]
	isEOF := false
	chunk := make([]byte, chunkSize)
	// calculate hash
	for {
		n, err := f.Read(chunk)
		if err == io.EOF {
			isEOF = true
			err = nil
		}
		if err != nil {
			return err
		}

		writeLen += int64(n)
		h.Write(chunk[:n])
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
	return nil
}

// CompareHashValuesWithFileName calculate the file hashes and return the last continuous hit HashValue.
// The offset in the HashValues must equal chunkSize * N, and N greater than zero
func CompareHashValuesWithFileName(path string, chunkSize int64, hvs HashValues) (*HashValue, error) {
	if len(path) == 0 {
		return nil, errEmptyPath
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return compareHashValuesWithFile(f, chunkSize, hvs)
}

func compareHashValuesWithFile(f *os.File, chunkSize int64, hvs HashValues) (eq *HashValue, err error) {
	if chunkSize <= 0 {
		return nil, errChunkSizeInvalid
	}
	if len(hvs) == 0 {
		return nil, nil
	}
	h := md5.New()
	var writeLen int64
	hvi := 0
	hv := hvs[0]
	isEOF := false
	chunk := make([]byte, chunkSize)
	// calculate hash
	for {
		n, err := f.Read(chunk)
		if err == io.EOF {
			isEOF = true
			err = nil
		}
		if err != nil {
			return nil, err
		}

		writeLen += int64(n)
		h.Write(chunk[:n])
		if writeLen >= hv.Offset {
			if writeLen != hv.Offset || hv.Hash != hex.EncodeToString(h.Sum(nil)) {
				return eq, nil
			}
			eq = hv
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
	return eq, nil
}
