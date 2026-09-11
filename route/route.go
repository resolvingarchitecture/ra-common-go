package route

import "encoding/json"

// Route is any route variant.
type Route interface {
	Type() string
	Meta() *RouteMeta
}

func strPtr(s string) *string { return &s }

// SimpleRoute is the plain service+operation hop.
type SimpleRoute struct {
	MetaValue RouteMeta `json:"meta"`
}

func SimpleRouteOf(service, operation string) *SimpleRoute {
	return &SimpleRoute{MetaValue: NewRouteMeta(strPtr(service), strPtr(operation))}
}

func (r *SimpleRoute) Type() string     { return "simple" }
func (r *SimpleRoute) Meta() *RouteMeta { return &r.MetaValue }

func (r SimpleRoute) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string    `json:"type"`
		Meta RouteMeta `json:"meta"`
	}{Type: "simple", Meta: r.MetaValue})
}
