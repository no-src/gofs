package eventlog

import (
	"fmt"
	"io"
)

// EventLog the event log recorder
type EventLog interface {
	// Write write event info to output writer
	Write(event Event) error
}

type eventLog struct {
	w io.Writer
}

// New create an EventLog instance with io.Writer
func New(w io.Writer) EventLog {
	return &eventLog{
		w: w,
	}
}

func (el *eventLog) Write(event Event) error {
	if el.w == nil {
		return nil
	}
	_, err := el.w.Write([]byte(fmt.Sprintf("%s\n", event.String())))
	return err
}
