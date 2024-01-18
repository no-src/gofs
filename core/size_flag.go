package core

import (
	"strings"

	"github.com/no-src/nsgo/unit"
)

// SizeVar defines a Size flag with specified name, default value, and usage string.
// The argument p points to a Size variable in which to store the value of the flag.
func (f *FlagSet) SizeVar(p *Size, name string, value string, usage string) {
	size, err := newSize(value)
	if err != nil {
		panic(err)
	}
	*p = *size
	f.Var(p, name, usage)
}

func newSize(s string) (*Size, error) {
	size := new(Size)
	if err := size.Set(s); err != nil {
		return nil, err
	}
	return size, nil
}

// Set implement the Set function for the flag.Value interface
func (s *Size) Set(str string) error {
	v, _, err := unit.ParseBytes(str)
	if err != nil {
		return err
	}
	s.bytes = v
	s.origin = strings.ReplaceAll(str, " ", "")
	return nil
}

// String implement the String function for the flag.Value interface
func (s *Size) String() string {
	return s.origin
}
