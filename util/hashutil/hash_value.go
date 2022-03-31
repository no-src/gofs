package hashutil

// HashValue the file hash info
type HashValue struct {
	// Offset the file data to calculate the hash value from zero to offset
	Offset int64
	// Hash the file hash value
	Hash string
}

// HashValues the list of *HashValue
type HashValues []*HashValue

// Last returns the last element of HashValues
func (hvs HashValues) Last() *HashValue {
	if len(hvs) > 0 {
		return hvs[len(hvs)-1]
	}
	return nil
}
