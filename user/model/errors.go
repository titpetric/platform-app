package model

import "errors"

// Errors returned by user validation and storage operations.
var (
	ErrSessionExpired    = errors.New("Your session has expired")
	ErrUsernameMissing   = errors.New("username is required")
	ErrUsernameMinLength = errors.New("username must be more than 3 characters")
	ErrUsernameInvalid   = errors.New("username must contain only lowercase letters, numbers, underscores and dashes, and must not begin or end with underscore or dash")
	ErrUsernameMaxLength = errors.New("username must be 20 characters or less")
	ErrUsernameTaken     = errors.New("username is already taken")

	// ErrInvalidActivationToken is returned by UserStorage.Activate when
	// the supplied token does not match any pending activation row.
	ErrInvalidActivationToken = errors.New("invalid activation token")

	// ErrUserAlreadyActivated is returned when a resend or activate
	// operation finds the user is already past the activation gate.
	ErrUserAlreadyActivated = errors.New("user already activated")

	// ErrUserNotActivated is the error surfaced by API handlers when
	// EmailActivationEnabled is on and a user attempts to log in or
	// otherwise act on an account that has not been activated yet.
	ErrUserNotActivated = errors.New("user has not activated their account")
)
