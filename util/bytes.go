package util

import "encoding/binary"

// PackBigEndian interprets the first four bytes of b as a big-endian signed 32-bit int.
func PackBigEndian(b []byte) int32 {
	return int32(binary.BigEndian.Uint32(b))
}

// UnpackBigEndian encodes x as four big-endian bytes.
func UnpackBigEndian(x int32) []byte {
	out := make([]byte, 4)
	binary.BigEndian.PutUint32(out, uint32(x))
	return out
}

// PackLittleEndian interprets the first four bytes of b as a little-endian signed 32-bit int.
func PackLittleEndian(b []byte) int32 {
	return int32(binary.LittleEndian.Uint32(b))
}

// UnpackLittleEndian encodes x as four little-endian bytes.
func UnpackLittleEndian(x int32) []byte {
	out := make([]byte, 4)
	binary.LittleEndian.PutUint32(out, uint32(x))
	return out
}
