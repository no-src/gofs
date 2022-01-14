package clist

import (
	"container/list"
	"sync"
)

type CList struct {
	v  *list.List
	mu sync.Mutex
}

func New() *CList {
	return &CList{
		v:  list.New(),
		mu: sync.Mutex{},
	}
}

func (cl *CList) PushBack(v interface{}) *list.Element {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	return cl.v.PushBack(v)
}

func (cl *CList) Front() *list.Element {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	return cl.v.Front()
}

func (cl *CList) Remove(e *list.Element) interface{} {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	return cl.v.Remove(e)
}

func (cl *CList) Len() int {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	return cl.v.Len()
}
