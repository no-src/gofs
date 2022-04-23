package auth

import (
	"bytes"
	"testing"

	"github.com/no-src/gofs/util/hashutil"
)

func TestParseAuthCommandData(t *testing.T) {
	testCases := []struct {
		name     string
		authData []byte
		expect   string
	}{
		{"username 111", append(append([]byte("auth"), authVersion...), []byte("698d51a19d8a121cbcbe3365e6ac95ea20220222072118")...), hashutil.MD5FromString("111")[:userNameHashLength]},
		{"username root", append(append([]byte("auth"), authVersion...), []byte("63a9f0ea7bb980507b24afc8bc80e54820220222072947")...), hashutil.MD5FromString("root")[:userNameHashLength]},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u, err := ParseAuthCommandData(tc.authData)
			if err != nil {
				t.Errorf("authData => %s error= %s", string(tc.authData), err.Error())
				return
			}
			if u == nil {
				t.Errorf("get a nil user => %s", string(tc.authData))
				return
			}
			actual := u.UserNameHash
			if actual != tc.expect {
				t.Errorf("[%s] => expect: %v, but actual: %v", tc.authData, tc.expect, actual)
			}
		})
	}
}

func TestParseAuthCommandData_ReturnError(t *testing.T) {
	authData := func() []byte {
		return append(append([]byte("auth"), authVersion...), []byte("63a9f0ea7bb980507b24afc8bc80e54820220222072947")...)
	}
	testCases := []struct {
		name     string
		authData []byte
	}{
		{"auth data too short length", authData()[1:]},
		{"auth data nil", nil},
		{"auth data empty", []byte("")},
		{"auth data too long length", append(authData(), []byte("x")...)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseAuthCommandData(tc.authData)
			if err == nil {
				t.Errorf("expect to get an error but get nil => %s", string(tc.authData))
				return
			}
		})
	}
}

func TestGenerateAuthCommandData(t *testing.T) {
	testCases := []struct {
		name   string
		user   *HashUser
		expect []byte
	}{
		{"698d51a19d8a121c bcbe3365e6ac95ea", NewHashUser("698d51a19d8a121c", "bcbe3365e6ac95ea", DefaultPerm), append(append([]byte("auth"), authVersion...), []byte("698d51a19d8a121cbcbe3365e6ac95ea20220222072118")...)},
		{"63a9f0ea7bb98050 7b24afc8bc80e548", NewHashUser("63a9f0ea7bb98050", "7b24afc8bc80e548", DefaultPerm), append(append([]byte("auth"), authVersion...), []byte("63a9f0ea7bb980507b24afc8bc80e54820220222072947")...)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := GenerateAuthCommandData(tc.user)
			if len(actual) != len(tc.expect) || !bytes.Equal(tc.expect[:len(tc.expect)-expiresLength], actual[:len(actual)-expiresLength]) {
				t.Errorf("[%s] => expect: %v, but actual: %v", tc.user.UserNameHash, tc.expect, actual)
			}
		})
	}

}

func TestGenerateAuthCommandData_WithNilUser(t *testing.T) {
	actual := GenerateAuthCommandData(nil)
	if actual != nil {
		t.Errorf("should be return nil")
	}
}
