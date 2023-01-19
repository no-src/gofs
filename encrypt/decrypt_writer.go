package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
)

type decryptWriter struct {
	w      io.Writer
	key    []byte
	iv     []byte
	stream cipher.Stream
}

func (w *decryptWriter) Write(p []byte) (nn int, err error) {
	dst := make([]byte, len(p))
	w.stream.XORKeyStream(dst, p)
	return w.w.Write(dst)
}

func newDecryptWriter(w io.Writer, key []byte, iv []byte) (io.Writer, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if err = checkAESIV(iv); err != nil {
		return nil, err
	}
	stream := cipher.NewCFBDecrypter(block, iv)
	return &decryptWriter{
		w:      w,
		key:    key,
		iv:     iv,
		stream: stream,
	}, nil
}
