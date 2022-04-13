package clist

import (
	"container/list"
	"sync"
)

// CList a concurrent safe list
type CList struct {
	v  *list.List
	mu sync.Mutex
}

// New create an instance of CList
func New() *CList {
	return &CList{
		v: list.New(),
	}
}

// PushBack inserts a new element e with value v at the back of list l and returns e.
func (cl *CList) PushBack(v any) *list.Element {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	return cl.v.PushBack(v)
}

// Front returns the first element of list l or nil if the list is empty.
func (cl *CList) Front() *list.Element {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	return cl.v.Front()
}

// Remove removes e from l if e is an element of list l.
// It returns the element value e.Value.
// The element must not be nil.
func (cl *CList) Remove(e *list.Element) any {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	return cl.v.Remove(e)
}

// Len returns the number of elements of list l.
// The complexity is O(1).
func (cl *CList) Len() int {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	return cl.v.Len()
}
