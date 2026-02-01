package model

import (
	"fmt"
	"time"
)

func NewUserSession() *UserSession {
	return &UserSession{}
}

func (u *UserSession) Ok() bool {
	return u.Validate() == nil
}

func (u *UserSession) Validate() error {
	if u == nil {
		return fmt.Errorf("no session")
	}
	if u.ID == "" || u.UserID == "" {
		return fmt.Errorf("session is empty")
	}
	if u.ExpiresAt != nil && time.Since(*u.ExpiresAt) > 0 {
		return fmt.Errorf("session is expired")
	}
	return nil
}
