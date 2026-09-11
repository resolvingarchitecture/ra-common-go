package messaging

import "encoding/json"

type DocumentMessage struct {
	BaseMessage
	Data []map[string]any `json:"data"`
}

func NewDocumentMessage() *DocumentMessage { return &DocumentMessage{Data: []map[string]any{{}}} }

func (m *DocumentMessage) Kind() string { return "document" }

func (m *DocumentMessage) Primary() map[string]any {
	if len(m.Data) == 0 {
		m.Data = append(m.Data, map[string]any{})
	}
	return m.Data[0]
}

func (m *DocumentMessage) Get(key string) any {
	if len(m.Data) == 0 {
		return nil
	}
	return m.Data[0][key]
}

func (m *DocumentMessage) Put(key string, value any) { m.Primary()[key] = value }

func (m DocumentMessage) MarshalJSON() ([]byte, error) {
	type alias DocumentMessage
	return json.Marshal(struct {
		Kind string `json:"kind"`
		alias
	}{Kind: "document", alias: alias(m)})
}
