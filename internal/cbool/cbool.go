package cbool

import "sync"

type CBool struct {
	v  bool
	mu sync.RWMutex
}

func New(v bool) *CBool {
	return &CBool{
		mu: sync.RWMutex{},
		v:  v,
	}
}

func (cb *CBool) Get() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.v
}

func (cb *CBool) Set(v bool) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.v = v
}
