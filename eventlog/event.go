package eventlog

import "fmt"

type event struct {
	name string
	op   string
}

// NewEvent create an event instance
func NewEvent(name, op string) event {
	return event{
		name: name,
		op:   op,
	}
}

// String return the format event info
func (e *event) String() string {
	return fmt.Sprintf("[%s][%s]", e.op, e.name)
}
