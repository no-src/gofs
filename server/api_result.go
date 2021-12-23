package server

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

func NewErrorApiResult(code int, message string) ApiResult {
	return NewApiResult(code, message, nil)
}