package admin

import (
	"log"
	"net/http"

	"github.com/titpetric/platform/pkg/telemetry"
)

// Error represents an HTTP error with a status code and optional cause.
type Error struct {
	StatusCode int
	Message    string
	Cause      error `json:"-"`
}

// Error returns the error message.
func (e *Error) Error() string {
	return e.Message
}

// Unwrap returns the underlying cause.
func (e *Error) Unwrap() error {
	return e.Cause
}

// NewError creates an Error with the given status, message, and cause.
func NewError(status int, message string, cause error) *Error {
	return &Error{
		StatusCode: status,
		Message:    message,
		Cause:      cause,
	}
}

// ErrBadRequest creates a 400 Bad Request error.
func ErrBadRequest(message string, cause error) *Error {
	return NewError(http.StatusBadRequest, message, cause)
}

// ErrUnauthorized creates a 401 Unauthorized error.
func ErrUnauthorized(message string, cause error) *Error {
	return NewError(http.StatusUnauthorized, message, cause)
}

// ErrForbidden creates a 403 Forbidden error.
func ErrForbidden(message string, cause error) *Error {
	return NewError(http.StatusForbidden, message, cause)
}

// ErrNotFound creates a 404 Not Found error.
func ErrNotFound(message string, cause error) *Error {
	return NewError(http.StatusNotFound, message, cause)
}

// ErrInternal creates a 500 Internal Server Error.
func ErrInternal(message string, cause error) *Error {
	return NewError(http.StatusInternalServerError, message, cause)
}

func (h *Handlers) errorHandler(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	ctx := r.Context()

	switch val := err.(type) {
	case *Error:
		if val.Cause != nil {
			log.Printf("error: %s (cause: %v)", val.Message, val.Cause)
			telemetry.CaptureError(ctx, val.Cause)
		} else {
			log.Printf("error: %s", val.Message)
		}
		http.Error(w, val.Message, val.StatusCode)
	default:
		log.Printf("error: %v", err)
		telemetry.CaptureError(ctx, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
