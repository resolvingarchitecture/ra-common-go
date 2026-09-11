package util

import (
	"crypto/rand"
	"encoding/binary"
	"math/big"
)

const alphanumeric = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// RandomAlphanumeric returns a cryptographically random alphanumeric string of the given length.
func RandomAlphanumeric(length int) string {
	bytes := RandomBytesOf(length)
	out := make([]byte, length)
	for i := 0; i < length; i++ {
		out[i] = alphanumeric[int(bytes[i])%len(alphanumeric)]
	}
	return string(out)
}

// RandomBytesOf returns count cryptographically random bytes.
func RandomBytesOf(count int) []byte {
	out := make([]byte, count)
	if count > 0 {
		if _, err := rand.Read(out); err != nil {
			panic(err) // crypto/rand failing means the OS entropy source is broken; nothing sensible to return
		}
	}
	return out
}

// NextIntIn returns a random signed 32-bit int in [lower, upper).
func NextIntIn(lower, upper int32) int32 {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(upper-lower)))
	if err != nil {
		panic(err)
	}
	return lower + int32(n.Int64())
}

// NextLong returns a random 64-bit correlation id (used for routing-slip route ids).
func NextLong() int64 {
	buf := RandomBytesOf(8)
	return int64(binary.BigEndian.Uint64(buf))
}
