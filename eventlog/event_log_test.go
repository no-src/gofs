package eventlog

import (
	"io"
	"os"
	"testing"
)

func TestEventLogStdoutWriter(t *testing.T) {
	testEventLogNilWriter(t, os.Stdout)
}

func TestEventLogNilWriter(t *testing.T) {
	testEventLogNilWriter(t, nil)
}

func testEventLogNilWriter(t *testing.T, w io.Writer) {
	el := New(w)
	e := NewEvent("gofs.txt", "Write")
	err := el.Write(e)
	if err != nil {
		t.Errorf("write event error => %s", err)
	}
}
