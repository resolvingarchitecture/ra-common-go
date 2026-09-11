package messaging

import (
	"encoding/json"

	"github.com/resolvingarchitecture/ra-common-go/identity"
	"github.com/resolvingarchitecture/ra-common-go/multipart"
	"github.com/resolvingarchitecture/ra-common-go/route"
	"github.com/resolvingarchitecture/ra-common-go/servicestatus"
	"github.com/resolvingarchitecture/ra-common-go/util"
)

const (
	HeaderAuthorization           = "Authorization"
	HeaderContentDisposition      = "Content-Disposition"
	HeaderContentTransferEncoding = "Content-Transfer-Encoding"
	HeaderContentType             = "Content-Type"
	HeaderContentTypeJSON         = "application/json"
	HeaderUserAgent               = "User-Agent"
)

type Action string

const (
	ActionPost   Action = "Post"
	ActionPut    Action = "Put"
	ActionDelete Action = "Delete"
	ActionGet    Action = "Get"
)

// Envelope wraps everything passed around the application so there is always
// a place for header/routing metadata. Ports ra.common.Envelope (and folds
// in the useful parts of the deprecated ra.common.DLC static helpers as
// methods). Equality is by ID.
type Envelope struct {
	ID                 string                     `json:"id"`
	DynamicRoutingSlip route.DynamicRoutingSlip   `json:"dynamic_routing_slip"`
	RouteValue         route.Route                `json:"route,omitempty"`
	Markers            []string                   `json:"markers"`
	Did                identity.Did               `json:"did"`
	Client             *string                    `json:"client,omitempty"`
	ReplyToClient      bool                       `json:"reply_to_client"`
	ClientReplyAction  *string                    `json:"client_reply_action,omitempty"`
	URL                *string                    `json:"url,omitempty"`
	Multipart          *multipart.Multipart       `json:"multipart,omitempty"`
	ActionValue        *Action                    `json:"action,omitempty"`
	CommandPath        *string                    `json:"command_path,omitempty"`
	Headers            map[string]any             `json:"headers"`
	MessageValue       Message                    `json:"message,omitempty"`
	Sensitivity        int                        `json:"sensitivity"`
	Delayed            bool                       `json:"delayed"`
	MinDelay           int                        `json:"min_delay"`
	MaxDelay           int                        `json:"max_delay"`
	Copy               bool                       `json:"copy"`
	MaxCopies          int                        `json:"max_copies"`
	MinCopies          int                        `json:"min_copies"`
	ServiceLevel       servicestatus.ServiceLevel `json:"service_level"`
}

func newEnvelope(id *string, message Message) *Envelope {
	envID := util.RandomAlphanumeric(32)
	if id != nil {
		envID = *id
	}
	return &Envelope{
		ID:           envID,
		Markers:      []string{},
		Did:          identity.NewDid(),
		Headers:      map[string]any{},
		MessageValue: message,
		Sensitivity:  1,
		ServiceLevel: servicestatus.AtLeastOnce,
	}
}

func CommandEnvelope() *Envelope                 { return newEnvelope(nil, NewCommandMessage(nil)) }
func DocumentEnvelope() *Envelope                { return newEnvelope(nil, NewDocumentMessage()) }
func DocumentEnvelopeWithID(id string) *Envelope { return newEnvelope(&id, NewDocumentMessage()) }
func HeadersOnlyEnvelope() *Envelope             { return newEnvelope(nil, nil) }
func EventEnvelope(eventType string) *Envelope   { return newEnvelope(nil, NewEventMessage(eventType)) }
func TextEnvelope() *Envelope                    { return newEnvelope(nil, &TextMessage{}) }

// ---- headers -------------------------------------------------

func (e *Envelope) SetHeader(name string, value any) { e.Headers[name] = value }
func (e *Envelope) HeaderExists(name string) bool    { _, ok := e.Headers[name]; return ok }
func (e *Envelope) RemoveHeader(name string)         { delete(e.Headers, name) }
func (e *Envelope) Header(name string) any           { return e.Headers[name] }

func (e *Envelope) ContentType() *string {
	if v, ok := e.Headers[HeaderContentType].(string); ok {
		return &v
	}
	return nil
}

func (e *Envelope) SetContentType(contentType string) { e.Headers[HeaderContentType] = contentType }

// ---- routing -----------------------------------------------

