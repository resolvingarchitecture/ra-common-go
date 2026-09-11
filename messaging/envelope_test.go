package messaging

import "testing"

func TestFactoriesSetMessageKind(t *testing.T) {
	if _, ok := DocumentEnvelope().MessageValue.(*DocumentMessage); !ok {
		t.Error("DocumentEnvelope() should carry a *DocumentMessage")
	}
	if _, ok := CommandEnvelope().MessageValue.(*CommandMessage); !ok {
		t.Error("CommandEnvelope() should carry a *CommandMessage")
	}
	if _, ok := EventEnvelope(EventTypeBusStatus).MessageValue.(*EventMessage); !ok {
		t.Error("EventEnvelope() should carry an *EventMessage")
	}
	if HeadersOnlyEnvelope().MessageValue != nil {
		t.Error("HeadersOnlyEnvelope() should carry no message")
	}
}

func TestContentRoundTripsThroughADocument(t *testing.T) {
	e := DocumentEnvelope()
	if !e.AddContent("hello") {
		t.Fatal("AddContent should succeed on a document envelope")
	}
	if e.Content() != "hello" {
		t.Errorf("Content() = %v, want hello", e.Content())
	}
	if CommandEnvelope().AddContent(nil) {
		t.Error("AddContent should fail on a command envelope")
	}
}

func TestExceptionsAccumulate(t *testing.T) {
	e := DocumentEnvelope()
	e.AddException("first")
	e.AddException("second")
	got := e.Exceptions()
	if len(got) != 2 || got[0] != "first" || got[1] != "second" {
		t.Errorf("Exceptions() = %v, want [first second]", got)
	}
}

func TestRatchetWalksTheSlipLifo(t *testing.T) {
	e := DocumentEnvelope()
	e.AddRoute("ra.a.ServiceA", "OP")
	e.AddRoute("ra.b.ServiceB", "OP")
	e.Ratchet()
	if svc := *e.GetRoute().Meta().Service; svc != "ra.b.ServiceB" {
		t.Errorf("first hop = %q, want ra.b.ServiceB", svc)
	}
	e.Ratchet()
	if svc := *e.GetRoute().Meta().Service; svc != "ra.a.ServiceA" {
		t.Errorf("second hop = %q, want ra.a.ServiceA", svc)
	}
}

func TestJSONRoundTrip(t *testing.T) {
	e := DocumentEnvelope()
	e.SetContentType(HeaderContentTypeJSON)
	e.AddContent(float64(42)) // JSON numbers round-trip as float64
	e.AddRoute("ra.x.Svc", "DO")
	e.Mark("seen")

	text, err := e.ToJSONString()
	if err != nil {
		t.Fatalf("ToJSONString: %v", err)
	}
	back, err := EnvelopeFromJSONString(text)
	if err != nil {
		t.Fatalf("EnvelopeFromJSONString: %v", err)
	}
	if back.ID != e.ID {
		t.Errorf("ID = %q, want %q", back.ID, e.ID)
	}
	if back.ContentType() == nil || *back.ContentType() != HeaderContentTypeJSON {
		t.Error("ContentType round trip failed")
	}
	if back.Content() != float64(42) {
		t.Errorf("Content() = %v, want 42", back.Content())
	}
	if !back.MarkerPresent("seen") {
		t.Error("marker 'seen' should be present after round trip")
	}
	if back.DynamicRoutingSlip.NumberRemainingRoutes() != 1 {
		t.Errorf("NumberRemainingRoutes = %d, want 1", back.DynamicRoutingSlip.NumberRemainingRoutes())
	}
}

func TestEqualityIsByID(t *testing.T) {
	if !DocumentEnvelopeWithID("same").Equals(DocumentEnvelopeWithID("same")) {
		t.Error("envelopes with the same id should be equal")
	}
}
