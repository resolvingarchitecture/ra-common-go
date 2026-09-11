package route

import (
	"encoding/json"
	"testing"
)

func TestRouteSlipLifoAndRoundTrip(t *testing.T) {
	slip := NewDynamicRoutingSlip()
	slip.AddRoute(SimpleRouteOf("a", "op"))
	slip.AddRoute(SimpleRouteOf("b", "op"))
	slip.AddRoute(SimpleRouteOf("c", "op"))
	if slip.NumberRemainingRoutes() != 3 {
		t.Fatalf("NumberRemainingRoutes = %d, want 3", slip.NumberRemainingRoutes())
	}
	if svc := *slip.NextRoute().Meta().Service; svc != "c" {
		t.Errorf("first NextRoute().Service = %q, want c", svc)
	}
	if svc := *slip.NextRoute().Meta().Service; svc != "b" {
		t.Errorf("second NextRoute().Service = %q, want b", svc)
	}
	if svc := *slip.PeekAtNextRoute().Meta().Service; svc != "a" {
		t.Errorf("PeekAtNextRoute().Service = %q, want a", svc)
	}

	data, err := json.Marshal(slip)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	var back DynamicRoutingSlip
	if err := json.Unmarshal(data, &back); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if back.NumberRemainingRoutes() != 1 {
		t.Errorf("round-tripped NumberRemainingRoutes = %d, want 1", back.NumberRemainingRoutes())
	}

	rebuilt, err := UnmarshalRoute(data)
	if err != nil {
		t.Fatalf("UnmarshalRoute: %v", err)
	}
	if _, ok := rebuilt.(*DynamicRoutingSlip); !ok {
		t.Errorf("UnmarshalRoute returned %T, want *DynamicRoutingSlip", rebuilt)
	}
}
