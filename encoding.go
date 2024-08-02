package main

import (
	"errors"
	"fmt"
	"math/bits"
)

var (
	ErrInvalidHighBits = errors.New("invalid high bits")
	ErrInvalidLength = errors.New("invalid length")
	ErrInsufficientLength = errors.New("insufficient length")
)

// Shockwave encoding

// VL64EncodeLen returns the number of bytes required to represent a variable-length base64-encoded integer.
func VL64EncodeLen(v int) int {
	if v < 0 {
		v *= -1
	}
	return (bits.Len32(uint32(v)) + 9) / 6
}

// VL64DecodeLen returns the byte length of a variable-length base64-encoded integer, given (and including) the first byte.
func VL64DecodeLen(b byte) (length int, err error) {
	if b&0xc0 != 0x40 {
		err = fmt.Errorf("%w in vl64: %q (0x%02[2]x)", ErrInvalidHighBits, b)
		return
	}
	length = int(b >> 3 & 7)
	if length == 0 || length > 6 {
		err = fmt.Errorf("%w in vl64: %q (%d)", ErrInvalidLength, b, length)
	}
	return
}

// VL64Encode encodes an integer to variable-length base64 into the specified byte slice.
func VL64Encode(b []byte, v int) {
	abs := v
	if abs < 0 {
		abs *= -1
	}
	n := VL64EncodeLen(v)

	b[0] = 64 | (byte(n)&7)<<3 | byte(abs&3)
	if v < 0 {
		b[0] |= 4
	}
	for i := 1; i < n; i++ {
		b[i] = 64 | byte((abs>>(2+6*(i-1)))&0x3f)
	}
}

// VL64Decode decodes a variable-length base64-encoded integer from the specified byte slice.
func VL64Decode(b []byte) (value, n int, err error) {
	value = int(b[0] & 3)

	n, err = VL64DecodeLen(b[0])
	if len(b) < n {
		err = fmt.Errorf("%w: need %d bytes, have %d", ErrInsufficientLength, n, len(b))
		return
	}
	for i := 1; i < n; i++ {
		if b[i]&0xc0 != 0x40 {
			err = fmt.Errorf("%w in vl64: %q (0x%02[2]x)", ErrInvalidHighBits, b[i])
			return
		}
		value |= int(b[i]&0x3f) << (2 + 6*(i-1))
	}

	if b[0]&4 != 0 {
		value *= -1
	}
	return
}

// B64Encode encodes an integer to base64 into the specified byte slice.
func B64Encode(b []byte, v int) {
	for i := 0; i < len(b); i++ {
		b[i] = 64 | byte(v>>((len(b)-i-1)*6)&0x3f)
	}
}

// B64Decode decodes a base64-encoded integer from the specified byte slice.
func B64Decode(b []byte) (v int, err error) {
	for i := 0; i < len(b); i++ {
		if b[i]&0xc0 != 0x40 {
			err = fmt.Errorf("%w in b64: %q (0x%02[2]x)", ErrInvalidHighBits, b[i])
			return
		}
		v |= int(b[i]&0x3f) << ((len(b) - i - 1) * 6)
	}
	return
}
