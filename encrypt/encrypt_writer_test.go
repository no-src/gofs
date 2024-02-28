package encrypt

import (
	"bytes"
	"crypto/aes"
	"errors"
	"strings"
	"testing"
)

func TestNewEncryptWriter_CheckIV(t *testing.T) {
	for _, tc := range aesIVTestCases {
		t.Run(tc.iv, func(t *testing.T) {
			_, err := newEncryptWriter(bytes.NewBuffer(nil), "test.txt", []byte(secret), []byte(tc.iv))
			if tc.valid && err != nil {
				t.Errorf("init encrypt writer error, err=%v", err)
				return
			}
			if !tc.valid && !errors.Is(err, errInvalidAESIV) {
				t.Errorf("init encrypt expect get error %v, but get err %v", errInvalidAESIV, err)
			}
		})
	}
}

func TestNewEncryptWriter_CheckKey(t *testing.T) {
	for _, tc := range aesKeyTestCases {
		t.Run(tc.key, func(t *testing.T) {
			_, err := newEncryptWriter(bytes.NewBuffer(nil), "test.txt", []byte(tc.key), aesIV)
			if tc.valid && err != nil {
				t.Errorf("init encrypt writer error, err=%v", err)
				return
			}
			expect := aes.KeySizeError(len(tc.key))
			if !tc.valid && !errors.As(err, &expect) {
				t.Errorf("init encrypt writer expect get error %v, but get err %v", expect, err)
			}
		})
	}
}

func TestNewEncryptWriter_CreateHeaderError(t *testing.T) {
	_, err := newEncryptWriter(bytes.NewBuffer(nil), strings.Repeat("x", 1<<16), []byte(secret), aesIV)
	if err == nil {
		t.Errorf("init encrypt writer expect get an error, but get nil")
		return
	}
	expect := "zip: FileHeader.Name too long"
	if !strings.Contains(err.Error(), expect) {
		t.Errorf("init encrypt writer expect get error %s, but get err %v", expect, err)
	}
}
