package contract

// Status the api response status info
type Status struct {
	// Code current status code
	Code Code `json:"code"`
	// Message current status code description
	Message string `json:"message"`
	// ApiType mark current api type
	ApiType ApiType `json:"api_type"`
}

// SuccessStatus create an instance of success status with specified api type
func SuccessStatus(apiType ApiType) Status {
	return NewStatus(Success, SuccessDesc, apiType)
}

// FailStatus create an instance of fail status with specified api type
func FailStatus(apiType ApiType) Status {
	return NewStatus(Fail, FailDesc, apiType)
}

// UnauthorizedStatus create an instance of unauthorized status with specified api type
func UnauthorizedStatus(apiType ApiType) Status {
	return NewStatus(Unauthorized, UnauthorizedDesc, apiType)
}

// NewStatus create an instance of Status
func NewStatus(code Code, message string, apiType ApiType) Status {
	return Status{
		Code:    code,
		Message: message,
		ApiType: apiType,
	}
}
