package util

import (
	"math"
	"testing"
)

func TestBytesStringsVersionNonceUniqueID(t *testing.T) {
	for _, v := range []int32{0, 1, -1, 42, math.MinInt32, math.MaxInt32} {
		if got := PackBigEndian(UnpackBigEndian(v)); got != v {
			t.Errorf("PackBigEndian(UnpackBigEndian(%d)) = %d", v, got)
		}
	}

	if got := Capitalize("one two three"); got != "One Two Three" {
		t.Errorf("Capitalize = %q", got)
	}

	if VersionCompare("1.8", "1.11") != -1 {
		t.Error("VersionCompare(1.8, 1.11) != -1")
	}
	if VersionCompare("2.0", "2.0.0") != -1 {
		t.Error("VersionCompare(2.0, 2.0.0) != -1")
	}
	if VersionCompare("8ea", "8") != 0 {
		t.Error("VersionCompare(8ea, 8) != 0")
	}
	if VersionCompare("8-ea", "8") != 1 {
		t.Error("VersionCompare(8-ea, 8) != 1")
	}

	nonce := DefaultNonce()
	if !nonce.ContinueOn("1") {
		t.Error("first ContinueOn(1) should be true")
	}
	if nonce.ContinueOn("1") {
		t.Error("second ContinueOn(1) should be false")
	}

	uid := RandomUniqueID()
	if len(uid.ToBase64()) != 44 {
		t.Errorf("UniqueID base64 length = %d, want 44", len(uid.ToBase64()))
	}
	back, err := UniqueIDFromBase64(uid.ToBase64())
	if err != nil {
		t.Fatalf("UniqueIDFromBase64: %v", err)
	}
	if back.Compare(uid) != 0 {
		t.Error("UniqueID round trip mismatch")
	}
}
