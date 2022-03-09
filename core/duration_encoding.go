package core

import "time"

// MarshalText implement interface encoding.TextMarshaler
func (d Duration) MarshalText() (text []byte, err error) {
	return []byte(d.Duration().String()), nil
}

// UnmarshalText implement interface encoding.TextUnmarshaler
func (d *Duration) UnmarshalText(data []byte) error {
	od, err := time.ParseDuration(string(data))
	if err != nil {
		return err
	}
	*d = Duration(od)
	return nil
}
