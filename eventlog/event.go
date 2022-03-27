package eventlog

import (
	"fmt"
	"github.com/no-src/gofs/util/timeutil"
)

// Event the description of file change event
type Event struct {
	// Name the path of file change
	Name string `json:"name"`
	// Op the operation of file change
	Op string `json:"op"`
	// Time the time of file change
	Time timeutil.Time `json:"time"`
}

// NewEvent create an event instance
func NewEvent(name, op string) Event {
	return Event{
		Name: name,
		Op:   op,
		Time: timeutil.Now(),
	}
}

// String return the format event info
func (e *Event) String() string {
	return fmt.Sprintf("[%s][%s][%s]", e.Time.String(), e.Op, e.Name)
}
