// Package encoding provides Base32 and Base58 string codecs. Ports
// ra.common.Base32 and ra.common.Base58. Since this port is not
// wire-compatible we use standard implementations: base32 is RFC 4648,
// uppercase A-Z2-7, no padding; base58 is the Bitcoin alphabet. Base64 is not
// here - use Go's stdlib encoding/base64 directly, unlike every other port
// (except ra-common-cpp) which had to hand-roll it.
package encoding

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/resolvingarchitecture/ra-common-go/raerror"
)

const base32Alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
const base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

func Base32Encode(data []byte) string {
	bits := 0
	value := 0
	var out strings.Builder
	for _, b := range data {
		value = (value << 8) | int(b)
		bits += 8
		for bits >= 5 {
			out.WriteByte(base32Alphabet[(value>>(bits-5))&0x1f])
			bits -= 5
		}
	}
	if bits > 0 {
		out.WriteByte(base32Alphabet[(value<<(5-bits))&0x1f])
	}
	return out.String()
}

func Base32Decode(text string) ([]byte, error) {
	bits := 0
	value := 0
	out := make([]byte, 0, len(text))
	for _, ch := range text {
		idx := strings.IndexRune(base32Alphabet, ch)
		if idx < 0 {
			return nil, raerror.DecodeErr(fmt.Sprintf("invalid base32 character: %c", ch))
		}
		value = (value << 5) | idx
		bits += 5
		if bits >= 8 {
			out = append(out, byte((value>>(bits-8))&0xff))
			bits -= 8
		}
	}
	return out, nil
}

func Base58Encode(data []byte) string {
	n := new(big.Int).SetBytes(data)
	mod := big.NewInt(58)
	var out []byte
	zero := big.NewInt(0)
	rem := new(big.Int)
	for n.Cmp(zero) > 0 {
		n.DivMod(n, mod, rem)
		out = append([]byte{base58Alphabet[rem.Int64()]}, out...)
	}
	pad := 0
	for _, b := range data {
		if b == 0 {
			pad++
		} else {
			break
		}
	}
	return strings.Repeat(string(base58Alphabet[0]), pad) + string(out)
}

func Base58Decode(text string) ([]byte, error) {
	n := new(big.Int)
	mod := big.NewInt(58)
	for _, ch := range text {
		idx := strings.IndexRune(base58Alphabet, ch)
		if idx < 0 {
			return nil, raerror.DecodeErr(fmt.Sprintf("invalid base58 character: %c", ch))
		}
		n.Mul(n, mod)
		n.Add(n, big.NewInt(int64(idx)))
	}
	decoded := n.Bytes()

	pad := 0
	for _, ch := range text {
		if byte(ch) == base58Alphabet[0] {
			pad++
		} else {
			break
		}
	}
	out := make([]byte, pad+len(decoded))
	copy(out[pad:], decoded)
	return out, nil
}
