package core

import (
	"fmt"
	"strings"
	"testing"

	"github.com/no-src/nsgo/jsonutil"
)

func TestNewSize(t *testing.T) {
	testCases := []struct {
		bytes uint64
	}{
		{0},
		{1},
		{1000 * 1000},
		{1024 * 1024},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%d", tc.bytes), func(t *testing.T) {
			size := NewSize(tc.bytes)
			if size.Bytes() != int64(tc.bytes) {
				t.Errorf("test NewSize error, expect:%d, actual:%d", tc.bytes, size.Bytes())
			}
		})
	}
}

func TestSize_MarshalText(t *testing.T) {
	testCases := []struct {
		size   Size
		expect string
	}{
		{NewSize(0), `"0"`},
		{NewSize(1), `"1"`},
		{NewSize(1024), `"1024"`},
	}

	for _, tc := range testCases {
		t.Run(tc.expect, func(t *testing.T) {
			actual, err := jsonutil.Marshal(tc.size)
			if err != nil {
				t.Errorf("test Size MarshalText error =>%s", err)
				return
			}
			if string(actual) != tc.expect {
				t.Errorf("test Size MarshalText error, expect:%s, actual:%s", tc.expect, string(actual))
			}
		})
	}
}

func TestSize_UnmarshalText(t *testing.T) {
	testCases := []struct {
		s      string
		expect int64
	}{
		{`"1MB"`, 1000 * 1000},
		{`"1MiB"`, 1024 * 1024},
		{`"1B"`, 1},
		{`"1"`, 1},
		{`"0"`, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.s, func(t *testing.T) {
			data := []byte(tc.s)
			var actual Size
			err := jsonutil.Unmarshal(data, &actual)
			if err != nil {
				t.Errorf("test Size unmarshal error =>%s", err)
				return
			}
			if actual.Bytes() != tc.expect {
				t.Errorf("test Size unmarshal error, expect:%d, actual:%d", tc.expect, actual.Bytes())
			}
		})
	}
}

func TestSize_UnmarshalText_ReturnError(t *testing.T) {
	testCases := []struct {
		s string
	}{
		{`""`},
		{""},
		{`"1x"`},
		{"2x"},
	}

	for _, tc := range testCases {
		t.Run(tc.s, func(t *testing.T) {
			var actual Size
			data := []byte(tc.s)
			if err := jsonutil.Unmarshal(data, &actual); err == nil {
				t.Errorf("test Size unmarshal error, expect to get an error but get nil")
			}
		})
	}
}

func TestSizeVar(t *testing.T) {
	testCases := []struct {
		s      string
		expect int64
	}{
		{"1MB", 1000 * 1000},
		{"1MiB", 1024 * 1024},
		{"1B", 1},
		{"1 MB", 1000 * 1000},
		{"1 MiB", 1024 * 1024},
		{"1 B", 1},
		{"1", 1},
		{"0", 0},
	}

	for _, tc := range testCases {
		t.Run(tc.s, func(t *testing.T) {
			var actual Size
			defaultValue := "1MB"
			flagName := "core_test_size" + tc.s
			testCommandLine.SizeVar(&actual, flagName, defaultValue, "test size")
			parseFlag(fmt.Sprintf("-%s=%s", flagName, tc.s))
			if actual.Bytes() != tc.expect {
				t.Errorf("test SizeVar error, Bytes() expect:%d, actual:%d", tc.expect, actual.Bytes())
			}
			expectStr := strings.ReplaceAll(tc.s, " ", "")
			if actual.String() != expectStr {
				t.Errorf("test SizeVar error, String() expect:%s, actual:%s", expectStr, actual.String())
			}
		})
	}
}

func TestSizeVar_ReturnError(t *testing.T) {
	testCases := []struct {
		s string
	}{
		{"1x"},
		{""},
	}
	for _, tc := range testCases {
		t.Run(tc.s, func(t *testing.T) {
			var actual Size
			flagName := "core_test_size_error" + tc.s
			getPanic := isPanic(func() {
				testCommandLine.SizeVar(&actual, flagName, tc.s, "test size")
			})
			if !getPanic {
				t.Errorf("test SizeVar error, expect to get panic but actual not")
			}
		})
	}
}

func isPanic(f func()) (b bool) {
	defer func() {
		if x := recover(); x != nil {
			b = true
		}
	}()
	f()
	return
}
