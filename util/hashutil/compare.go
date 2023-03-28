package hashutil

import (
	"encoding/hex"
	"io"
	"os"
	"time"
)

func (dh *defaultHash) Compare(chunkSize int64, checkpointCount int, sourceFile *os.File, sourceSize int64, dest string, destSize int64, offset *int64) (equal bool) {
	if destSize <= 0 {
		return false
	}
	hvs, _ := dh.CheckpointsHashFromFile(sourceFile, chunkSize, checkpointCount)
	if len(hvs) > 0 && hvs.Last().Offset == sourceSize {
		// if source and dest is the same file, ignore the following steps and return directly
		equal, hv := dh.CompareHashValues(dest, sourceSize, hvs.Last().Hash, chunkSize, hvs)
		if equal {
			return equal
		}

		if hv != nil {
			*offset = hv.Offset
		}
	}
	return false
}

func (dh *defaultHash) QuickCompare(forceChecksum bool, sourceSize, destSize int64, sourceModTime, destModTime time.Time) (equal bool) {
	if !forceChecksum && sourceSize == destSize && sourceModTime == destModTime {
		return true
	}
	return false
}

func (dh *defaultHash) CompareHashValues(dstPath string, sourceSize int64, sourceHash string, chunkSize int64, hvs HashValues) (equal bool, hv *HashValue) {
	if sourceSize > 0 {
		// calculate the entire file hash value
		if len(hvs) == 0 || hvs.Last().Offset < sourceSize {
			hvs = append(hvs, NewHashValue(sourceSize, sourceHash))
		}
		hv, err := dh.CompareHashValuesWithFileName(dstPath, chunkSize, hvs)
		if err == nil && hv != nil {
			return hv.Offset == sourceSize && hv.Hash == sourceHash && len(sourceHash) > 0, hv
		}
	}
	return false, nil
}

func (dh *defaultHash) CompareHashValuesWithFileName(path string, chunkSize int64, hvs HashValues) (eq *HashValue, err error) {
	f, err := dh.open(path)
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
