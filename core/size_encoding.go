package core

// MarshalText implement interface encoding.TextMarshaler
func (s Size) MarshalText() (text []byte, err error) {
	return []byte(s.String()), nil
}

// UnmarshalText implement interface encoding.TextUnmarshaler
func (s *Size) UnmarshalText(data []byte) error {
	size, err := newSize(string(data))
	if err != nil {
		return err
	}
	*s = *size
	return nil
}
