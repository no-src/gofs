package hashutil

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"hash/adler32"
	"hash/crc32"
	"hash/crc64"
	"hash/fnv"
	"strings"
)

var (
	factories map[string]hashFactory
)

type hashFactory func() hash.Hash

// NewHash return the hash implementation by the specified hash algorithm
func NewHash(algorithm string) (Hash, error) {
	f, err := getFactory(algorithm)
	if err != nil {
		return nil, err
	}
	return &defaultHash{factory: f}, nil
}

func getFactory(algorithm string) (hashFactory, error) {
	algorithm = strings.ToLower(algorithm)
	f, ok := factories[algorithm]
	if ok {
		return f, nil
	}
	return nil, fmt.Errorf("unsupported hash algorithm => %s", algorithm)
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

func register() {
	factories = map[string]hashFactory{
		MD5Hash:    md5.New,
		SHA1Hash:   sha1.New,
		SHA256Hash: sha256.New,
		SHA512Hash: sha512.New,
		CRC32Hash: func() hash.Hash {
			return crc32.NewIEEE()
		},
		CRC64Hash: func() hash.Hash {
			return crc64.New(crc64.MakeTable(crc64.ISO))
		},
		Adler32Hash: func() hash.Hash {
			return adler32.New()
		},
		FNV132Hash: func() hash.Hash {
			return fnv.New32()
		},
		FNV1A32Hash: func() hash.Hash {
			return fnv.New32a()
		},
		FNV164Hash: func() hash.Hash {
			return fnv.New64()
		},
		FNV1A64Hash: func() hash.Hash {
			return fnv.New64a()
		},
		FNV1128Hash:  fnv.New128,
		FNV1A128Hash: fnv.New128a,
	}
}

func init() {
	register()
}
