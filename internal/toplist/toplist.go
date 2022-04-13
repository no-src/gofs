package toplist

import (
	"errors"
	"sync"

	"github.com/no-src/gofs/util/jsonutil"
)

// TopList store some elements in list with specified capacity, the oldest elements that exceed specified capacity will be discarded
type TopList struct {
	capacity int
	length   int
	begin    int
	end      int
	data     []any
	mu       sync.RWMutex
	asc      bool
}

var errInvalidCapacity = errors.New("capacity must be greater than zero")

// New create a TopList with specified capacity and the capacity must be greater than zero, order by add time desc
func New(capacity int) (*TopList, error) {
	return newTopList(capacity, false)
}

// NewOrderByAsc create a TopList with specified capacity and the capacity must be greater than zero, order by add time asc
func NewOrderByAsc(capacity int) (*TopList, error) {
	return newTopList(capacity, true)
}

func newTopList(capacity int, asc bool) (*TopList, error) {
	if capacity <= 0 {
		return nil, errInvalidCapacity
	}
	tl := &TopList{
		capacity: capacity,
		begin:    0,
		end:      -1,
		data:     make([]any, capacity),
		asc:      asc,
	}
	return tl, nil
}

// Add add an element to the list
func (tl *TopList) Add(element any) {
	tl.mu.Lock()
	defer tl.mu.Unlock()
	nextEnd := tl.getNextEnd()
	tl.data[nextEnd] = element
	tl.end = nextEnd
	if tl.length < tl.capacity {
		tl.length++
	} else {
		tl.begin = tl.getNextBegin()
	}
}

// Len returns the length of the list
func (tl *TopList) Len() int {
	tl.mu.RLock()
	defer tl.mu.RUnlock()
	return tl.length
}

// Cap returns the capacity of the list
func (tl *TopList) Cap() int {
	tl.mu.RLock()
	defer tl.mu.RUnlock()
	return tl.capacity
}

// Get get an element by a specified index in the list, if the index is greater than TopList.Len or less than 0, return a nil element always
func (tl *TopList) Get(index int) any {
	tl.mu.RLock()
	defer tl.mu.RUnlock()
	return tl.get(index)
}

func (tl *TopList) get(index int) any {
	if index >= tl.length || index < 0 {
		return nil
	}
	if !tl.asc {
		index = tl.length - 1 - index
	}

	pos := tl.begin + index
	if pos >= tl.capacity {
		pos -= tl.capacity
	}
	return tl.data[pos]
}

// Last returns the last element of the list
func (tl *TopList) Last() any {
	tl.mu.RLock()
	defer tl.mu.RUnlock()
	if tl.length <= 0 {
		return nil
	}
	if tl.asc {
		return tl.data[tl.end]
	}
	return tl.data[tl.begin]
}

// Top returns the latest elements by top n
func (tl *TopList) Top(n int) (list []any) {
	tl.mu.RLock()
	defer tl.mu.RUnlock()
	if n <= 0 || tl.length == 0 {
		return list
	}

	if tl.length < n {
		n = tl.length
	}
	for i := 0; i < n; i++ {
		list = append(list, tl.get(i))
	}
	return list
}

// MarshalJSON implement interface json.Marshaler
func (tl *TopList) MarshalJSON() (text []byte, err error) {
	return jsonutil.Marshal(tl.Top(tl.Cap()))
}

func (tl *TopList) getNextPos(pos int) int {
	next := pos + 1
	if next >= tl.capacity {
		next = 0
	}
	return next
}

func (tl *TopList) getNextBegin() int {
	return tl.getNextPos(tl.begin)
}

func (tl *TopList) getNextEnd() int {
	return tl.getNextPos(tl.end)
}
