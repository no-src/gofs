//go:build !go1.19

package cbool

import "sync"

// CBool a concurrent safe bool
type CBool struct {
	v  bool
	mu sync.RWMutex
}

// New create an instance of CBool
func New(v bool) *CBool {
	return &CBool{
		v: v,
	}
}

// Get return the bool value
func (cb *CBool) Get() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.v
}

// Set to set the bool value
func (cb *CBool) Set(v bool) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.v = v
}

// SetC to set the bool value and return a closed channel
func (cb *CBool) SetC(v bool) <-chan struct{} {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.v = v
	c := make(chan struct{})
	close(c)
	return c
}