// GetRoute returns the current route, resolving it from the routing slip on first access.
func (e *Envelope) GetRoute() route.Route {
	if e.RouteValue == nil {
		e.RouteValue = e.DynamicRoutingSlip.CurrentRoute()
	}
	return e.RouteValue
}

func (e *Envelope) Ratchet() { e.RouteValue = e.DynamicRoutingSlip.NextRoute() }

func (e *Envelope) AddRoute(service, operation string) {
	e.DynamicRoutingSlip.AddRoute(route.SimpleRouteOf(service, operation))
}

func (e *Envelope) AddExternalRoute(service, operation string) {
	e.DynamicRoutingSlip.AddRoute(route.SimpleExternalRouteOf(service, operation))
}

// ---- document payload accessors --------------------------

func (e *Envelope) doc() *DocumentMessage {
	d, _ := e.MessageValue.(*DocumentMessage)
	return d
}

func (e *Envelope) AddContent(content any) bool {
	d := e.doc()
	if d == nil {
		return false
	}
	d.Put(Content, content)
	return true
}

func (e *Envelope) Content() any {
	if d := e.doc(); d != nil {
		return d.Get(Content)
	}
	return nil
}

func (e *Envelope) AddEntity(entity any) bool {
	d := e.doc()
	if d == nil {
		return false
	}
	d.Put(Entity, entity)
	return true
}

func (e *Envelope) Entity() any {
	if d := e.doc(); d != nil {
		return d.Get(Entity)
	}
	return nil
}

func (e *Envelope) AddException(msg string) bool {
	d := e.doc()
	if d == nil {
		return false
	}
	bucket := d.Primary()
	existing, _ := bucket[Exceptions].([]any)
	bucket[Exceptions] = append(existing, msg)
	return true
}

