package push

const (
	// FileInfo basic push file info
	FileInfo = "file_info"
	// UpFile the field name of upload file
	UpFile = "up_file"
	// Offset the offset relative to the origin of the file.
	//
	// The offset less than zero means to compare file size and hash value only.
	// The offset equals zero means the first chunk or only one chunk.
	// The offset greater than zero means that have two chunks at least.
	Offset = "offset"
)
