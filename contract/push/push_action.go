package push

// PushAction the file upload action
type PushAction int

const (
	// PushActionUnknown the unknown push action
	PushActionUnknown PushAction = iota
	// PushActionCompareFile compare the file hash value before upload the file
	PushActionCompareFile
	// PushActionCompareChunk compare the file chunk hash value before upload the file chunk
	PushActionCompareChunk
	// PushActionCompareFileAndChunk compare the file hash value and first file chunk hash value before upload the file
	PushActionCompareFileAndChunk
	// PushActionWrite upload the file or file chunk
	PushActionWrite
	// PushActionTruncate truncate the file with the specific size
	PushActionTruncate
)
