package timeutil

import (
	"time"
)

// Time a custom time to change the behavior of marshal
type Time int64

const defaultTimeFormat = "2006-01-02 15:04:05"

// NewTime create a Time instance with the time.Time
func NewTime(t time.Time) Time {
	return Time(t.Unix())
}

// Unix returns the unix time, the number of seconds elapsed since January 1, 1970 UTC
func (t Time) Unix() int64 {
	return int64(t)
}

// String returns the default formatted time string
func (t Time) String() string {
	return time.Unix(int64(t), 0).Format(defaultTimeFormat)
}

// MarshalText implement interface encoding.TextMarshaler
func (t Time) MarshalText() (text []byte, err error) {
	return []byte(t.String()), nil
}
