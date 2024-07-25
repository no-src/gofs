package contract

import (
	"fmt"
	"testing"
)

func TestFsDirValue(t *testing.T) {
	testCases := []struct {
		isDir FsDirValue
		desc  string
		b     bool
	}{
		{FsIsDir, "1", true},
		{FsNotDir, "0", false},
		{FsUnknown, "-1", false},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			s := tc.isDir.String()
			if s != tc.desc {
				t.Errorf("FsDirValue.String() [%d] => expect: %v, but actual: %v \n", tc.isDir, tc.desc, s)
				return
			}

			is := tc.isDir.Is(tc.desc)
			if !is {
				t.Errorf("FsDirValue[%d].Is(%s)  => expect: true, but actual: false \n", tc.isDir, tc.desc)
				return
			}

			not := tc.isDir.Not(tc.desc)
			if not {
				t.Errorf("FsDirValue[%d].Not(%s)  => expect: false, but actual: true \n", tc.isDir, tc.desc)
				return
			}

			b := tc.isDir.Bool()
			if not {
				t.Errorf("FsDirValue[%d].Bool()  => expect: %v, but actual: %v \n", tc.isDir, tc.b, b)
			}
		})
	}
}

func TestParseFsDirValue(t *testing.T) {
	testCases := []struct {
		b      bool
		expect FsDirValue
	}{
		{true, FsIsDir},
		{false, FsNotDir},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v", tc.b), func(t *testing.T) {
			actual := ParseFsDirValue(tc.b)
			if actual != tc.expect {
				t.Errorf("ParseFsDirValue(%v)  => expect: %v, but actual: %v \n", tc.b, tc.expect, actual)
			}
		})
	}
}
