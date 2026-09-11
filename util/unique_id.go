package util

import (
	"bytes"
	"encoding/base64"

	"github.com/resolvingarchitecture/ra-common-go/raerror"
)

const uniqueIDLength = 32

// UniqueID is a fixed 32-byte identifier, rendered as standard padded base64 (44 chars).
type UniqueID struct {
	bytes []byte
}

func NewUniqueID(b []byte) (*UniqueID, error) {
	if len(b) != uniqueIDLength {
		return nil, raerror.InvalidErr("UniqueId must be 32 bytes")
	}
	return &UniqueID{bytes: b}, nil
}

func RandomUniqueID() *UniqueID {
	id, _ := NewUniqueID(RandomBytesOf(uniqueIDLength)) // length is exactly right by construction
	return id
}

func UniqueIDFromBase64(text string) (*UniqueID, error) {
	raw, err := base64.StdEncoding.DecodeString(text)
	if err != nil || len(raw) != uniqueIDLength {
		return nil, raerror.DecodeErr("UniqueId must be 32 bytes")
	}
	return &UniqueID{bytes: raw}, nil
}

func (u *UniqueID) Bytes() []byte { return u.bytes }

func (u *UniqueID) ToBase64() string { return base64.StdEncoding.EncodeToString(u.bytes) }

func (u *UniqueID) String() string { return u.ToBase64() }

func (u *UniqueID) Compare(other *UniqueID) int { return bytes.Compare(u.bytes, other.bytes) }
