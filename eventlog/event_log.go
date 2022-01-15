package eventlog

import (
	"fmt"
	"io"
	"time"
)

type EventLog interface {
	Write(event event)
}

type eventLog struct {
	w io.Writer
}

func New(w io.Writer) EventLog {
	return &eventLog{
		w: w,
	}
}

func (el *eventLog) Write(event event) {
	el.w.Write([]byte(fmt.Sprintf("[%s]%s\n", time.Now().Format("2006-01-02 15:04:05.999"), event.String())))
}
