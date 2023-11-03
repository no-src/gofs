package flag

import (
	"os"
	"testing"

	"github.com/no-src/nsgo/hashutil"
)

func TestParseFlags(t *testing.T) {
	testCases := []struct {
		name   string
		args   []string
		expect string
	}{
		{"with default value", []string{os.Args[0], "-server_addr=:443"}, hashutil.DefaultHash},
		{"with empty value", []string{os.Args[0], "-checksum_algorithm="}, ""},
		{"with sha512", []string{os.Args[0], "-checksum_algorithm=" + hashutil.SHA512Hash}, hashutil.SHA512Hash},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := ParseFlags(tc.args)
			if c.ChecksumAlgorithm != tc.expect {
				t.Errorf("parse flags error, checksum_algorithm expect to get %s, but actual get %s", tc.expect, c.ChecksumAlgorithm)
			}
		})
	}
}

func TestParseFlags_PanicWithZeroArgument(t *testing.T) {
	defer func() {
		e := recover()
		if e == nil {
			t.Errorf("parse the flags with zero argument error, expect to panic but not")
		}
	}()
	ParseFlags([]string{})
}

func TestParseFlags_PanicWithOneArgument(t *testing.T) {
	defer func() {
		e := recover()
		if e == nil {
			t.Errorf("parse the flags with one argument error, expect to panic but not")
		}
	}()
	ParseFlags(os.Args[:1])
}
