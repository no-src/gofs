package encrypt

import (
	"bufio"
	"io"
)

type bufferWriter struct {
	bw *bufio.Writer
}

func (w *bufferWriter) Close() error {
	return w.bw.Flush()
}

func (w *bufferWriter) Write(p []byte) (n int, err error) {
	return w.bw.Write(p)
}

func newBufferWriter(w io.Writer) io.WriteCloser {
	return &bufferWriter{
		bw: bufio.NewWriter(w),
	}
}
