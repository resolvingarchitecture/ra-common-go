// Package raerror is the package-wide error type. Replaces the family of
// checked exception classes in ra-common-java (ServiceNotFoundException,
// FileCreationFailedException, ...).
package raerror

import "fmt"

type Kind string

const (
	Decode                   Kind = "decode"
	Crypto                   Kind = "crypto"
	Invalid                  Kind = "invalid"
	ServiceNotFound          Kind = "service-not-found"
	ServiceNotAccessible     Kind = "service-not-accessible"
	ServiceNotSupported      Kind = "service-not-supported"
	ServiceAlreadyRegistered Kind = "service-already-registered"
	FileCreationFailed       Kind = "file-creation-failed"
	IO                       Kind = "io"
)

type RaError struct {
	Kind    Kind
	Message string
}

func (e *RaError) Error() string {
	return fmt.Sprintf("%s: %s", e.Kind, e.Message)
}

func New(kind Kind, message string) *RaError {
	return &RaError{Kind: kind, Message: message}
}

func DecodeErr(message string) *RaError  { return New(Decode, message) }
func CryptoErr(message string) *RaError  { return New(Crypto, message) }
func InvalidErr(message string) *RaError { return New(Invalid, message) }
