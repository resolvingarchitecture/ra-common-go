package messaging

import (
	"encoding/json"
	"testing"
)

func TestMessagingTaggedRoundTrip(t *testing.T) {
	start := CmdStart
	shutdown := CmdGracefullyShutdown
	msgs := []Message{
		NewDocumentMessage(),
		NewCommandMessage(&start),
		NewCommandMessage(&shutdown),
	}
	for _, msg := range msgs {
		data, err := json.Marshal(msg)
		if err != nil {
			t.Fatalf("Marshal: %v", err)
		}
		var probe struct {
			Kind string `json:"kind"`
		}
		if err := json.Unmarshal(data, &probe); err != nil || probe.Kind == "" {
			t.Fatalf("expected a kind field, got %v (%v)", probe, err)
		}
		back, err := UnmarshalMessage(data)
		if err != nil {
			t.Fatalf("UnmarshalMessage: %v", err)
		}
		if got, want := formatType(back), formatType(msg); got != want {
			t.Errorf("UnmarshalMessage type = %s, want %s", got, want)
		}
	}

	data, _ := json.Marshal(NewCommandMessage(&shutdown))
	var decoded map[string]any
	_ = json.Unmarshal(data, &decoded)
	if decoded["command"] != "GracefullyShutdown" {
		t.Errorf(`command field = %v, want "GracefullyShutdown"`, decoded["command"])
	}
}

func formatType(m Message) string {
	switch m.(type) {
	case *DocumentMessage:
		return "document"
	case *CommandMessage:
		return "command"
	case *EventMessage:
		return "event"
	case *TextMessage:
		return "text"
	default:
		return "unknown"
	}
}
