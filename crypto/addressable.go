package crypto

// Addressable is implemented by anything with a fingerprint/address pair (e.g. PublicKey).
type Addressable interface {
	Fingerprint() *string
	Address() *string
}
