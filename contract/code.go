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
	// NotModified  the resource is not modified
	NotModified Code = -7
	// ChunkNotModified the chunk is not modified
	ChunkNotModified Code = -8
	// Modified  the resource is modified
	Modified Code = -9
	// ChunkModified the chunk is modified
	ChunkModified Code = -10
)

const (
	// UnknownDesc the description of Unknown code
	UnknownDesc = "unknown"
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
	// NotModifiedDesc the description of NotModified code
	NotModifiedDesc = "not modified"
	// ChunkNotModifiedDesc the description of ChunkNotModified code
	ChunkNotModifiedDesc = "chunk not modified"
	// ModifiedDesc the description of Modified code
	ModifiedDesc = "modified"
	// ChunkModifiedDesc the description of ChunkModified code
	ChunkModifiedDesc = "chunk modified"
)

// String return the code description name
func (code Code) String() string {
	desc := ""
	switch code {
	case Unknown:
		desc = UnknownDesc
	case Success:
		desc = SuccessDesc
	case Fail:
		desc = FailDesc
	case Unauthorized:
		desc = UnauthorizedDesc
	case NotFound:
		desc = NotFoundDesc
	case NoPermission:
		desc = NoPermissionDesc
	case ServerError:
		desc = ServerErrorDesc
	case AccessDeny:
		desc = AccessDenyDesc
	case NotModified:
		desc = NotModifiedDesc
	case ChunkNotModified:
		desc = ChunkNotModifiedDesc
	case Modified:
		desc = ModifiedDesc
	case ChunkModified:
		desc = ChunkModifiedDesc
	default:
		desc = UnknownDesc
	}
	return desc
}
