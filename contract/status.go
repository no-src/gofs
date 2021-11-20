package contract

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
	return Status{
		Code:    Success,
		Message: SuccessDesc,
		ApiType: apiType,
	}
}
