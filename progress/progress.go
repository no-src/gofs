package progress

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

// NewWriter wrap the io.Writer to support print write progress
func NewWriter(w io.Writer, size int64, desc string) io.Writer {
	if w == nil || size == 0 {
		return w
	}

	bar := progressbar.NewOptions64(
		size,
		progressbar.OptionSetDescription(desc),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionSetTheme(progressbar.Theme{Saucer: "=", SaucerHead: ">", SaucerPadding: "-", BarStart: "[", BarEnd: "]"}),
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
