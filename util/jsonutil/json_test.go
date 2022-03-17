package jsonutil

import "testing"

func TestMarshal(t *testing.T) {
	jd := jsonData{
		Code:    1,
		Message: "success",
		Secret:  "xyz",
		c:       100,
	}
	data, err := Marshal(jd)
	if err != nil {
		t.Errorf("test Marshal error => %s", err)
		return
	}
	expect := `{"code":1,"message":"success"}`
	actual := string(data)
	if expect != actual {
		t.Errorf("test Marshal error, expect:%s, actual:%s", expect, actual)
	}
}

func TestUnmarshal(t *testing.T) {
	data := []byte(`{"code":1,"message":"success"}`)
	var jd jsonData
	err := Unmarshal(data, &jd)
	if err != nil {
		t.Errorf("test Unmarshal error => %s", err)
		return
	}
	expectCode := 1
	expectMessage := "success"
	if jd.Code != expectCode || jd.Message != expectMessage {
		t.Errorf("test Unmarshal error, expect code:%d, actual code:%d, expect message:%s, actual message:%s", expectCode, jd.Code, expectMessage, jd.Message)
	}
}

type jsonData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Secret  string `json:"-"`
	c       int
}
