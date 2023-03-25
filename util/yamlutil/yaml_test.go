package yamlutil

import "testing"

func TestMarshal(t *testing.T) {
	jd := getTestData()
	d, err := Marshal(jd)
	if err != nil {
		t.Errorf("test Marshal error => %s", err)
		return
	}
	expect := "code: 1\nmessage: success\n"
	actual := string(d)
	if expect != actual {
		t.Errorf("test Marshal error, expect:%s, actual:%s", expect, actual)
	}
}

func TestUnmarshal(t *testing.T) {
	d := []byte("code: 1\nmessage: success\n")
	var yd data
	err := Unmarshal(d, &yd)
	if err != nil {
		t.Errorf("test Unmarshal error => %s", err)
		return
	}
	expectCode := 1
	expectMessage := "success"
	expectSecret := ""
	expectC := 0
	if yd.Code != expectCode || yd.Message != expectMessage || yd.Secret != expectSecret || yd.c != expectC {
		t.Errorf("test Unmarshal error, expect code:%d, actual code:%d, expect message:%s, actual message:%s, expect secret:%s, actual secret:%s, expect c:%d, actual c:%d", expectCode, yd.Code, expectMessage, yd.Message, expectSecret, yd.Secret, expectC, yd.c)
	}
}

type data struct {
	Code    int    `yaml:"code"`
	Message string `yaml:"message"`
	Secret  string `yaml:"-"`
	c       int
}

func getTestData() data {
	return data{
		Code:    1,
		Message: "success",
		Secret:  "xyz",
		c:       100,
	}
}
