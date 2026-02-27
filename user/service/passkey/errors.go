package passkey

// Error represents an HTTP error with a status code.
type Error struct {
	Status int
	Err    error
}

// Error returns the underlying error message.
func (e *Error) Error() string {
	return e.Err.Error()
}
