package model

import "time"

// Status constants
const (
	StatusPending = "pending"
	StatusSent    = "sent"
	StatusFailed  = "failed"
)

// NewEmail creates a new email
func NewEmail(recipient, subject, body string) *Email {
	email := &Email{
		Recipient:  recipient,
		Subject:    subject,
		Body:       body,
		Status:     StatusPending,
		RetryCount: 0,
	}
	email.SetCreatedAt(time.Now())
	return email
}
