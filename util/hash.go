package util

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"os"
)

func MD5FromFile(file *os.File, bufSize int) (hash string, err error) {
	if file == nil {
		err = errors.New("file is nil")
		return hash, err
	}
	block := make([]byte, bufSize)
	md5Provider := md5.New()
	reader := bufio.NewReader(file)
	for {
		_, err = reader.Read(block)
		if err == io.EOF {
			break
		}
		if err != nil {
			return hash, err
		}
		_, err = md5Provider.Write(block)
		if err != nil {
			return hash, err
		}
	}
	sum := md5Provider.Sum(nil)
	hash = hex.EncodeToString(sum)
	return hash, nil
}
