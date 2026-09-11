package network

import (
	"encoding/json"
	"testing"
)

func TestPeerEquality(t *testing.T) {
	a := NewPeer(Tor)
	b := NewPeer(Tor)
	if a.Equals(b) {
		t.Error("two fresh peers with no address/fingerprint should not be equal")
	}

	addr, fp := "addr", "fp"
	a.Did.PublicKey.AddressValue = &addr
	a.Did.PublicKey.FingerprintValue = &fp
	b.Did.PublicKey.AddressValue = &addr
	b.Did.PublicKey.FingerprintValue = &fp
	if !a.Equals(b) {
		t.Error("peers with matching address+fingerprint should be equal")
	}

	data, err := json.Marshal(a)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	var back Peer
	if err := json.Unmarshal(data, &back); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if back.NetworkValue != Tor {
		t.Errorf("NetworkValue = %q, want Tor", back.NetworkValue)
	}
}
