package encrypt

import (
	"io"
)

type decryptWriter struct {
	w      io.Writer
	secret []byte
	index  int
}

func (w *decryptWriter) Write(p []byte) (nn int, err error) {
	pLen := len(p)
	secretLen := len(w.secret)
	for i := 0; i < pLen; i++ {
		p[i] = p[i] ^ w.secret[w.index%secretLen]
		w.index++
	}
	return w.w.Write(p)
}

func newDecryptWriter(w io.Writer, secret []byte) io.Writer {
	return &decryptWriter{
		w:      w,
		secret: secret,
	}
}
