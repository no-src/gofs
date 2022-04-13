package toplist

import (
	"bytes"
	"testing"

	"github.com/no-src/gofs/util/jsonutil"
)

func TestTopListWithCapZero(t *testing.T) {
	// desc
	capacity := 0
	_, err := New(capacity)
	if err != errInvalidCapacity {
		t.Errorf("[desc] test toplist with zero capacity error, expect get an error")
	}

	// asc
	capacity = 0
	_, err = NewOrderByAsc(capacity)
	if err != errInvalidCapacity {
		t.Errorf("[asc] test toplist with zero capacity error, expect get an error")
	}
}

func TestTopListWithCapOne(t *testing.T) {
	// desc
	capacity := 1
	tl, err := New(capacity)
	if err != nil {
		t.Errorf("[desc] test toplist with one capacity error, get an error => %v", err)
		return
	}

	topListLast(t, tl, nil)

	s := "hello"
	tl.Add(s)
	topListGet(t, tl, 0, s)

	s = "world"
	tl.Add(s)
	topListGet(t, tl, 0, s)
	topListGet(t, tl, -1, nil)
	topListGet(t, tl, capacity, nil)
	topListLast(t, tl, s)
	topListLen(t, tl, capacity)
	topListCap(t, tl, capacity)
	topListTop(t, tl, 0)
	topListTop(t, tl, 1, s)
	topListTop(t, tl, 2, s)

	// asc
	capacity = 1
	tl, err = NewOrderByAsc(capacity)
	if err != nil {
		t.Errorf("[asc] test toplist with one capacity error, get an error => %v", err)
		return
	}

	topListLast(t, tl, nil)

	s = "hello"
	tl.Add(s)
	topListGet(t, tl, 0, s)

	s = "world"
	tl.Add(s)
	topListGet(t, tl, 0, s)
	topListGet(t, tl, -1, nil)
	topListGet(t, tl, capacity, nil)
	topListLast(t, tl, s)
	topListLen(t, tl, capacity)
	topListCap(t, tl, capacity)
	topListTop(t, tl, 0)
	topListTop(t, tl, 1, s)
	topListTop(t, tl, 2, s)
}

func TestTopListWithCapTwo(t *testing.T) {
	// desc
	capacity := 2
	tl, err := New(capacity)
	if err != nil {
		t.Errorf("[desc] test toplist with one capacity error, get an error => %v", err)
		return
	}

	topListLast(t, tl, nil)

	s1 := "hello"
	tl.Add(s1)
	topListGet(t, tl, 0, s1)
	topListLen(t, tl, 1)

	s2 := "world"
	tl.Add(s2)
	topListGet(t, tl, 0, s2)
	topListGet(t, tl, 1, s1)
	topListGet(t, tl, -1, nil)
	topListGet(t, tl, capacity, nil)
	topListLast(t, tl, s1)
	topListLen(t, tl, capacity)
	topListCap(t, tl, capacity)

	s3 := "golang"
	tl.Add(s3)
	topListGet(t, tl, 0, s3)

	topListGet(t, tl, 1, s2)
	topListGet(t, tl, -1, nil)
	topListGet(t, tl, capacity, nil)
	topListLast(t, tl, s2)
	topListLen(t, tl, capacity)
	topListCap(t, tl, capacity)
	topListTop(t, tl, 0)
	topListTop(t, tl, 1, s3)
	topListTop(t, tl, 2, s3, s2)

	// asc
	capacity = 2
	tl, err = NewOrderByAsc(capacity)
	if err != nil {
		t.Errorf("[asc] test toplist with one capacity error, get an error => %v", err)
		return
	}

	topListLast(t, tl, nil)

	s1 = "hello"
	tl.Add(s1)
	topListGet(t, tl, 0, s1)
	topListLen(t, tl, 1)

	s2 = "world"
	tl.Add(s2)
	topListGet(t, tl, 0, s1)
	topListGet(t, tl, 1, s2)
	topListGet(t, tl, -1, nil)
	topListGet(t, tl, capacity, nil)
	topListLast(t, tl, s2)
	topListLen(t, tl, capacity)
	topListCap(t, tl, capacity)

	s3 = "golang"
	tl.Add(s3)
	topListGet(t, tl, 0, s2)

	topListGet(t, tl, 1, s3)
	topListGet(t, tl, -1, nil)
	topListGet(t, tl, capacity, nil)
	topListLast(t, tl, s3)
	topListLen(t, tl, capacity)
	topListCap(t, tl, capacity)
	topListTop(t, tl, 0)
	topListTop(t, tl, 1, s2)
	topListTop(t, tl, 2, s2, s3)
}

func TestTopListMarshalJSON(t *testing.T) {
	capacity := 2
	tl, err := New(capacity)
	if err != nil {
		t.Errorf("create new toplist error, get an error => %v", err)
		return
	}
	tl.Add("hello")
	tl.Add("world")
	jsonBytes, err := jsonutil.Marshal(tl)
	if err != nil {
		t.Errorf("test marshal toplist error => %v", err)
		return
	}
	expectJsonBytes := []byte(`["world","hello"]`)
	if !bytes.Equal(expectJsonBytes, jsonBytes) {
		t.Errorf("test marshal toplist error, expect:%s, actual:%s", string(expectJsonBytes), string(jsonBytes))
	}
}

func topListGet(t *testing.T, tl *TopList, index int, expect any) {
	order := "desc"
	if tl.asc {
		order = "asc"
	}
	actual := tl.Get(index)
	if actual != expect {
		t.Errorf("[%s] get toplist element that index=0 error, expect:%v, actual:%v", order, expect, actual)
	}
}

func topListLast(t *testing.T, tl *TopList, expect any) {
	order := "desc"
	if tl.asc {
		order = "asc"
	}
	actual := tl.Last()
	if actual != expect {
		t.Errorf("[%s] get toplist last element error, expect:%v, actual:%s", order, expect, actual)
	}
}

func topListLen(t *testing.T, tl *TopList, expect int) {
	order := "desc"
	if tl.asc {
		order = "asc"
	}
	actual := tl.Len()
	if actual != expect {
		t.Errorf("[%s] get toplist length error, expect:%d, actual:%d", order, expect, actual)
	}
}

func topListCap(t *testing.T, tl *TopList, expect int) {
	order := "desc"
	if tl.asc {
		order = "asc"
	}
	actual := tl.Cap()
	if actual != expect {
		t.Errorf("[%s] get toplist capacity error, expect:%d, actual:%d", order, expect, actual)
	}
}

func topListTop(t *testing.T, tl *TopList, top int, expect ...any) {
	order := "desc"
	if tl.asc {
		order = "asc"
	}
	actual := tl.Top(top)
	if len(actual) != len(expect) {
		t.Errorf("[%s] get toplist top error, expect length:%d, actual length:%d", order, len(expect), len(actual))
		return
	}

	for index, expectElement := range expect {
		actualElement := actual[index]
		if actualElement != expectElement {
			t.Errorf("[%s] get toplist top error, expect element:%v, actual element:%v", order, expectElement, actualElement)
			return
		}
	}
}
