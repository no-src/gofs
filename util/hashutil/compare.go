package hashutil

import (
	"encoding/hex"
	"io"
	"os"
	"time"
)

// Compare whether the source file is equal to the destination file
func Compare(chunkSize int64, checkpointCount int, sourceFile *os.File, sourceSize int64, dest string, destSize int64, offset *int64) (equal bool) {
	if destSize <= 0 {
		return false
	}
	hvs, _ := CheckpointsHashFromFile(sourceFile, chunkSize, checkpointCount)
	if len(hvs) > 0 && hvs.Last().Offset == sourceSize {
		// if source and dest is the same file, ignore the following steps and return directly
		equal, hv := CompareHashValues(dest, sourceSize, hvs.Last().Hash, chunkSize, hvs)
		if equal {
			return equal
		}

		if hv != nil {
			*offset = hv.Offset
		}
	}
	return false
}

// QuickCompare if the forceChecksum is false, check whether the size and time are both equal, otherwise return false
func QuickCompare(forceChecksum bool, sourceSize, destSize int64, sourceModTime, destModTime time.Time) (equal bool) {
	if !forceChecksum && sourceSize == destSize && sourceModTime == destModTime {
		return true
	}
	return false
}

// CompareHashValues compare the HashValues from source file with the destination file
func CompareHashValues(dstPath string, sourceSize int64, sourceHash string, chunkSize int64, hvs HashValues) (equal bool, hv *HashValue) {
	if sourceSize > 0 {
		// calculate the entire file hash value
		if len(hvs) == 0 || hvs.Last().Offset < sourceSize {
			hvs = append(hvs, NewHashValue(sourceSize, sourceHash))
		}
		hv, err := CompareHashValuesWithFileName(dstPath, chunkSize, hvs)
		if err == nil && hv != nil {
			return hv.Offset == sourceSize && hv.Hash == sourceHash && len(sourceHash) > 0, hv
		}
	}
	return false, nil
}

// CompareHashValuesWithFileName calculate the file hashes and return the last continuous hit HashValue.
// The offset in the HashValues must equal chunkSize * N, and N greater than zero
func CompareHashValuesWithFileName(path string, chunkSize int64, hvs HashValues) (eq *HashValue, err error) {
	f, err := open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if chunkSize <= 0 {
		return nil, errChunkSizeInvalid
	}
	if len(hvs) == 0 {
		return nil, nil
	}
	h := New()
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
