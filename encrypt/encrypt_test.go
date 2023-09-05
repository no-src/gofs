//go:build encrypt_test

package encrypt

import (
	"bytes"
	"crypto/aes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/no-src/gofs/conf"
)

var (
	secret          = "encrypt_secure_encrypt_secure_01"
	sourcePath      = "./"
	encryptPath     = "./"
	decryptPath     = "./testdata/encrypt"
	decryptOut      = "./testdata/decrypt_out"
	originPath      = "./encrypt_test.go"
	encryptFilePath = "./testdata/encrypt/encrypt_test.go.data"
)

func TestEncrypt(t *testing.T) {
	encryptOpt := NewOption(conf.Config{
		Encrypt:       true,
		EncryptPath:   encryptPath,
		EncryptSecret: secret,
	})

	decryptOpt := NewOption(conf.Config{
		Decrypt:       true,
		DecryptPath:   decryptPath,
		DecryptSecret: secret,
		DecryptOut:    decryptOut,
	})

	err := testEncrypt(encryptOpt, decryptOpt, sourcePath, originPath, encryptFilePath)
	if err != nil {
		t.Errorf("test encrypt and decrypt error err=%v", err)
	}
}

func TestEncrypt_Disabled(t *testing.T) {
	encryptOpt := NewOption(conf.Config{
		Encrypt:       true,
		EncryptPath:   encryptPath,
		EncryptSecret: secret,
	})

	decryptOpt := NewOption(conf.Config{
		Decrypt:       true,
		DecryptPath:   decryptPath,
		DecryptSecret: secret,
		DecryptOut:    decryptOut,
	})

	// disable encryption and decryption
	encryptOpt.Encrypt = false
	decryptOpt.Decrypt = false

	err := testEncrypt(encryptOpt, decryptOpt, sourcePath, originPath, encryptFilePath)
	if err != nil {
		t.Errorf("test encrypt and decrypt error err=%v", err)
	}
}

func TestEncrypt_NotSubPath(t *testing.T) {
	encryptOpt := NewOption(conf.Config{
		Encrypt:       true,
		EncryptPath:   "../",
		EncryptSecret: secret,
	})

	decryptOpt := NewOption(conf.Config{
		Decrypt:       true,
		DecryptPath:   decryptPath,
		DecryptSecret: secret,
		DecryptOut:    decryptOut,
	})

	err := testEncrypt(encryptOpt, decryptOpt, sourcePath, originPath, encryptFilePath)
	if !errors.Is(err, errNotSubDir) {
		t.Errorf("expect get error => [%v] but get [%v]", errNotSubDir, err)
	}
}

func TestEncrypt_EmptyOption(t *testing.T) {
	encryptOpt := NewOption(conf.Config{
		Encrypt:       false,
		EncryptPath:   encryptPath,
		EncryptSecret: secret,
	})

	if len(encryptOpt.EncryptPath) != 0 || len(encryptOpt.EncryptSecret) != 0 {
		t.Errorf("expect to get an empty option but not")
		return
	}

	decryptOpt := NewOption(conf.Config{
		Decrypt:       false,
		DecryptPath:   decryptPath,
		DecryptSecret: secret,
		DecryptOut:    decryptOut,
	})

	if len(decryptOpt.EncryptPath) != 0 || len(decryptOpt.EncryptSecret) != 0 || len(decryptOpt.DecryptOut) != 0 {
		t.Errorf("expect to get an empty option but not")
		return
	}
}

func TestNewEncrypt_CheckKey(t *testing.T) {
	for _, tc := range aesKeyTestCases {
		t.Run(tc.key, func(t *testing.T) {
			encryptOpt := NewOption(conf.Config{
				Encrypt:       true,
				EncryptPath:   encryptPath,
				EncryptSecret: tc.key,
			})

			_, err := NewEncrypt(encryptOpt, sourcePath)
			if tc.valid && err != nil {
				t.Errorf("init encrypt component error, err=%v", err)
				return
			}
			expect := aes.KeySizeError(len(tc.key))
			if !tc.valid && !errors.As(err, &expect) {
				t.Errorf("init encrypt expect get error %v, but get err %v", expect, err)
			}
		})
	}
}

//gocyclo:ignore
func testEncrypt(encryptOpt Option, decryptOpt Option, sourcePath string, originPath string, encryptFilePath string) error {
	// encrypt
	enc, err := NewEncrypt(encryptOpt, sourcePath)
	if err != nil {
		return err
	}

	originFile, err := os.Open(originPath)
	if err != nil {
		return err
	}
	defer originFile.Close()

	err = os.MkdirAll(filepath.Dir(encryptFilePath), fs.ModePerm)
	if err != nil {
		return err
	}
	encFile, err := os.Create(encryptFilePath)
	if err != nil {
		return err
	}
	defer encFile.Close()

	encStat, err := encFile.Stat()
	if err != nil {
		return err
	}
	decFileName := encStat.Name()

	w, err := enc.NewWriter(encFile, originPath, decFileName)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, originFile)
	if err != nil {
		return err
	}
	w.Close()

	// encrypt temp
	_, removeFunc, err := enc.CreateEncryptTemp(originPath)
	if err != nil {
		return err
	}
	err = removeFunc()
	if err != nil {
		return err
	}

	// decrypt
	dec, err := NewDecrypt(decryptOpt)
	if err != nil {
		return err
	}
	err = dec.Decrypt()
	if err != nil {
		return err
	}

	// check result
	if !encryptOpt.Encrypt || !decryptOpt.Decrypt {
		return nil
	}
	originFile.Seek(0, 0)
	originContent, err := io.ReadAll(originFile)
	if err != nil {
		return err
	}

	decContent, err := os.ReadFile(filepath.Join(decryptOpt.DecryptOut, decFileName))
	if err != nil {
		return err
	}
	if bytes.Compare(originContent, decContent) != 0 {
		return fmt.Errorf("the origin file not equals to decrypt file")
	}
	return nil
}
