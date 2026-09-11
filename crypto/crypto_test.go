package crypto

import (
	"strings"
	"testing"
)

func TestHashPasswordHashcashMultihash(t *testing.T) {
	if !NewHash("abc", Sha256).Equals(NewHash("abc", Sha1)) {
		t.Error("Hash equality should ignore algorithm")
	}

	alice := []byte("Alice")
	h, err := GenerateHash(alice, Sha256)
	if err != nil {
		t.Fatalf("GenerateHash: %v", err)
	}
	ok, err := VerifyHash(alice, h, Sha256)
	if err != nil || !ok {
		t.Errorf("VerifyHash(alice) = %v, %v", ok, err)
	}
	ok, _ = VerifyHash([]byte("Bob"), h, Sha256)
	if ok {
		t.Error("VerifyHash(bob) should be false")
	}

	pw := GeneratePasswordHash("hunter2")
	if !strings.HasPrefix(pw, "1000_") {
		t.Errorf("password hash should start with 1000_, got %q", pw)
	}
	if !VerifyPasswordHash("hunter2", pw) {
		t.Error("VerifyPasswordHash(hunter2) should be true")
	}
	if VerifyPasswordHash("hunter3", pw) {
		t.Error("VerifyPasswordHash(hunter3) should be false")
	}

	digest := make([]byte, 32)
	for i := range digest {
		digest[i] = 0xab
	}
	m, err := NewMultihash(MhSha2_256, digest)
	if err != nil {
		t.Fatalf("NewMultihash: %v", err)
	}
	if back, err := MultihashFromBytes(m.ToBytes()); err != nil || !back.Equals(m) {
		t.Errorf("MultihashFromBytes round trip failed: %v", err)
	}
	if back, err := MultihashFromHexString(m.ToHexString()); err != nil || !back.Equals(m) {
		t.Errorf("MultihashFromHexString round trip failed: %v", err)
	}
	if back, err := MultihashFromBase58(m.ToBase58()); err != nil || !back.Equals(m) {
		t.Errorf("MultihashFromBase58 round trip failed: %v", err)
	}

	hc, err := MintHashCash("brian@resolvingarchitecture.io", 10)
	if err != nil {
		t.Fatalf("MintHashCash: %v", err)
	}
	if hc.ComputedBits() < 10 {
		t.Error("ComputedBits should be >= 10")
	}
	if !hc.IsValidFor("brian@resolvingarchitecture.io", 10) {
		t.Error("IsValidFor(correct resource) should be true")
	}
	if hc.IsValidFor("nope", 10) {
		t.Error("IsValidFor(wrong resource) should be false")
	}
	parsed, err := ParseHashCash(hc.Token())
	if err != nil || parsed.Resource() != hc.Resource() {
		t.Errorf("ParseHashCash round trip failed: %v", err)
	}

	if LeadingZeroBits([]byte{0, 0, 0x0f}) != 20 {
		t.Error("LeadingZeroBits({0,0,0x0f}) != 20")
	}
	if LeadingZeroBits([]byte{0xff}) != 0 {
		t.Error("LeadingZeroBits({0xff}) != 0")
	}
}
