package auth

import (
	"bytes"
	"github.com/no-src/gofs/util/hashutil"
	"testing"
)

func TestParseAuthCommandData(t *testing.T) {
	expect := hashutil.MD5FromString("111")
	authData := append(append([]byte("auth"), authVersion...), []byte("698d51a19d8a121cbcbe3365e6ac95ea20220222072118")...)
	testParseAuthCommandData(t, expect[:userNameHashLength], authData)

	expect = hashutil.MD5FromString("root")
	authData = append(append([]byte("auth"), authVersion...), []byte("63a9f0ea7bb980507b24afc8bc80e54820220222072947")...)
	testParseAuthCommandData(t, expect[:userNameHashLength], authData)

	testParseAuthCommandDataFail(t, authData[1:])
	testParseAuthCommandDataFail(t, nil)
	testParseAuthCommandDataFail(t, []byte(""))
	testParseAuthCommandDataFail(t, append(authData, []byte("x")...))
}

func testParseAuthCommandData(t *testing.T, expect string, authData []byte) {
	u, err := ParseAuthCommandData(authData)
	if err != nil {
		t.Errorf("ParseAuthCommandData error => %s error= %s", string(authData), err.Error())
		return
	}
	if u == nil {
		t.Errorf("ParseAuthCommandData error, get a nil user => %s", string(authData))
		return
	}
	actual := u.UserNameHash
	if actual != expect {
		t.Errorf("[%s] => expect: %v, but actual: %v \n", authData, expect, actual)
	}
}

func testParseAuthCommandDataFail(t *testing.T, authData []byte) {
	_, err := ParseAuthCommandData(authData)
	if err == nil {
		t.Errorf("ParseAuthCommandData shuold be error => %s", string(authData))
		return
	}
}

func TestGenerateAuthCommandData(t *testing.T) {
	authData := append(append([]byte("auth"), authVersion...), []byte("698d51a19d8a121cbcbe3365e6ac95ea20220222072118")...)
	testGenerateAuthCommandData(t, authData, NewHashUser("698d51a19d8a121c", "bcbe3365e6ac95ea", DefaultPerm))

	authData = append(append([]byte("auth"), authVersion...), []byte("63a9f0ea7bb980507b24afc8bc80e54820220222072947")...)
	testGenerateAuthCommandData(t, authData, NewHashUser("63a9f0ea7bb98050", "7b24afc8bc80e548", DefaultPerm))
}

func testGenerateAuthCommandData(t *testing.T, expect []byte, user *HashUser) {
	actual := GenerateAuthCommandData(user)
	if len(actual) != len(expect) || !bytes.Equal(expect[:len(expect)-expiresLength], actual[:len(actual)-expiresLength]) {
		t.Errorf("[%s] => expect: %v, but actual: %v \n", user.UserNameHash, expect, actual)
	}
}

func TestGenerateAuthCommandDataNilUser(t *testing.T) {
	actual := GenerateAuthCommandData(nil)
	if actual != nil {
		t.Errorf("GenerateAuthCommandData should be return nil")
	}
}
