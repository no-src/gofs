package stringutil

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/no-src/gofs/util/jsonutil"
)

func TestString(t *testing.T) {
	alwaysErrorMarshal := func(v any) ([]byte, error) {
		return nil, errors.New("marshal error test")
	}
	testCases := []struct {
		name    string
		marshal func(v any) ([]byte, error)
		v       any
		expect  string
	}{
		{"custom", json.Marshal, str{name: "golang"}, "golang"},
		{"custom with error", alwaysErrorMarshal, []string{"hello"}, `[hello]`},
		{"string", json.Marshal, "hello", "hello"},
		{"int", json.Marshal, 100, "100"},
		{"uint64", json.Marshal, uint64(200), "200"},
		{"int64", json.Marshal, int64(300), "300"},
		{"bool", json.Marshal, true, "true"},
		{"error", json.Marshal, errors.New("test error info"), "test error info"},
		{"null", json.Marshal, nil, "null"},
		{"other", json.Marshal, []string{"hello", "world"}, `["hello","world"]`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonutil.Marshal = tc.marshal
			testString(t, tc.v, tc.expect)
		})
	}

	// reset
	jsonutil.Marshal = json.Marshal
}

func testString(t *testing.T, v any, expect string) {
	actual := String(v)
	if actual != expect {
		t.Errorf("test String error, expect:%s, actual:%s", expect, actual)
	}
}

func TestInt64(t *testing.T) {
	testCases := []struct {
		str    string
		expect int64
	}{
		{"-1", -1},
		{"0", 0},
		{"1", 1},
		{"0005", 5},
		{"+10", 10},
		{"100001", 100001},
	}

	for _, tc := range testCases {
		t.Run(tc.str, func(t *testing.T) {
			actual, err := Int64(tc.str)
			if err != nil {
				t.Errorf("test Int64 error [%s] => %s", tc.str, err)
				return
			}
			if actual != tc.expect {
				t.Errorf("test Int64 error, expect:%d,actual:%d", tc.expect, actual)
			}
		})
	}
}

func TestInt64_ReturnError(t *testing.T) {
	testCases := []struct {
		str string
	}{
		{""},
		{"		"},
		{"\t"},
		{"\r"},
		{"\n"},
		{"a"},
		{"abc"},
		{"100001x"},
		{"@#()"},
		{"123@#()"},
		{"--1"},
		{" 11"},
	}

	for _, tc := range testCases {
		t.Run(tc.str, func(t *testing.T) {
			_, err := Int64(tc.str)
			if err == nil {
				t.Errorf("test Int64 error, expect to get an error but get nil => %s", tc.str)
			}
		})
	}
}

func TestIsEmpty(t *testing.T) {
	testCases := []struct {
		str    string
		expect bool
	}{
		{"", true},
		{" ", true},
		{"		", true},
		{"\t", true},
		{"\r", true},
		{"\n", true},
		{"hello", false},
	}

	for _, tc := range testCases {
		t.Run("["+tc.str+"]", func(t *testing.T) {
			if IsEmpty(tc.str) != tc.expect {
				t.Errorf("test IsEmpty error, expect %v", tc.expect)
			}
		})
	}
}

type str struct {
	name string
}

func (s str) String() string {
	return s.name
}
