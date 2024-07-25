package version

import (
	"fmt"
	"testing"
)

func TestPrintVersion(t *testing.T) {
	commit := Commit
	defer func() {
		Commit = commit
	}()

	testCases := []struct {
		name   string
		commit string
	}{
		{"gofs-test-1", ""},
		{"gofs-test-with-commit", "abcdefg"},
		{"gofs-test-2", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var s string
			Commit = tc.commit
			PrintVersion(tc.name, func(format string, args ...any) {
				s += fmt.Sprintf(format+"\n", args...)
			})
			fmt.Println(s)
		})
	}
}
