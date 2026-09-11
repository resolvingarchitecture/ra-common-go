package route

import (
	"encoding/json"
	"fmt"
)

// DynamicRoutingSlip is a LIFO stack of routes walked one hop at a time.
// AddRoute/NextRoute/PeekAtNextRoute all operate on the end of the slice -
// append-to-end + pop-from-end is a stack, so there's no need for the
// insert-at-front ("unshift") the JS-targeting ports use; same LIFO
// behaviour, cheaper (O(1) at a Go slice's natural end instead of O(n) at its front).
type DynamicRoutingSlip struct {
	MetaValue RouteMeta `json:"meta"`
	routes    []Route
	current   Route
}

func NewDynamicRoutingSlip() *DynamicRoutingSlip {
	return &DynamicRoutingSlip{MetaValue: NewRouteMeta(nil, nil)}
}

func (s *DynamicRoutingSlip) Type() string     { return "routing_slip" }
func (s *DynamicRoutingSlip) Meta() *RouteMeta { return &s.MetaValue }

// AddRoute pushes route onto the stack, stamping it with this slip's route id.
func (s *DynamicRoutingSlip) AddRoute(r Route) {
	r.Meta().RouteID = s.MetaValue.RouteID
	s.routes = append(s.routes, r)
}

func (s *DynamicRoutingSlip) NumberRemainingRoutes() int { return len(s.routes) }

func (s *DynamicRoutingSlip) CurrentRoute() Route {
	if s.current == nil {
		s.NextRoute()
	}
	return s.current
}

func (s *DynamicRoutingSlip) NextRoute() Route {
	if len(s.routes) == 0 {
		s.current = nil
		return nil
	}
	last := len(s.routes) - 1
	s.current = s.routes[last]
	s.routes = s.routes[:last]
	return s.current
}

func (s *DynamicRoutingSlip) PeekAtNextRoute() Route {
	if len(s.routes) == 0 {
		return nil
	}
	return s.routes[len(s.routes)-1]
}

func (s DynamicRoutingSlip) MarshalJSON() ([]byte, error) {
	routes := make([]json.RawMessage, len(s.routes))
	for i, r := range s.routes {
		data, err := json.Marshal(r)
		if err != nil {
			return nil, err
		}
		routes[i] = data
	}
	aux := struct {
		Type    string            `json:"type"`
		Meta    RouteMeta         `json:"meta"`
		Routes  []json.RawMessage `json:"routes"`
		Current json.RawMessage   `json:"current_route,omitempty"`
	}{Type: "routing_slip", Meta: s.MetaValue, Routes: routes}

	if s.current != nil {
		data, err := json.Marshal(s.current)
		if err != nil {
			return nil, err
		}
		aux.Current = data
	}
	return json.Marshal(aux)
}

func (s *DynamicRoutingSlip) UnmarshalJSON(data []byte) error {
	var aux struct {
		Meta    RouteMeta         `json:"meta"`
		Routes  []json.RawMessage `json:"routes"`
		Current json.RawMessage   `json:"current_route"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	s.MetaValue = aux.Meta
	s.routes = make([]Route, 0, len(aux.Routes))
	for _, raw := range aux.Routes {
		r, err := UnmarshalRoute(raw)
		if err != nil {
			return err
		}
		s.routes = append(s.routes, r)
	}
	if len(aux.Current) > 0 {
		r, err := UnmarshalRoute(aux.Current)
		if err != nil {
			return err
		}
		s.current = r
	}
	return nil
}

// UnmarshalRoute rebuilds a polymorphic Route from its `type` tag.
func UnmarshalRoute(data []byte) (Route, error) {
	var probe struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return nil, err
	}
	switch probe.Type {
	case "simple":
		var r SimpleRoute
		if err := json.Unmarshal(data, &r); err != nil {
			return nil, err
		}
		return &r, nil
	case "routing_slip":
		var r DynamicRoutingSlip
		if err := json.Unmarshal(data, &r); err != nil {
			return nil, err
		}
		return &r, nil
	case "simple_external":
		var r SimpleExternalRoute
		if err := json.Unmarshal(data, &r); err != nil {
			return nil, err
		}
		return &r, nil
	case "relayed_external":
		var r RelayedExternalRoute
		if err := json.Unmarshal(data, &r); err != nil {
			return nil, err
		}
		return &r, nil
	default:
		return nil, fmt.Errorf("unknown route type: %s", probe.Type)
	}
}
