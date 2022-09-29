package encrypt

import (
	"archive/zip"
	"bufio"
	"io"
)

type encryptWriter struct {
	bw     *bufio.Writer
	zw     *zip.Writer
	secret []byte
	index  int
}

func (w *encryptWriter) Write(p []byte) (nn int, err error) {
	pLen := len(p)
	secretLen := len(w.secret)
	for i := 0; i < pLen; i++ {
		p[i] = p[i] ^ w.secret[w.index%secretLen]
		w.index++
	}
	return w.bw.Write(p)
}

func (w *encryptWriter) Close() error {
	if err := w.bw.Flush(); err != nil {
		return err
	}
	return w.zw.Close()
}

func newEncryptWriter(w io.Writer, name string, secret []byte) (io.WriteCloser, error) {
	zw := zip.NewWriter(w)
	ew, err := zw.Create(name)
	if err != nil {
		return nil, err
	}
	bw := bufio.NewWriter(ew)
	return &encryptWriter{
		bw:     bw,
		zw:     zw,
		secret: secret,
	}, err
}
