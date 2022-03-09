package core

import (
	"time"
)

// Duration a duration with custom encoding
type Duration time.Duration

// Duration return the origin time.Duration
func (d Duration) Duration() time.Duration {
	return time.Duration(d)
}
