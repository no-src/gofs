package stringutil

import (
	"encoding/json"
	"errors"
	"github.com/no-src/gofs/util/jsonutil"
	"testing"
)

func TestString(t *testing.T) {
	// custom
	testString(t, str{
		name: "golang",
	}, "golang")

	// custom with error
	jsonutil.Marshal = func(v interface{}) ([]byte, error) {
		return nil, errors.New("marshal error test")
	}
	testString(t, []string{"hello"}, `[hello]`)
	jsonutil.Marshal = json.Marshal

	// string
	testString(t, "hello", "hello")

	// int
	testString(t, int(100), "100")

	// uint64
	testString(t, uint64(200), "200")

	// int64
	testString(t, int64(300), "300")

	// bool
	testString(t, true, "true")

	// error
	testString(t, errors.New("test error info"), "test error info")

	// null
	testString(t, nil, "null")

	// other
	testString(t, []string{"hello", "world"}, `["hello","world"]`)
}

func testString(t *testing.T, v interface{}, expect string) {
	actual := String(v)
	if actual != expect {
		t.Errorf("test String error, expect:%s, actual:%s", expect, actual)
	}
}

func TestInt64(t *testing.T) {
	s := "100001"
	var expect int64 = 100001
	actual, err := Int64(s)
	if err != nil {
		t.Errorf("test Int64 error [%s] => %s", s, err)
		return
	}
	if actual != expect {
		t.Errorf("test Int64 error, expect:%d,actual:%d", expect, actual)
	}
}

func TestInt64Error(t *testing.T) {
	s := "100001x"
	_, err := Int64(s)
	if err == nil {
		t.Errorf("test Int64 error, should get an error => %s", s)
		return
	}
}

func TestIsEmpty(t *testing.T) {
	expect := true
	testIsEmpty(t, "", expect)
	testIsEmpty(t, " ", expect)
	testIsEmpty(t, "		", expect)
	testIsEmpty(t, "\t", expect)
	testIsEmpty(t, "\r", expect)
	testIsEmpty(t, "\n", expect)
	expect = false
	testIsEmpty(t, "hello", expect)
}

func testIsEmpty(t *testing.T, s string, expect bool) {
	if IsEmpty(s) != expect {
		t.Errorf("test IsEmpty error, expect %v", expect)
	}
}

type str struct {
	name string
}

func (s str) String() string {
	return s.name
}
