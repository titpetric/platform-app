package model

import "errors"

var (
	ErrSessionExpired    = errors.New("Your session has expired")
	ErrUsernameMissing   = errors.New("username is required")
	ErrUsernameMinLength = errors.New("username must be more than 3 characters")
	ErrUsernameInvalid   = errors.New("username must contain only lowercase letters, numbers, underscores and dashes, and must not begin or end with underscore or dash")
	ErrUsernameTaken     = errors.New("username is already taken")
)
