package progress

import (
	"io"

	"github.com/schollz/progressbar/v3"
)

// NewWriter wrap the io.Writer to support print write progress
func NewWriter(w io.Writer, size int64, desc string) io.Writer {
	if w == nil || size == 0 {
		return w
	}
	bar := progressbar.DefaultBytes(
		size,
		desc,
	)
	return io.MultiWriter(w, bar)
}

// NewWriterWithEnable wrap the io.Writer to support print write progress, if enable is false, then return the origin io.Writer
func NewWriterWithEnable(w io.Writer, size int64, desc string, enable bool) io.Writer {
	if !enable {
		return w
	}
	return NewWriter(w, size, desc)
}
