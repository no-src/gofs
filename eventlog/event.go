package eventlog

import (
	"fmt"
	"github.com/no-src/gofs/util/timeutil"
	"time"
)

// Event the description of file change event
type Event struct {
	Name string
	Op   string
	Time timeutil.Time
}

// NewEvent create an event instance
func NewEvent(name, op string) Event {
	return Event{
		Name: name,
		Op:   op,
		Time: timeutil.NewTime(time.Now()),
	}
}

// String return the format event info
func (e *Event) String() string {
	return fmt.Sprintf("[%s][%s][%s]", e.Time.String(), e.Op, e.Name)
}
