package messaging

import (
	"encoding/json"

	"github.com/resolvingarchitecture/ra-common-go/util"
)

const (
	EventTypeError              = "ERROR"
	EventTypeException          = "EXCEPTION"
	EventTypeBusStatus          = "BUS_STATUS"
	EventTypePeerStatus         = "PEER_STATUS"
	EventTypeServiceStatus      = "SERVICE_STATUS"
	EventTypeDidStatus          = "DID_STATUS"
	EventTypeNetworkStateUpdate = "NETWORK_STATE_UPDATE"
	EventTypePriceChange        = "PRICE_CHANGE"
)

type EventMessage struct {
	BaseMessage
	ID           string  `json:"id"`
	EventType    string  `json:"event_type"`
	Name         *string `json:"name,omitempty"`
	MessageValue any     `json:"message,omitempty"`
}

func NewEventMessage(eventType string) *EventMessage {
	return &EventMessage{ID: util.RandomAlphanumeric(32), EventType: eventType}
}

func (m *EventMessage) Kind() string { return "event" }

func (m EventMessage) MarshalJSON() ([]byte, error) {
	type alias EventMessage
	return json.Marshal(struct {
		Kind string `json:"kind"`
		alias
	}{Kind: "event", alias: alias(m)})
}
