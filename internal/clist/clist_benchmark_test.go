package clist

import (
	"testing"
)

func BenchmarkCList_PushBack(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	cl := New()
	for i := 0; i < b.N; i++ {
		cl.PushBack(i)
	}
}

func BenchmarkCList_Front(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	cl := newTestCList()
	for i := 0; i < b.N; i++ {
		cl.Front()
	}
}

func BenchmarkCList_Len(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	cl := newTestCList()
	for i := 0; i < b.N; i++ {
		cl.Len()
	}
}

func BenchmarkCList(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	cl := New()
	for i := 0; i < b.N; i++ {
		cl.PushBack(i)
		cl.Len()
		cl.Remove(cl.Front())
	}
}

func newTestCList() *CList {
	cl := New()
	cl.PushBack("hello")
	cl.PushBack("world")
	for i := 0; i < 100; i++ {
		cl.PushBack(i)
	}
	return cl
}
