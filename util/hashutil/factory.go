package hashutil

import (
	"crypto/md5"
	"hash"
)

var (
	factory hashFactory
)

type hashFactory func() hash.Hash

// InitDefaultHash initial default hash factory
func InitDefaultHash(algorithm string) {
	switch algorithm {
	case MD5Hash:
		factory = func() hash.Hash {
			return md5.New()
		}
	}
}

// New return default hash implementation
func New() hash.Hash {
	return factory()
}

const (
	// MD5Hash the MD5 hash algorithm
	MD5Hash = "md5"
)

func init() {
	InitDefaultHash(MD5Hash)
}
