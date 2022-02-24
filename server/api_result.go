package server

import "github.com/no-src/gofs/contract"

// ApiResult the common result of api response
type ApiResult struct {
	Code    contract.Code `json:"code"`
	Message string        `json:"message"`
	Data    interface{}   `json:"data"`
}

// NewApiResult create an instance of the ApiResult
func NewApiResult(code contract.Code, message string, data interface{}) ApiResult {
	r := ApiResult{
		Code:    code,
		Message: message,
		Data:    data,
	}
	return r
}

// NewErrorApiResult create an instance of the ApiResult that contains error info
func NewErrorApiResult(code contract.Code, message string) ApiResult {
	return NewApiResult(code, message, nil)
}

// NewServerErrorResult create an instance of the ApiResult that means server error
func NewServerErrorResult() ApiResult {
	return NewErrorApiResult(contract.ServerError, contract.ServerErrorDesc)
}
