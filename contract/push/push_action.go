package push

// PushAction the file upload action
type PushAction int

const (
	// UnknownPushAction the unknown push action
	UnknownPushAction PushAction = iota
	// CompareFilePushAction compare the file hash value before upload the file
	CompareFilePushAction
	// CompareChunkPushAction compare the file chunk hash value before upload the file chunk
	CompareChunkPushAction
	// CompareFileAndChunkPushAction compare the file hash value and first file chunk hash value before upload the file
	CompareFileAndChunkPushAction
	// WritePushAction upload the file or file chunk
	WritePushAction
	// TruncatePushAction truncate the file with the specific size
	TruncatePushAction
)
