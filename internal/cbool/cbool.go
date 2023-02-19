package cbool

import (
	"sync/atomic"
)

// CBool a concurrent safe bool
type CBool struct {
	v atomic.Bool
}

// New create an instance of CBool
func New(v bool) *CBool {
	cb := &CBool{}
	cb.v.Store(v)
	return cb
}

// Get return the bool value
func (cb *CBool) Get() bool {
	return cb.v.Load()
}

// Set to set the bool value
func (cb *CBool) Set(v bool) {
	cb.v.Store(v)
}

// SetC to set the bool value and return a closed channel
func (cb *CBool) SetC(v bool) <-chan struct{} {
	cb.Set(v)
	c := make(chan struct{})
	close(c)
	return c
}
