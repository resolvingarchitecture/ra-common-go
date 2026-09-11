package identity

import "github.com/resolvingarchitecture/ra-common-go/crypto"

type DidStatus string

const (
	DidInactive  DidStatus = "Inactive"
	DidActive    DidStatus = "Active"
	DidSuspended DidStatus = "Suspended"
	DidPrivate   DidStatus = "Private"
)

type DidType string

const (
	DidTypeContact  DidType = "Contact"
	DidTypeIdentity DidType = "Identity"
	DidTypeNode     DidType = "Node"
)

// Did is a decentralized identity: a username, an optional passphrase (+ its
// hash), and a PublicKey. Deliberately does not follow the W3C DID spec - RA
// models each key as its own identity rather than grouping keys.
type Did struct {
	Username                string               `json:"username"`
	Passphrase              *string              `json:"passphrase,omitempty"`
	Passphrase2             *string              `json:"passphrase2,omitempty"`
	PassphraseHash          *crypto.Hash         `json:"passphrase_hash,omitempty"`
	PassphraseHashAlgorithm crypto.HashAlgorithm `json:"passphrase_hash_algorithm"`
	Description             string               `json:"description"`
	Status                  DidStatus            `json:"status"`
	DidType                 DidType              `json:"did_type"`
	Verified                bool                 `json:"verified"`
	Authenticated           bool                 `json:"authenticated"`
	PublicKey               PublicKey            `json:"public_key"`
}

func NewDid() Did {
	return Did{
		Username:                "Anon",
		PassphraseHashAlgorithm: crypto.Pbkdf2HmacSha1,
		Status:                  DidInactive,
		DidType:                 DidTypeIdentity,
	}
}

func WithUsername(username string) Did {
	d := NewDid()
	d.Username = username
	return d
}

func (d *Did) EffectivePassphraseHashAlgorithm() crypto.HashAlgorithm {
	if d.PassphraseHash != nil {
		return d.PassphraseHash.Algorithm
	}
	return d.PassphraseHashAlgorithm
}

func (d *Did) ClearSensitive() {
	d.Username = ""
	d.Passphrase = nil
	d.Passphrase2 = nil
	d.Description = ""
	d.Status = DidPrivate
	d.Verified = false
	d.Authenticated = false
}
