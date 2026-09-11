package messaging

const MimetypeTextPlain = "text/plain"

type Email struct {
	To          *string `json:"to,omitempty"`
	From        *string `json:"from,omitempty"`
	Subject     *string `json:"subject,omitempty"`
	Message     *string `json:"message,omitempty"`
	ID          *int64  `json:"id,omitempty"`
	MessageType string  `json:"message_type"`
	Flag        int     `json:"flag"`
}

func NewEmail() Email { return Email{MessageType: MimetypeTextPlain} }

func AnonymousEmail(to, subject, message string) Email {
	e := NewEmail()
	e.To = &to
	e.Subject = &subject
	e.Message = &message
	return e
}
