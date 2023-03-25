package core

import (
	"time"
)

// DurationVar defines a core.Duration flag with specified name, default value, and usage string.
// The argument p points to a core.Duration variable in which to store the value of the flag.
// The flag accepts a value acceptable to time.ParseDuration.
func (f *FlagSet) DurationVar(p *Duration, name string, value time.Duration, usage string) {
	dp := (*time.Duration)(p)
	f.FlagSet.DurationVar(dp, name, value, usage)
}
