// Package route: routing slips and external/relayed routes. Ports the
// ra.common.route package. The Java abstract base + reflective
// Class.forName polymorphism becomes a Route interface + a `type` tag on the
// wire (simple, routing_slip, simple_external, relayed_external), rebuilt by
// UnmarshalRoute.
package route

import "github.com/resolvingarchitecture/ra-common-go/util"

type RouteMeta struct {
	Service   *string `json:"service,omitempty"`
	Operation *string `json:"operation,omitempty"`
	Routed    bool    `json:"routed"`
	// RouteID is a plain JSON number here (unlike the string-encoded route_id
	// in the JS-targeting ports): Go's encoding/json round-trips int64
	// exactly, so it doesn't need the precision workaround JS's float64-based
	// JSON numbers require.
	RouteID int64 `json:"route_id"`
}

func NewRouteMeta(service, operation *string) RouteMeta {
	return RouteMeta{Service: service, Operation: operation, RouteID: util.NextLong()}
}
