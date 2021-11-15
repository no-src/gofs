package server

import "encoding/json"

type ApiResult struct {
	Code    int
	Message string
	Data    interface{}
}

func NewApiResult(code int, message string, data interface{}) ApiResult {
	r := ApiResult{
		Code:    code,
		Message: message,
		Data:    data,
	}
	return r
}

func NewApiResultBytes(code int, message string, data interface{}) []byte {
	r := NewApiResult(code, message, data)
	bytes, err := json.Marshal(r)
	if err != nil {
		return nil
	}
	return bytes
}

func NewErrorApiResultBytes(code int, message string) []byte {
	return NewApiResultBytes(code, message, nil)
}
