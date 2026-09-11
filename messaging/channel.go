package messaging

import "github.com/resolvingarchitecture/ra-common-go/lifecycle"

// Client is a caller that a service can send a reply Envelope back to.
type Client interface {
	Reply(envelope *Envelope)
}

type MessageProducer interface {
	Send(envelope *Envelope) bool
	SendWithCallback(envelope *Envelope, callback Client) bool
	DeadLetter(envelope *Envelope) bool
}

type MessageConsumer interface {
	Receive(envelope *Envelope) bool
}

type MessageChannel interface {
	MessageProducer
	lifecycle.LifeCycle
	Name() string
	IsPubSub() bool
	Queued() int
	Ack(envelope *Envelope)
}

type MessageBus interface {
	lifecycle.LifeCycle
	RegisterChannel(name string, serviceLevel *string) bool
	Publish(envelope *Envelope) bool
	PublishWithCallback(envelope *Envelope, callback Client) bool
	Completed(envelope *Envelope) bool
}
