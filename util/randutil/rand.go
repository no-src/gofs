package randutil

import (
	"crypto/rand"
	mathrand "math/rand"
	"time"
)

var (
	read      = rand.Read
	innerRand = mathrand.New(mathrand.NewSource(time.Now().UnixNano()))
)

// RandomString generate a random string, max length is 20
func RandomString(length int) (s string) {
	max := 20
	if length > max {
		length = max
	}
	bytes := make([]byte, 32)
	_, err := read(bytes)
	if err != nil {
		for i := 0; i < 32; i++ {
			bytes[i] = byte(innerRand.Intn(256))
		}
	}
	var h [32]byte
	copy(h[:32], bytes)
	s = hashToString(h)
	s = s[:length]
	return s
}

func hashToString(h [32]byte) string {
	const b64 = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
	const chunks = 5
	var dst [chunks * 4]byte
	for i := 0; i < chunks; i++ {
		v := uint32(h[3*i])<<16 | uint32(h[3*i+1])<<8 | uint32(h[3*i+2])
		dst[4*i+0] = b64[(v>>18)&0x3F]
		dst[4*i+1] = b64[(v>>12)&0x3F]
		dst[4*i+2] = b64[(v>>6)&0x3F]
		dst[4*i+3] = b64[v&0x3F]
	}
	return string(dst[:])
}
