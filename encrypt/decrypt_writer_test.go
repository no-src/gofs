package encrypt

import (
	"bytes"
	"crypto/aes"
	"errors"
	"testing"
)

func TestNewDecryptWriter_CheckIV(t *testing.T) {
	for _, tc := range aesIVTestCases {
		t.Run(tc.iv, func(t *testing.T) {
			_, err := newDecryptWriter(bytes.NewBuffer(nil), []byte(secret), []byte(tc.iv))
			if tc.valid && err != nil {
				t.Errorf("init decrypt writer error, err=%v", err)
				return
			}
			if !tc.valid && !errors.Is(err, errInvalidAESIV) {
				t.Errorf("init decrypt expect get error %v, but get err %v", errInvalidAESIV, err)
			}
		})
	}
}

func TestNewDecryptWriter_CheckKey(t *testing.T) {
	for _, tc := range aesKeyTestCases {
		t.Run(tc.key, func(t *testing.T) {
			_, err := newDecryptWriter(bytes.NewBuffer(nil), []byte(tc.key), aesIV)
			if tc.valid && err != nil {
				t.Errorf("init decrypt writer error, err=%v", err)
				return
			}
			expect := aes.KeySizeError(len(tc.key))
			if !tc.valid && !errors.As(err, &expect) {
				t.Errorf("init decrypt writer expect get error %v, but get err %v", expect, err)
			}
		})
	}
}
