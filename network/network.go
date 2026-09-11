// Package network holds the minimal network identity types needed by route
// and messaging. Ports ra.common.network.Network, NetworkStatus and
// NetworkPeer. The full network service layer is deferred to a later phase.
package network

import "github.com/resolvingarchitecture/ra-common-go/identity"

type Network string

const (
	Card      Network = "Card"
	Nfc       Network = "Nfc"
	Http      Network = "Http"
	Tor       Network = "Tor"
	I2p       Network = "I2p"
	WiFi      Network = "WiFi"
	Bluetooth Network = "Bluetooth"
	Satellite Network = "Satellite"
	FsRadio   Network = "FsRadio"
	LiFi      Network = "LiFi"
)

type Status string

const (
	NotInstalled Status = "NotInstalled"
	Closed       Status = "Closed"
	Error        Status = "Error"
	PortConflict Status = "PortConflict"
	Waiting      Status = "Waiting"
	Warmup       Status = "Warmup"
	Connecting   Status = "Connecting"
	Connected    Status = "Connected"
	Verified     Status = "Verified"
	Hanging      Status = "Hanging"
	Failed       Status = "Failed"
	Blocked      Status = "Blocked"
	Disconnected Status = "Disconnected"
)

// Peer is a peer in a peer-to-peer network. Equality follows the Java
// version: two peers are equal iff both have a public-key address and
// fingerprint and both match.
type Peer struct {
	NetworkValue Network      `json:"network"`
	Did          identity.Did `json:"did"`
	ID           *string      `json:"id,omitempty"`
	Port         *int         `json:"port,omitempty"`
	Services     []string     `json:"services,omitempty"`
}

func NewPeer(network Network) Peer {
	return Peer{NetworkValue: network, Did: identity.NewDid()}
}

func WithCredentials(network Network, username string, passphrase *string) Peer {
	p := NewPeer(network)
	p.Did = identity.WithUsername(username)
	p.Did.Passphrase = passphrase
	return p
}

func (p Peer) key() *string {
	if p.Did.PublicKey.AddressValue != nil && p.Did.PublicKey.FingerprintValue != nil {
		k := *p.Did.PublicKey.AddressValue + " " + *p.Did.PublicKey.FingerprintValue
		return &k
	}
	return nil
}

func (p Peer) Equals(other Peer) bool {
	a := p.key()
	b := other.key()
	return a != nil && b != nil && *a == *b
}
