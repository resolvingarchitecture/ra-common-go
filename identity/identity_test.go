package identity

import (
	"encoding/json"
	"testing"

	"github.com/resolvingarchitecture/ra-common-go/crypto"
)

func TestDidSignaturePublicKey(t *testing.T) {
	d := WithUsername("alice")
	secret := "secret"
	d.Passphrase = &secret
	d.Authenticated = true
	d.ClearSensitive()
	if d.Username != "" {
		t.Errorf("Username after ClearSensitive = %q, want empty", d.Username)
	}
	if d.Status != DidPrivate {
		t.Errorf("Status after ClearSensitive = %q, want Private", d.Status)
	}

	bob := WithUsername("bob")
	bob.PublicKey = NewPublicKeyFromAddress("addr")
	bob.PassphraseHash = &crypto.Hash{HashValue: "deadbeef", Algorithm: crypto.Sha256}

	data, err := json.Marshal(bob)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	// Unmarshal into a NewDid() so any missing/omitted fields keep their
	// documented defaults, matching the other ports' FromJson behaviour -
	// see DESIGN.md.
	back := NewDid()
	if err := json.Unmarshal(data, &back); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if back.Username != "bob" {
		t.Errorf("Username = %q, want bob", back.Username)
	}
	if back.PublicKey.AddressValue == nil || *back.PublicKey.AddressValue != "addr" {
		t.Error("PublicKey.AddressValue round trip failed")
	}
	if back.PassphraseHash == nil || back.PassphraseHash.HashValue != "deadbeef" {
		t.Error("PassphraseHash round trip failed")
	}

	var a, b Signature
	if a.Equals(b) {
		t.Error("two empty signatures should not be equal")
	}
	addr := "x"
	a.SignedByAddress = &addr
	b.SignedByAddress = &addr
	if !a.Equals(b) {
		t.Error("signatures with the same SignedByAddress should be equal")
	}

	pk := NewPublicKeyFromAddress("B32")
	signer := "signer"
	pk.AddSignedAttribute("email", Signature{SignedByAddress: &signer})
	if len(pk.SignedAttributes["email"]) != 1 {
		t.Fatalf("expected 1 signed attribute, got %d", len(pk.SignedAttributes["email"]))
	}
	pk.RemoveSignature("email", "signer")
	if len(pk.SignedAttributes["email"]) != 0 {
		t.Errorf("expected 0 signed attributes after removal, got %d", len(pk.SignedAttributes["email"]))
	}
}
