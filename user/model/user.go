package model

import "fmt"

// NewUser creates a new empty User.
func NewUser() *User {
	return &User{}
}

// String returns a string representation of the User.
func (u *User) String() string {
	if u.DeletedAt != nil {
		return "Deleted user"
	}
	return u.FullName
}

// Ok reports whether the User is valid.
func (u *User) Ok() bool {
	return u.Validate() == nil
}

// Validate checks that the User is non-nil, has an ID, and is not deleted.
func (u *User) Validate() error {
	if u == nil {
		return fmt.Errorf("user is empty")
	}
	if u.ID == "" {
		return fmt.Errorf("user id is empty")
	}
	if u.DeletedAt != nil {
		return fmt.Errorf("user is deleted")
	}
	return nil
}
