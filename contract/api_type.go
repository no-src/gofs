package contract

type ApiType int

const (
	UnknownApi     ApiType = 0
	SyncMessageApi ApiType = 1
	InfoApi        ApiType = 2
)
