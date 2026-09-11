package identity

// PublicKey. IsIdentityKey/IsEncryptionKey/etc. use `omitempty` so a false
// value (the common case) is dropped from the wire, matching the other
// ports' "only emit true flags" behaviour.
type PublicKey struct {
	Alias            *string                `json:"alias,omitempty"`
	FingerprintValue *string                `json:"fingerprint,omitempty"`
	AddressValue     *string                `json:"address,omitempty"`
	KeyType          *string                `json:"type,omitempty"`
	IsIdentityKey    bool                   `json:"is_identity_key,omitempty"`
	IsEncryptionKey  bool                   `json:"is_encryption_key,omitempty"`
	IsBase64Encoded  bool                   `json:"is_base64_encoded,omitempty"`
	IsBase58Encoded  bool                   `json:"is_base58_encoded,omitempty"`
	IsPem            bool                   `json:"is_pem,omitempty"`
	IsHex            bool                   `json:"is_hex,omitempty"`
	Attributes       map[string]any         `json:"attributes,omitempty"`
	SignedAttributes map[string][]Signature `json:"signed_attributes,omitempty"`
}

// Fingerprint and Address implement crypto.Addressable.
func (pk *PublicKey) Fingerprint() *string { return pk.FingerprintValue }
func (pk *PublicKey) Address() *string     { return pk.AddressValue }

func NewPublicKeyFromAddress(address string) PublicKey {
	return PublicKey{AddressValue: &address}
}

func (pk *PublicKey) AddAttribute(name string, value any) {
	if pk.Attributes == nil {
		pk.Attributes = map[string]any{}
	}
	pk.Attributes[name] = value
}

func (pk *PublicKey) Attribute(name string) any { return pk.Attributes[name] }

func (pk *PublicKey) AddSignedAttribute(name string, signature Signature) {
	if pk.SignedAttributes == nil {
		pk.SignedAttributes = map[string][]Signature{}
	}
	pk.SignedAttributes[name] = append(pk.SignedAttributes[name], signature)
}

func (pk *PublicKey) RemoveSignature(name, signedByAddress string) {
	sigs, ok := pk.SignedAttributes[name]
	if !ok {
		return
	}
	filtered := sigs[:0]
	for _, s := range sigs {
		if s.SignedByAddress == nil || *s.SignedByAddress != signedByAddress {
			filtered = append(filtered, s)
		}
	}
	pk.SignedAttributes[name] = filtered
}
