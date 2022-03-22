package hashutil

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"os"
)

// MD5FromFile calculate the hash value of the file
// If you reuse the file reader, please set its offset to start position first, like os.File.Seek
func MD5FromFile(file io.Reader) (hash string, err error) {
	if file == nil {
		err = errors.New("file is nil")
		return hash, err
	}
	md5Provider := md5.New()
	reader := bufio.NewReader(file)
	_, err = reader.WriteTo(md5Provider)
	if err != nil {
		return hash, err
	}
	sum := md5Provider.Sum(nil)
	hash = hex.EncodeToString(sum)
	return hash, nil
}

// MD5FromFileName calculate the hash value of the file
func MD5FromFileName(path string) (hash string, err error) {
	if len(path) == 0 {
		err = errors.New("file path can't be empty")
		return hash, err
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
		err = errors.New("file path can't be empty")
		return hash, err
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
	md5Provider := md5.New()
	md5Provider.Write(bytes)
	sum := md5Provider.Sum(nil)
	hash = hex.EncodeToString(sum)
	return hash
}

// MD5FromString calculate the hash value of the string
func MD5FromString(s string) (hash string) {
	return MD5([]byte(s))
}
