package encrypt

import (
	"crypto/aes"
	"fmt"
)

var (
	aesIV = []byte("nosrc-gofs-aesiv")
)

func checkAESKey(key []byte) error {
	length := len(key)
	if length == 16 || length == 24 || length == 32 {
		return nil
	}
	return aes.KeySizeError(length)
}

func checkAESIV(iv []byte) error {
	length := len(iv)
	if length != aes.BlockSize {
		return fmt.Errorf("IV length must equal block size, iv length=%d, block size=%d", length, aes.BlockSize)
	}
	return nil
}
