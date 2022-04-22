package eventlog

import (
	"io"
	"os"
	"testing"
)

func TestEventLog_Write(t *testing.T) {
	testCases := []struct {
		name string
		w    io.Writer
	}{
		{"stdout", os.Stdout},
		{"nil writer", nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			el := New(tc.w)
			e := NewEvent("gofs.txt", "Write")
			err := el.Write(e)
			if err != nil {
				t.Errorf("write event error => %s", err)
			}
		})
	}
}
