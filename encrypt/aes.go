package encrypt

import (
	"crypto/aes"
	"errors"
	"fmt"
)

var (
	aesIV = []byte("nosrc-gofs-aesiv")

	errInvalidAESIV = errors.New("IV length must equal block size")
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
		return fmt.Errorf("%w, iv length=%d, block size=%d", errInvalidAESIV, length, aes.BlockSize)
	}
	return nil
}
