package api

// RequestError represents an HTTP error with a status code.
type RequestError struct {
	StatusCode int

	Err error
}

// Error returns the underlying error message.
func (r *RequestError) Error() string {
	return r.Err.Error()
}
