package cbool

import "testing"

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
