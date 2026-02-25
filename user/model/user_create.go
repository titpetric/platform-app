package model

import (
	"regexp"
	"strings"
)

// UserCreateRequest holds the fields required to create a new user.
type UserCreateRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username,omitempty"`
}

// Valid reports whether the UserCreateRequest has all required fields.
func (r *UserCreateRequest) Valid() bool {
	if r == nil {
		return false
	}
	if r.FullName == "" {
		return false
	}
	if r.Email == "" || r.Password == "" {
		return false
	}
	if r.Username == "" {
		return false
	}
	return r.ValidateUsername() == nil
}

// ValidateUsername checks that the username meets length and format requirements.
func (r *UserCreateRequest) ValidateUsername() error {
	if r.Username == "" {
		return ErrUsernameMissing
	}
	if len(r.Username) < 4 {
		return ErrUsernameMinLength
	}
	if len(r.Username) > 20 {
		return ErrUsernameMaxLength
	}
	if !regexp.MustCompile(`^[a-z0-9][a-z0-9_-]*[a-z0-9]$|^[a-z0-9]{3}$`).MatchString(r.Username) {
		return ErrUsernameInvalid
	}
	return nil
}

// User converts the UserCreateRequest into a User.
func (r *UserCreateRequest) User() *User {
	username := r.Username
	if username == "" {
		username = r.Email
	}
	return &User{
		FullName: r.FullName,
		Username: username,
		Slug:     strings.ToLower(username),
	}
}

// UserAuth converts the UserCreateRequest into a UserAuth.
func (r *UserCreateRequest) UserAuth() *UserAuth {
	return &UserAuth{
		Email:    r.Email,
		Password: r.Password,
	}
}
