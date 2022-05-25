package hashutil

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"hash/adler32"
	"hash/crc32"
	"hash/crc64"
	"hash/fnv"
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
	case SHA1Hash:
		factory = func() hash.Hash {
			return sha1.New()
		}
	case SHA256Hash:
		factory = func() hash.Hash {
			return sha256.New()
		}
	case SHA512Hash:
		factory = func() hash.Hash {
			return sha512.New()
		}
	case CRC32Hash:
		factory = func() hash.Hash {
			return crc32.NewIEEE()
		}
	case CRC64Hash:
		factory = func() hash.Hash {
			return crc64.New(crc64.MakeTable(crc64.ISO))
		}
	case Adler32Hash:
		factory = func() hash.Hash {
			return adler32.New()
		}
	case FNV132Hash:
		factory = func() hash.Hash {
			return fnv.New32()
		}
	case FNV1A32Hash:
		factory = func() hash.Hash {
			return fnv.New32a()
		}
	case FNV164Hash:
		factory = func() hash.Hash {
			return fnv.New64()
		}
	case FNV1A64Hash:
		factory = func() hash.Hash {
			return fnv.New64a()
		}
	case FNV1128Hash:
		factory = func() hash.Hash {
			return fnv.New128()
		}
	case FNV1A128Hash:
		factory = func() hash.Hash {
			return fnv.New128a()
		}
	}
}

// New return default hash implementation
func New() hash.Hash {
	return factory()
}

const (
	// DefaultHash the default hash algorithm
	DefaultHash = MD5Hash
	// MD5Hash the MD5 hash algorithm
	MD5Hash = "md5"
	// SHA1Hash the SHA-1 hash algorithm
	SHA1Hash = "sha1"
	// SHA256Hash the SHA256 hash algorithm
	SHA256Hash = "sha256"
	// SHA512Hash the SHA-512 hash algorithm
	SHA512Hash = "sha512"
	// CRC32Hash the CRC-32 checksum
	CRC32Hash = "crc32"
	// CRC64Hash the CRC-64 checksum
	CRC64Hash = "crc64"
	// Adler32Hash the Adler-32 checksum
	Adler32Hash = "adler32"
	// FNV132Hash the 32-bit FNV-1 non-cryptographic hash function
	FNV132Hash = "fnv-1-32"
	// FNV1A32Hash the 32-bit FNV-1a non-cryptographic hash function
	FNV1A32Hash = "fnv-1a-32"
	// FNV164Hash the 64-bit FNV-1 non-cryptographic hash function
	FNV164Hash = "fnv-1-64"
	// FNV1A64Hash the 64-bit FNV-1a non-cryptographic hash function
	FNV1A64Hash = "fnv-1a-64"
	// FNV1128Hash the 128-bit FNV-1 non-cryptographic hash function
	FNV1128Hash = "fnv-1-128"
	// FNV1A128Hash the 128-bit FNV-1a non-cryptographic hash function
	FNV1A128Hash = "fnv-1a-128"
)

func init() {
	InitDefaultHash(DefaultHash)
}
