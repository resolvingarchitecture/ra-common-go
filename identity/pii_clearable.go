package identity

// PiiClearable is implemented by anything carrying personally-identifiable
// information that can be scrubbed.
type PiiClearable interface {
	ClearSensitive()
}
