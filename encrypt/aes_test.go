package encrypt

import (
	"crypto/aes"
	"errors"
	"testing"
)

var aesKeyTestCases = []struct {
	key   string
	valid bool
}{
	{"", false},
	{"1", false},
	{"123456789012345", false},
	{"1234567890123456", true},
	{"12345678901234567", false},
	{"12345678901234567890123", false},
	{"123456789012345678901234", true},
	{"1234567890123456789012345", false},
	{"1234567890123456789012345678901", false},
	{"12345678901234567890123456789012", true},
	{"123456789012345678901234567890123", false},
}

var aesIVTestCases = []struct {
	iv    string
	valid bool
}{
	{"", false},
	{"1", false},
	{"123456789012345", false},
	{"1234567890123456", true},
	{"12345678901234567", false},
}

func TestCheckAESKey(t *testing.T) {
	for _, tc := range aesKeyTestCases {
		t.Run(tc.key, func(t *testing.T) {
			err := checkAESKey([]byte(tc.key))
			if tc.valid && err != nil {
				t.Errorf("check AES key error, err=%v", err)
				return
			}
			expect := aes.KeySizeError(len(tc.key))
			if !tc.valid && !errors.As(err, &expect) {
				t.Errorf("check AES key expect get error %v, but get err %v", expect, err)
			}
		})
	}
}

func TestCheckAESIV(t *testing.T) {
	for _, tc := range aesIVTestCases {
		t.Run(tc.iv, func(t *testing.T) {
			err := checkAESIV([]byte(tc.iv))
			if tc.valid && err != nil {
				t.Errorf("check AES IV error, err=%v", err)
				return
			}
			if !tc.valid && !errors.Is(err, errInvalidAESIV) {
				t.Errorf("check AES IV expect get error %v, but get err %v", errInvalidAESIV, err)
			}
		})
	}
}
