package model

import "strings"

type UserCreateRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username,omitempty"`
}

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
	return true
}

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

func (r *UserCreateRequest) UserAuth() *UserAuth {
	return &UserAuth{
		Email:    r.Email,
		Password: r.Password,
	}
}
