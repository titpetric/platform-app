package user

import "errors"

// ErrLoginRequired is returned with RequireLoginError middleware.
var ErrLoginRequired = errors.New("login required")
