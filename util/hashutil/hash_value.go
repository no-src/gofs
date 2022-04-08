package hashutil

// HashValue the file hash info
type HashValue struct {
	// Offset the file data to calculate the hash value from zero to offset
	Offset int64 `json:"offset"`
	// Hash the file hash value
	Hash string `json:"hash"`
}

// NewHashValue returns an instance of HashValue
func NewHashValue(offset int64, hash string) *HashValue {
	return &HashValue{
		Offset: offset,
		Hash:   hash,
	}
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
