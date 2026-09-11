package crypto

import (
	"encoding/hex"
	"fmt"

	"github.com/resolvingarchitecture/ra-common-go/encoding"
	"github.com/resolvingarchitecture/ra-common-go/raerror"
)

type MultihashType string

const (
	MhSha1     MultihashType = "Sha1"
	MhSha2_256 MultihashType = "Sha2_256"
	MhSha2_512 MultihashType = "Sha2_512"
	MhSha3     MultihashType = "Sha3"
	MhBlake2b  MultihashType = "Blake2b"
	MhBlake2s  MultihashType = "Blake2s"
)

var multihashCodes = map[MultihashType]struct {
	code   byte
	length int
}{
	MhSha1:     {0x11, 20},
	MhSha2_256: {0x12, 32},
	MhSha2_512: {0x13, 64},
	MhSha3:     {0x14, 64},
	MhBlake2b:  {0x40, 64},
	MhBlake2s:  {0x41, 32},
}

func MultihashCode(t MultihashType) byte  { return multihashCodes[t].code }
func MultihashLength(t MultihashType) int { return multihashCodes[t].length }

func MultihashTypeFromCode(code byte) (MultihashType, error) {
	for t, v := range multihashCodes {
		if v.code == code {
			return t, nil
		}
	}
	return "", raerror.InvalidErr(fmt.Sprintf("unknown multihash type: 0x%x", code))
}

// Multihash pairs a hash-type tag with its digest. JSON shape: {"type": ..., "hash": [byte, byte, ...]}.
type Multihash struct {
	Kind   MultihashType `json:"type"`
	Digest []byte        `json:"hash"`
}

func NewMultihash(kind MultihashType, digest []byte) (Multihash, error) {
	if len(digest) > 127 {
		return Multihash{}, raerror.InvalidErr(fmt.Sprintf("unsupported hash size: %d", len(digest)))
	}
	if len(digest) != MultihashLength(kind) {
		return Multihash{}, raerror.InvalidErr(fmt.Sprintf("incorrect hash length: %d != %d", len(digest), MultihashLength(kind)))
	}
	return Multihash{Kind: kind, Digest: digest}, nil
}

func (m Multihash) ToBytes() []byte {
	out := make([]byte, 0, 2+len(m.Digest))
	out = append(out, MultihashCode(m.Kind), byte(len(m.Digest)))
	return append(out, m.Digest...)
}

func MultihashFromBytes(data []byte) (Multihash, error) {
	if len(data) < 2 {
		return Multihash{}, raerror.InvalidErr("multihash too short")
	}
	kind, err := MultihashTypeFromCode(data[0])
	if err != nil {
		return Multihash{}, err
	}
	length := int(data[1])
	if len(data) < 2+length {
		return Multihash{}, raerror.InvalidErr("multihash truncated")
	}
	return NewMultihash(kind, data[2:2+length])
}

func (m Multihash) ToHexString() string { return hex.EncodeToString(m.ToBytes()) }

func MultihashFromHexString(text string) (Multihash, error) {
	data, err := hex.DecodeString(text)
	if err != nil {
		return Multihash{}, raerror.DecodeErr("bad hex: " + text)
	}
	return MultihashFromBytes(data)
}

func (m Multihash) ToBase58() string { return encoding.Base58Encode(m.ToBytes()) }

func MultihashFromBase58(text string) (Multihash, error) {
	data, err := encoding.Base58Decode(text)
	if err != nil {
		return Multihash{}, err
	}
	return MultihashFromBytes(data)
}

func (m Multihash) Equals(other Multihash) bool { return m.ToHexString() == other.ToHexString() }

func (m Multihash) String() string { return m.ToBase58() }
