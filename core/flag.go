package core

import "flag"

// FlagSet is a custom flag.FlagSet
type FlagSet struct {
	*flag.FlagSet
}

// NewFlagSet returns a new, empty flag set with the specified name and
// error handling property. If the name is not empty, it will be printed
// in the default usage message and in error messages.
func NewFlagSet(name string, errorHandling flag.ErrorHandling) *FlagSet {
	return &FlagSet{
		FlagSet: flag.NewFlagSet(name, errorHandling),
	}
}
