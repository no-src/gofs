package contract

type Code int

const (
	Unknown      Code = 0
	Success      Code = 1
	Fail         Code = -1
	Unauthorized Code = -2
)

const (
	SuccessDesc      = "success"
	FailDesc         = "fail"
	UnauthorizedDesc = "unauthorized"
)
