//go:build encrypt_test

package encrypt

import (
	"archive/zip"
	"crypto/aes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/no-src/gofs/conf"
	"github.com/no-src/gofs/logger"
)

func TestDecrypt_DecryptOutNotDir(t *testing.T) {
	logger := logger.NewTestLogger()
	defer logger.Close()

	decryptOpt := NewOption(conf.Config{
		Decrypt:       true,
		DecryptPath:   decryptPath,
		DecryptSecret: secret,
		DecryptOut:    "./encrypt_test.go",
	}, logger)
	dec, err := NewDecrypt(decryptOpt)
	if err != nil {
		t.Errorf("init decrypt component error => %v", err)
		return
	}
	err = dec.Decrypt()
	if !errors.Is(err, errDecryptOutNotDir) {
		t.Errorf("expect get error => [%v] but get [%v]", errDecryptOutNotDir, err)
	}
}

func TestDecrypt_EvilFile(t *testing.T) {
	evilFile := "./testdata/zipslip.zip"
	err := os.MkdirAll(filepath.Dir(evilFile), fs.ModePerm)
	if err != nil {
		t.Errorf("mkdir evil path error err=%v", err)
		return
	}
	f, err := os.Create(evilFile)
	if err != nil {
		t.Errorf("create evil file error err=%v", err)
		return
	}
	zw := zip.NewWriter(f)
	_, err = zw.CreateHeader(&zip.FileHeader{
		Name:   "../zipslip/",
		Method: zip.Store,
	})
	if err != nil {
		t.Errorf("create zip file error err=%v", err)
		return
	}
	err = zw.Flush()
	if err != nil {
		t.Errorf("flush zip error err=%v", err)
		return
	}
	err = zw.Close()
	if err != nil {
		t.Errorf("close zip writer error err=%v", err)
		return
	}
	err = f.Close()
	if err != nil {
		t.Errorf("close file error err=%v", err)
		return
	}

	logger := logger.NewTestLogger()
	defer logger.Close()

	decryptOpt := NewOption(conf.Config{
		Decrypt:       true,
		DecryptPath:   evilFile,
		DecryptSecret: secret,
		DecryptOut:    decryptOut,
	}, logger)
	dec, err := NewDecrypt(decryptOpt)
	if err != nil {
		t.Errorf("init decrypt component error => %v", err)
		return
	}
	err = dec.Decrypt()
	if !errors.Is(err, errIllegalPath) {
		t.Errorf("expect get error => [%v] but get [%v]", errIllegalPath, err)
	}
}

func TestNewDecrypt_CheckKey(t *testing.T) {
	logger := logger.NewTestLogger()
	defer logger.Close()

	for _, tc := range aesKeyTestCases {
		t.Run(tc.key, func(t *testing.T) {
			decryptOpt := NewOption(conf.Config{
				Decrypt:       true,
				DecryptPath:   decryptPath,
				DecryptSecret: tc.key,
				DecryptOut:    decryptOut,
			}, logger)

			_, err := NewDecrypt(decryptOpt)
			if tc.valid && err != nil {
				t.Errorf("init decrypt component error, err=%v", err)
				return
			}
			expect := aes.KeySizeError(len(tc.key))
			if !tc.valid && !errors.As(err, &expect) {
				t.Errorf("init decrypt expect get error %v, but get err %v", expect, err)
			}
		})
	}
}
