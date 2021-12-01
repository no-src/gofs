package contract

type ApiType int

const (
	// UnknownApi unknown api type
	UnknownApi ApiType = 0
	// SyncMessageApi send remote sync message api
	SyncMessageApi ApiType = 1
	// InfoApi query file server info api
	InfoApi ApiType = 2
	// AuthApi tcp server auth api
	AuthApi ApiType = 3
)
