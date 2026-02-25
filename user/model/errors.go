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
)
