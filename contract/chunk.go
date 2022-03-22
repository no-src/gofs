package contract

// Chunk the file chunk info that is used to upload the file
type Chunk struct {
	// Offset the offset relative to the origin of the file
	//
	// The offset less than zero is invalid.
	// The offset equals zero means the first chunk or only one chunk.
	// The offset greater than zero means that have two chunks at least.
	Offset int64 `json:"offset"`

	// Hash calculate the file chunk hash value, if the path is a file
	Hash string `json:"hash"`

	// Size the size of file chunk for bytes
	Size int64 `json:"size"`
}
