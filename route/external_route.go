package route

import (
	"encoding/json"

	"github.com/resolvingarchitecture/ra-common-go/network"
)

type SimpleExternalRoute struct {
	MetaValue       RouteMeta     `json:"meta"`
	Origination     *network.Peer `json:"origination,omitempty"`
	Destination     *network.Peer `json:"destination,omitempty"`
	SendContentOnly bool          `json:"send_content_only"`
	StatusCode      int           `json:"status_code"`
}

func NewSimpleExternalRoute() *SimpleExternalRoute {
	return &SimpleExternalRoute{MetaValue: NewRouteMeta(nil, nil)}
}

func SimpleExternalRouteOf(service, operation string) *SimpleExternalRoute {
	r := &SimpleExternalRoute{MetaValue: NewRouteMeta(strPtr(service), strPtr(operation))}
	r.SendContentOnly = true
	return r
}

func (r *SimpleExternalRoute) Type() string     { return "simple_external" }
func (r *SimpleExternalRoute) Meta() *RouteMeta { return &r.MetaValue }

func (r SimpleExternalRoute) MarshalJSON() ([]byte, error) {
	type alias SimpleExternalRoute
	return json.Marshal(struct {
		Type string `json:"type"`
		alias
	}{Type: "simple_external", alias: alias(r)})
}

type RelayedExternalRoute struct {
	Base        SimpleExternalRoute `json:"base"`
	FromPeer    *network.Peer       `json:"from_peer,omitempty"`
	ToPeer      *network.Peer       `json:"to_peer,omitempty"`
	Delayed     bool                `json:"delayed"`
	MinDelay    int                 `json:"min_delay"`
	MaxDelay    int                 `json:"max_delay"`
	Copy        bool                `json:"copy"`
	MinCopies   int                 `json:"min_copies"`
	MaxCopies   int                 `json:"max_copies"`
	Sensitivity int                 `json:"sensitivity"`
}

func NewRelayedExternalRoute() *RelayedExternalRoute {
	return &RelayedExternalRoute{Base: *NewSimpleExternalRoute()}
}

func (r *RelayedExternalRoute) Type() string     { return "relayed_external" }
func (r *RelayedExternalRoute) Meta() *RouteMeta { return &r.Base.MetaValue }

func (r RelayedExternalRoute) MarshalJSON() ([]byte, error) {
	type alias RelayedExternalRoute
	return json.Marshal(struct {
		Type string `json:"type"`
		alias
	}{Type: "relayed_external", alias: alias(r)})
}

const (
	ExternalStatusDestinationPeerRequired     = 2
	ExternalStatusDestinationPeerWrongNetwork = 3
	ExternalStatusDestinationPeerNotFound     = 4
	ExternalStatusNoService                   = 7
	ExternalStatusNoOperation                 = 8
	ExternalStatusNoAddress                   = 9
	ExternalStatusNoFingerprint               = 10
	ExternalStatusNoPort                      = 11
)
