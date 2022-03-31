package hashutil

// HashValue the file hash info
type HashValue struct {
	// Offset the file data to calculate the hash value from zero to offset
	Offset int64
	// Hash the file hash value
	Hash string
}
