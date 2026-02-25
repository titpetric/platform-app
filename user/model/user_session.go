package model

import (
	"fmt"
	"time"
)

// NewUserSession creates a new empty UserSession.
func NewUserSession() *UserSession {
	return &UserSession{}
}

// Ok reports whether the UserSession is valid.
func (u *UserSession) Ok() bool {
	return u.Validate() == nil
}

// Validate checks that the UserSession is non-nil, populated, and not expired.
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
