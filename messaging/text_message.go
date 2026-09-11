package messaging

import (
	"encoding/json"

	"github.com/resolvingarchitecture/ra-common-go/identity"
)

type TextMessage struct {
	BaseMessage
	To   *identity.Did `json:"to,omitempty"`
	From *identity.Did `json:"from,omitempty"`
	Text *string       `json:"text,omitempty"`
}

func (m *TextMessage) Kind() string { return "text" }

func (m TextMessage) MarshalJSON() ([]byte, error) {
	type alias TextMessage
	return json.Marshal(struct {
		Kind string `json:"kind"`
		alias
	}{Kind: "text", alias: alias(m)})
}
