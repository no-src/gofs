package encrypt

import (
	"archive/zip"
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"io"
)

type encryptWriter struct {
	bw     *bufio.Writer
	zw     *zip.Writer
	key    []byte
	iv     []byte
	stream cipher.Stream
}

func (w *encryptWriter) Write(p []byte) (nn int, err error) {
	dst := make([]byte, len(p))
	w.stream.XORKeyStream(dst, p)
	return w.bw.Write(dst)
}

func (w *encryptWriter) Close() error {
	if err := w.bw.Flush(); err != nil {
		return err
	}
	return w.zw.Close()
}

func newEncryptWriter(w io.Writer, name string, key []byte, iv []byte) (io.WriteCloser, error) {
	zw := zip.NewWriter(w)
	ew, err := zw.CreateHeader(&zip.FileHeader{
		Name:   name,
		Method: zip.Store,
	})
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if err = checkAESIV(iv); err != nil {
		return nil, err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	bw := bufio.NewWriter(ew)
	return &encryptWriter{
		bw:     bw,
		zw:     zw,
		key:    key,
		iv:     iv,
		stream: stream,
	}, err
}
