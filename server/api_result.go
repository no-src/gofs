package server

import (
	"github.com/no-src/gofs/util"
)

type ApiResult struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
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
	bytes, err := util.Marshal(r)
	if err != nil {
		return nil
	}
	return bytes
}

func NewErrorApiResultBytes(code int, message string) []byte {
	return NewApiResultBytes(code, message, nil)
}
