package cbool

import "sync"

type CBool struct {
	v  bool
	mu sync.Mutex
}

func New(v bool) *CBool {
	return &CBool{
		mu: sync.Mutex{},
		v:  v,
	}
}

func (cb *CBool) Get() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.v
}

func (cb *CBool) Set(v bool) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.v = v
}
