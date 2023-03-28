package hashutil

import (
	"bufio"
	"encoding/hex"
	"errors"
	"hash"
	"io"
	"os"
	"time"
)

var (
	errNilFile          = errors.New("file is nil")
	errEmptyPath        = errors.New("file path can't be empty")
	errChunkSizeInvalid = errors.New("chunk size must be greater than zero")
)

const (
	defaultChunkSize = 4096
)

// Hash a hash calculate component
type Hash interface {
	// HashFromFile calculate the hash value of the file
	// If you reuse the file reader, please set its offset to start position first, like os.File.Seek
	HashFromFile(file io.Reader) (hashString string, err error)
	// HashFromFileName calculate the hash value of the file
	HashFromFileName(path string) (hash string, err error)
	// HashFromFileChunk calculate the hash value of the file chunk
	HashFromFileChunk(path string, offset int64, chunkSize int64) (hash string, err error)
	// Hash calculate the hash value of the bytes
	Hash(bytes []byte) (hashString string)
	// HashFromString calculate the hash value of the string
	HashFromString(s string) (hash string)
	// CheckpointsHashFromFileName calculate the hash value of the entire file and first chunk and some checkpoints
	// the first chunk hash is optional
	// the checkpoint hash is optional
	// the entire file hash is required
	CheckpointsHashFromFileName(path string, chunkSize int64, checkpointCount int) (hvs HashValues, err error)
	// CheckpointsHashFromFile calculate the hash value of the entire file and first chunk and some checkpoints
	CheckpointsHashFromFile(f *os.File, chunkSize int64, checkpointCount int) (hvs HashValues, err error)
	// GetFileSizeAndHashCheckpoints get the file size and hash checkpoints from the specified file
	GetFileSizeAndHashCheckpoints(path string, chunkSize int64, checkpointCount int) (size int64, hash string, hvs HashValues, err error)
	// Compare whether the source file is equal to the destination file
	Compare(chunkSize int64, checkpointCount int, sourceFile *os.File, sourceSize int64, dest string, destSize int64, offset *int64) (equal bool)
	// QuickCompare if the forceChecksum is false, check whether the size and time are both equal, otherwise return false
	QuickCompare(forceChecksum bool, sourceSize, destSize int64, sourceModTime, destModTime time.Time) (equal bool)
	// CompareHashValues compare the HashValues from source file with the destination file
	CompareHashValues(dstPath string, sourceSize int64, sourceHash string, chunkSize int64, hvs HashValues) (equal bool, hv *HashValue)
	// CompareHashValuesWithFileName calculate the file hashes and return the last continuous hit HashValue.
	// The offset in the HashValues must equal chunkSize * N, and N greater than zero
	CompareHashValuesWithFileName(path string, chunkSize int64, hvs HashValues) (eq *HashValue, err error)
}

type defaultHash struct {
	factory hashFactory
}

func (dh *defaultHash) new() hash.Hash {
	return dh.factory()
}

func (dh *defaultHash) HashFromFile(file io.Reader) (hashString string, err error) {
	if file == nil {
		return hashString, errNilFile
	}
	hash := dh.new()
	reader := bufio.NewReader(file)
	_, err = reader.WriteTo(hash)
	if err != nil {
		return hashString, err
	}
	sum := hash.Sum(nil)
	hashString = hex.EncodeToString(sum)
	return hashString, nil
}

func (dh *defaultHash) HashFromFileName(path string) (hash string, err error) {
	f, err := dh.open(path)
	if err != nil {
		return hash, err
	}
	defer f.Close()
	return dh.HashFromFile(f)
}

func (dh *defaultHash) HashFromFileChunk(path string, offset int64, chunkSize int64) (hash string, err error) {
	f, err := dh.open(path)
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
	return dh.Hash(chunk[:n]), nil
}

func (dh *defaultHash) Hash(bytes []byte) (hashString string) {
	hash := dh.new()
	hash.Write(bytes)
	sum := hash.Sum(nil)
	hashString = hex.EncodeToString(sum)
	return hashString
}

func (dh *defaultHash) HashFromString(s string) (hash string) {
	return dh.Hash([]byte(s))
}

func (dh *defaultHash) CheckpointsHashFromFileName(path string, chunkSize int64, checkpointCount int) (hvs HashValues, err error) {
	f, err := dh.open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return dh.CheckpointsHashFromFile(f, chunkSize, checkpointCount)
}

func (dh *defaultHash) CheckpointsHashFromFile(f *os.File, chunkSize int64, checkpointCount int) (hvs HashValues, err error) {
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := stat.Size()

	// add first chunk hash
	if chunkSize > 0 && fileSize > chunkSize {
		hvs = append(hvs, NewHashValue(chunkSize, ""))
	}

	// init default chunk size
	if chunkSize == 0 {
		chunkSize = defaultChunkSize
	}

	// add checkpoint hash
	hvs = append(hvs, dh.buildCheckpointHashValues(fileSize, chunkSize, checkpointCount)...)

	// add entire file hash
	if (len(hvs) > 0 && hvs.Last().Offset < fileSize) || len(hvs) == 0 {
		hvs = append(hvs, NewHashValue(fileSize, ""))
	}

	err = dh.calcHashValuesWithFile(f, chunkSize, hvs)
	return hvs, err
}

func (dh *defaultHash) buildCheckpointHashValues(fileSize int64, chunkSize int64, checkpointCount int) (hvs HashValues) {
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

func (dh *defaultHash) calcHashValuesWithFile(f *os.File, chunkSize int64, hvs HashValues) error {
	if chunkSize <= 0 {
		return errChunkSizeInvalid
	}
	if len(hvs) == 0 {
		return nil
	}
	h := dh.new()
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

func (dh *defaultHash) open(path string) (*os.File, error) {
	if len(path) == 0 {
		return nil, errEmptyPath
	}
	return os.Open(path)
}

func (dh *defaultHash) GetFileSizeAndHashCheckpoints(path string, chunkSize int64, checkpointCount int) (size int64, hash string, hvs HashValues, err error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return size, hash, hvs, err
	}
	if fileInfo.IsDir() {
		return size, hash, hvs, nil
	}
	size = fileInfo.Size()
	if size > 0 {
		hvs, err = dh.CheckpointsHashFromFileName(path, chunkSize, checkpointCount)
		if err == nil && len(hvs) > 0 {
			hash = hvs.Last().Hash
		}
	}
	return size, hash, hvs, err
}
