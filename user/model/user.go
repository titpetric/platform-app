package model

import "fmt"

func NewUser() *User {
	return &User{}
}

func (u *User) String() string {
	if u.DeletedAt != nil {
		return "Deleted user"
	}
	return u.FirstName + " " + u.LastName
}

func (u *User) Ok() bool {
	return u.Validate() == nil
}

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
