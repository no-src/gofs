package cbool

import (
	"sync"
	"testing"
)

func TestCBool(t *testing.T) {
	expect := true
	cb := New(expect)
	actual := cb.Get()
	if actual != expect {
		t.Errorf("test CBoll New and Get failed, expect:%v, actual:%v", expect, actual)
	}

	expect = false
	cb.Set(expect)
	actual = cb.Get()
	if actual != expect {
		t.Errorf("test CBoll Set and Get failed, expect:%v, actual:%v", expect, actual)
	}

	expect = true
	c := cb.SetC(expect)
	actual = cb.Get()
	if actual != expect {
		t.Errorf("test CBoll SetC and Get failed, expect:%v, actual:%v", expect, actual)
	}
	_, ok := <-c
	if ok {
		t.Errorf("test CBoll SetC value failed, channel should be closed")
	}
}

func TestCBool_Concurrent(t *testing.T) {
	cb := New(false)
	wg := sync.WaitGroup{}
	count := 10
	wg.Add(count * 3)
	for i := 0; i < count; i++ {
		go func() {
			cb.Get()
			wg.Done()
		}()

		go func() {
			cb.Set(true)
			wg.Done()
		}()

		go func() {
			<-cb.SetC(true)
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkCBool(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	cb := New(false)
	for i := 0; i < b.N; i++ {
		cb.Set(true)
		cb.Get()
	}
}
