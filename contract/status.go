package contract

type Status struct {
	Code    Code    `json:"code"`
	Message string  `json:"message"`
	ApiType ApiType `json:"api_type"`
}

func SuccessStatus(apiType ApiType) Status {
	return Status{
		Code:    Success,
		Message: SuccessDesc,
		ApiType: apiType,
	}
}
