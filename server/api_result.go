package server

import "github.com/no-src/gofs/contract"

type ApiResult struct {
	Code    contract.Code `json:"code"`
	Message string        `json:"message"`
	Data    interface{}   `json:"data"`
}

func NewApiResult(code contract.Code, message string, data interface{}) ApiResult {
	r := ApiResult{
		Code:    code,
		Message: message,
		Data:    data,
	}
	return r
}

func NewErrorApiResult(code contract.Code, message string) ApiResult {
	return NewApiResult(code, message, nil)
}

func NewServerErrorResult() ApiResult {
	return NewErrorApiResult(contract.ServerError, contract.ServerErrorDesc)
}