func (e *Envelope) Exceptions() []string {
	d := e.doc()
	if d == nil {
		return nil
	}
	raw, _ := d.Get(Exceptions).([]any)
	out := make([]string, 0, len(raw))
	for _, v := range raw {
		if s, ok := v.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

func (e *Envelope) AddErrorMessage(msg string) {
	if e.MessageValue != nil {
		e.MessageValue.AddErrorMessage(msg)
	}
}

func (e *Envelope) ErrorMessages() []string {
	if e.MessageValue != nil {
		return e.MessageValue.ErrorMessages()
	}
	return nil
}

func (e *Envelope) AddNvp(name string, value any) bool {
	d := e.doc()
	if d == nil {
		return false
	}
	d.Put(name, value)
	return true
}

func (e *Envelope) Value(name string) any {
	if d := e.doc(); d != nil {
		return d.Get(name)
	}
	return nil
}

func (e *Envelope) Values() map[string]any {
	if d := e.doc(); d != nil && len(d.Data) > 0 {
		return d.Data[0]
	}
	return nil
}

// ---- markers ---------------------------------------------

func (e *Envelope) MarkerPresent(marker string) bool {
	for _, m := range e.Markers {
		if m == marker {
			return true
		}
	}
	return false
}

func (e *Envelope) Mark(marker string) { e.Markers = append(e.Markers, marker) }

// ---- serialization ------------------------------------

func (e Envelope) MarshalJSON() ([]byte, error) {
	type alias struct {
		ID                 string                     `json:"id"`
		DynamicRoutingSlip route.DynamicRoutingSlip   `json:"dynamic_routing_slip"`
		Route              json.RawMessage            `json:"route,omitempty"`
		Markers            []string                   `json:"markers"`
		Did                identity.Did               `json:"did"`
		Client             *string                    `json:"client,omitempty"`
		ReplyToClient      bool                       `json:"reply_to_client"`
		ClientReplyAction  *string                    `json:"client_reply_action,omitempty"`
		URL                *string                    `json:"url,omitempty"`
		Multipart          *multipart.Multipart       `json:"multipart,omitempty"`
		Action             *Action                    `json:"action,omitempty"`
		CommandPath        *string                    `json:"command_path,omitempty"`
		Headers            map[string]any             `json:"headers"`
		Message            json.RawMessage            `json:"message,omitempty"`
		Sensitivity        int                        `json:"sensitivity"`
		Delayed            bool                       `json:"delayed"`
		MinDelay           int                        `json:"min_delay"`
		MaxDelay           int                        `json:"max_delay"`
		Copy               bool                       `json:"copy"`
		MaxCopies          int                        `json:"max_copies"`
		MinCopies          int                        `json:"min_copies"`
		ServiceLevel       servicestatus.ServiceLevel `json:"service_level"`
	}
	a := alias{
		ID:                 e.ID,
		DynamicRoutingSlip: e.DynamicRoutingSlip,
		Markers:            e.Markers,
		Did:                e.Did,
		Client:             e.Client,
		ReplyToClient:      e.ReplyToClient,
		ClientReplyAction:  e.ClientReplyAction,
		URL:                e.URL,
		Multipart:          e.Multipart,
		Action:             e.ActionValue,
		CommandPath:        e.CommandPath,
		Headers:            e.Headers,
		Sensitivity:        e.Sensitivity,
		Delayed:            e.Delayed,
		MinDelay:           e.MinDelay,
		MaxDelay:           e.MaxDelay,
		Copy:               e.Copy,
		MaxCopies:          e.MaxCopies,
		MinCopies:          e.MinCopies,
		ServiceLevel:       e.ServiceLevel,
	}
	if e.RouteValue != nil {
		data, err := json.Marshal(e.RouteValue)
		if err != nil {
			return nil, err
		}
		a.Route = data
	}
	if e.MessageValue != nil {
		data, err := json.Marshal(e.MessageValue)
		if err != nil {
			return nil, err
		}
		a.Message = data
	}
	return json.Marshal(a)
}

func (e *Envelope) UnmarshalJSON(data []byte) error {
	var aux struct {
		ID                 string                     `json:"id"`
		DynamicRoutingSlip route.DynamicRoutingSlip   `json:"dynamic_routing_slip"`
		Route              json.RawMessage            `json:"route"`
		Markers            []string                   `json:"markers"`
		Did                identity.Did               `json:"did"`
		Client             *string                    `json:"client"`
		ReplyToClient      bool                       `json:"reply_to_client"`
		ClientReplyAction  *string                    `json:"client_reply_action"`
		URL                *string                    `json:"url"`
		Multipart          *multipart.Multipart       `json:"multipart"`
		Action             *Action                    `json:"action"`
		CommandPath        *string                    `json:"command_path"`
		Headers            map[string]any             `json:"headers"`
		Message            json.RawMessage            `json:"message"`
		Sensitivity        int                        `json:"sensitivity"`
		Delayed            bool                       `json:"delayed"`
		MinDelay           int                        `json:"min_delay"`
		MaxDelay           int                        `json:"max_delay"`
		Copy               bool                       `json:"copy"`
		MaxCopies          int                        `json:"max_copies"`
		MinCopies          int                        `json:"min_copies"`
		ServiceLevel       servicestatus.ServiceLevel `json:"service_level"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	e.ID = aux.ID
	e.DynamicRoutingSlip = aux.DynamicRoutingSlip
	e.Markers = aux.Markers
	if e.Markers == nil {
		e.Markers = []string{}
	}
	e.Did = aux.Did
	e.Client = aux.Client
	e.ReplyToClient = aux.ReplyToClient
	e.ClientReplyAction = aux.ClientReplyAction
	e.URL = aux.URL
	e.Multipart = aux.Multipart
	e.ActionValue = aux.Action
	e.CommandPath = aux.CommandPath
	e.Headers = aux.Headers
	if e.Headers == nil {
		e.Headers = map[string]any{}
	}
	e.Sensitivity = aux.Sensitivity
	e.Delayed = aux.Delayed
	e.MinDelay = aux.MinDelay
	e.MaxDelay = aux.MaxDelay
	e.Copy = aux.Copy
	e.MaxCopies = aux.MaxCopies
	e.MinCopies = aux.MinCopies
	e.ServiceLevel = aux.ServiceLevel
	if e.ServiceLevel == "" {
		e.ServiceLevel = servicestatus.AtLeastOnce
	}
	if len(aux.Route) > 0 {
		r, err := route.UnmarshalRoute(aux.Route)
		if err != nil {
			return err
		}
		e.RouteValue = r
	}
	if len(aux.Message) > 0 {
		m, err := UnmarshalMessage(aux.Message)
		if err != nil {
			return err
		}
		e.MessageValue = m
	}
	return nil
}

func (e *Envelope) ToJSONString() (string, error) {
	data, err := json.MarshalIndent(e, "", "  ")
	return string(data), err
}

func EnvelopeFromJSONString(text string) (*Envelope, error) {
	var e Envelope
	if err := json.Unmarshal([]byte(text), &e); err != nil {
		return nil, err
	}
	return &e, nil
}

func (e *Envelope) Equals(other *Envelope) bool { return e.ID == other.ID }
