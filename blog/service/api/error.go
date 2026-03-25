package api

import (
	"log"
	"net/http"

	"github.com/titpetric/platform"
	"github.com/titpetric/platform/pkg/telemetry"
)

// Error represents an HTTP error with a status code and optional cause.
type Error struct {
	StatusCode int
	Message    string
	Cause      error `json:"-"`
}

func (e *Error) Error() string {
	return e.Message
}

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
		platform.Error(w, r, val.StatusCode, val)
	default:
		log.Printf("error: %v", err)
		telemetry.CaptureError(ctx, err)
		platform.Error(w, r, http.StatusInternalServerError, err)
	}
}
