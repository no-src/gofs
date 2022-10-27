package progress

import (
	"bytes"
	"io"
	"testing"
)

func TestNewWriterWithEnable(t *testing.T) {
	testCases := []struct {
		desc   string
		w      io.Writer
		size   int64
		enable bool
	}{
		{"disable progress", io.Discard, 100, false},
		{"normal progress", io.Discard, 100, true},
		{"progress with nil writer", nil, 100, true},
		{"progress with zero size", io.Discard, 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			w := NewWriterWithEnable(tc.w, tc.size, tc.desc, tc.enable)
			if tc.w == nil && w != nil {
				t.Errorf("expect to get a nil writer but not")
				return
			} else if tc.w == nil && w == nil {
				return
			}
			data := bytes.Repeat([]byte{90}, int(tc.size))
			_, err := w.Write(data)
			if err != nil {
				t.Errorf("write progress error => %v", err)
			}
		})
	}
}
