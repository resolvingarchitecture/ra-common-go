// Package messaging holds Message (and its concrete kinds), Envelope, and
// the producer/consumer/channel/bus interfaces together. Ports the
// ra.common.messaging package plus ra.common.Envelope. They live in one Go
// package (not split like the other ports' files) because Envelope embeds
// Message and the bus/channel interfaces take an *Envelope parameter - a
// real circular dependency every other port also has, but Go's build
// (unlike TS/C#/C++/Python) flatly rejects circular package imports, even
// for interface-only references. See DESIGN.md.
package messaging

import (
	"encoding/json"
	"fmt"
)

const (
	Content    = "CONTENT"
	Entity     = "ENTITY"
	Exceptions = "EXCEPTIONS"
)

// Message is the interface every message kind implements.
type Message interface {
	Kind() string
	AddErrorMessage(msg string)
	ClearErrorMessages()
	ErrorMessages() []string
}

// BaseMessage supplies the error-message bookkeeping shared by every kind,
// via embedding (Go's composition-over-inheritance stand-in for the Java
// BaseMessage abstract class).
type BaseMessage struct {
	ErrorMessagesValue []string `json:"error_messages,omitempty"`
}

func (m *BaseMessage) AddErrorMessage(msg string) {
	m.ErrorMessagesValue = append(m.ErrorMessagesValue, msg)
}
func (m *BaseMessage) ClearErrorMessages()     { m.ErrorMessagesValue = nil }
func (m *BaseMessage) ErrorMessages() []string { return m.ErrorMessagesValue }

// UnmarshalMessage rebuilds a polymorphic Message from its `kind` tag.
func UnmarshalMessage(data []byte) (Message, error) {
	var probe struct {
		Kind string `json:"kind"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return nil, err
	}
	switch probe.Kind {
	case "document":
		var m DocumentMessage
		if err := json.Unmarshal(data, &m); err != nil {
			return nil, err
		}
		return &m, nil
	case "command":
		var m CommandMessage
		if err := json.Unmarshal(data, &m); err != nil {
			return nil, err
		}
		return &m, nil
	case "event":
		var m EventMessage
		if err := json.Unmarshal(data, &m); err != nil {
			return nil, err
		}
		return &m, nil
	case "text":
		var m TextMessage
		if err := json.Unmarshal(data, &m); err != nil {
			return nil, err
		}
		return &m, nil
	default:
		return nil, fmt.Errorf("unknown message kind: %s", probe.Kind)
	}
}
