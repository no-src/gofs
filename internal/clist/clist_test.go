package clist

import (
	"sync"
	"testing"
)

func TestCList(t *testing.T) {
	cl := New()
	s1 := "hello"
	s2 := "world"
	s3 := "and"
	s4 := "gopher"
	cl.PushBack(s1)
	cl.PushBack(s2)
	cl.PushBack(s3)
	cl.PushBack(s4)
	expectLen := 4
	actualLen := cl.Len()
	if expectLen != actualLen {
		t.Errorf("test CList Len failed, expect:%d, actual:%d", expectLen, actualLen)
	}

	el := cl.Front()
	if el == nil || el.Value.(string) != s1 {
		t.Errorf("test CList Front failed, expect:%s, actual:%s", s1, el.Value.(string))
		return
	}

	expectLen = 4
	actualLen = cl.Len()
	if expectLen != actualLen {
		t.Errorf("test CList Len failed, expect:%d, actual:%d", expectLen, actualLen)
	}

	rs := cl.Remove(el).(string)
	if rs != s1 {
		t.Errorf("test CList Remove failed, expect:%s, actual:%s", s1, rs)
	}

	expectLen = 3
	actualLen = cl.Len()
	if expectLen != actualLen {
		t.Errorf("test CList Len failed, expect:%d, actual:%d", expectLen, actualLen)
	}
}

func TestCList_Concurrent(t *testing.T) {
	cl := New()
	wg := sync.WaitGroup{}
	count := 10
	wg.Add(count * 4)
	for i := 0; i < count; i++ {
		go func(data int) {
			cl.PushBack(data)
			wg.Done()
		}(i)

		go func() {
			elem := cl.Front()
			wg.Done()
			if elem != nil {
				cl.Remove(elem)
			}
			wg.Done()
		}()

		go func() {
			cl.Len()
			wg.Done()
		}()
	}
	wg.Wait()
}
