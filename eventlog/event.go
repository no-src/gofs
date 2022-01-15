package eventlog

import "fmt"

type event struct {
	name string
	op   string
}

func NewEvent(name, op string) event {
	return event{
		name: name,
		op:   op,
	}
}

func (e *event) String() string {
	return fmt.Sprintf("[%s][%s]", e.op, e.name)
}
