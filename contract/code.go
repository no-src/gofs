package contract

// Code the status code info
type Code int

const (
	// Unknown the unknown status code
	Unknown Code = 0
	// Success the success status code
	Success Code = 1
	// Fail the standard fail status code
	Fail Code = -1
	// Unauthorized the unauthorized status code, the current user needs to sign in
	Unauthorized Code = -2
	// NotFound the resource not found status code
	NotFound Code = -3
	// NoPermission the no permission status code, the user is authorized but has no permission
	NoPermission Code = -4
	// ServerError the server error status code
	ServerError Code = -5
	// AccessDeny deny current access
	AccessDeny Code = -6
	// Abort the abort code means to ignore and abort the next request
	Abort Code = -7
)

const (
	// SuccessDesc the description of Success code
	SuccessDesc = "success"
	// FailDesc the description of Fail code
	FailDesc = "fail"
	// UnauthorizedDesc the description of Unauthorized code
	UnauthorizedDesc = "unauthorized"
	// NotFoundDesc the description of NotFound code
	NotFoundDesc = "not found"
	// NoPermissionDesc the description of NoPermission code
	NoPermissionDesc = "no permission"
	// ServerErrorDesc the description of ServerError code
	ServerErrorDesc = "server internal error"
	// AccessDenyDesc the description of AccessDeny code
	AccessDenyDesc = "access deny"
	// AbortDesc the description of Abort code
	AbortDesc = "abort"
)
