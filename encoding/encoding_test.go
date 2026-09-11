package encoding

import (
	"bytes"
	"testing"
)

func TestEncodingRoundTrips(t *testing.T) {
	data := []byte("resolving architecture")
	encoded := Base32Encode(data)
	decoded, err := Base32Decode(encoded)
	if err != nil {
		t.Fatalf("Base32Decode: %v", err)
	}
	if !bytes.Equal(decoded, data) {
		t.Errorf("base32 round trip mismatch: got %v want %v", decoded, data)
	}

	hw := []byte("Hello World!")
	if got := Base58Encode(hw); got != "2NEpo7TZRRrLZSi2U" {
		t.Errorf("Base58Encode(%q) = %q, want 2NEpo7TZRRrLZSi2U", hw, got)
	}
	back, err := Base58Decode("2NEpo7TZRRrLZSi2U")
	if err != nil {
		t.Fatalf("Base58Decode: %v", err)
	}
	if !bytes.Equal(back, hw) {
		t.Errorf("base58 decode mismatch: got %v want %v", back, hw)
	}
}
