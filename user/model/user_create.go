package model

import "strings"

type UserCreateRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Username  string `json:"username,omitempty"`
}

func (r *UserCreateRequest) Valid() bool {
	if r == nil {
		return false
	}
	if r.FirstName == "" || r.LastName == "" {
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
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Username:  username,
		Slug:      strings.ToLower(username),
	}
}

func (r *UserCreateRequest) UserAuth() *UserAuth {
	return &UserAuth{
		Email:    r.Email,
		Password: r.Password,
	}
}
