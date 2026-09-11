package crypto

// Hash is a hash string with the algorithm used. Equality is on the hash
// string alone. Plain tagged struct - encoding/json handles (de)serialization
// directly, unlike the other ports which each need an explicit
// ToJson()/FromJson() pair (see DESIGN.md).
type Hash struct {
	HashValue string        `json:"hash"`
	Algorithm HashAlgorithm `json:"algorithm"`
}

func NewHash(hashValue string, algorithm HashAlgorithm) Hash {
	return Hash{HashValue: hashValue, Algorithm: algorithm}
}

func (h Hash) Equals(other Hash) bool { return h.HashValue == other.HashValue }

func (h Hash) String() string { return h.HashValue }
