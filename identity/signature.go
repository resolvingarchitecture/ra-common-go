// Package identity: Did, PublicKey, Signature. Ports the ra.common.identity package.
package identity

// Signature is a detached signature over some value. This package does not
// sign or verify; Signature is a metadata record carried inside PublicKey.
// Equality is on SignedByAddress.
type Signature struct {
	ValueSigned         *string `json:"value_signed,omitempty"`
	Algorithm           *string `json:"algorithm,omitempty"`
	SignedDate          *string `json:"signed_date,omitempty"` // RFC 3339
	SignedByUsername    *string `json:"signed_by_username,omitempty"`
	SignedByFingerprint *string `json:"signed_by_fingerprint,omitempty"`
	SignedByAddress     *string `json:"signed_by_address,omitempty"`
}

func (s Signature) Equals(other Signature) bool {
	return s.SignedByAddress != nil && other.SignedByAddress != nil && *s.SignedByAddress == *other.SignedByAddress
}
