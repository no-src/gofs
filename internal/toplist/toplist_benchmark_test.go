package toplist

import (
	"fmt"
	"testing"
)

func BenchmarkTopList_Add(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	tl := newTestTopList()
	for i := 0; i < b.N; i++ {
		tl.Add(i)
	}
}

func BenchmarkTopList_Get(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	tl := newTestTopList()
	for i := 0; i < b.N; i++ {
		tl.Get(i)
	}
}

func BenchmarkTopList_Top(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	tl := newTestTopList()
	for i := 0; i < b.N; i++ {
		tl.Top(5)
	}
}

func BenchmarkTopList_Len(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	tl := newTestTopList()
	for i := 0; i < b.N; i++ {
		tl.Len()
	}
}

func BenchmarkTopList_Cap(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	tl := newTestTopList()
	for i := 0; i < b.N; i++ {
		tl.Cap()
	}
}

func BenchmarkTopList_Last(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	tl := newTestTopList()
	for i := 0; i < b.N; i++ {
		tl.Last()
	}
}

func BenchmarkTopList_MarshalJSON(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	tl := newTestTopList()
	for i := 0; i < b.N; i++ {
		tl.MarshalJSON()
	}
}

func BenchmarkTopList(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	tl := newTestTopList()
	for i := 0; i < b.N; i++ {
		tl.Add(i)
		tl.Get(10)
		tl.Len()
		tl.Cap()
		tl.Top(5)
		tl.Last()
	}
}

func newTestTopList() *TopList {
	tl, err := New(10)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 100; i++ {
		tl.Add(fmt.Sprintf("default-%d", i))
	}
	return tl
}
