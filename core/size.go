package core

import "strconv"

// Size represent the size of data
type Size struct {
	bytes  uint64
	origin string
}

// NewSize create an instance of Size from bytes size
func NewSize(size uint64) Size {
	s, _ := newSize(strconv.FormatUint(size, 10))
	return *s
}

// Bytes return the bytes of the size
func (s Size) Bytes() int64 {
	return int64(s.bytes)
}
