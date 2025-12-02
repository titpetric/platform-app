package model

import "time"

// Email represents an email to be sent
type Email struct {
	ID         string     `db:"id"`
	Recipient  string     `db:"recipient"`
	Subject    string     `db:"subject"`
	Body       string     `db:"body"`
	Status     string     `db:"status"` // pending, sent, failed
	CreatedAt  time.Time  `db:"created_at"`
	SentAt     *time.Time `db:"sent_at"`
	Error      *string    `db:"error"`
	RetryCount int        `db:"retry_count"`
	LastError  *string    `db:"last_error"`
	LastRetry  *time.Time `db:"last_retry"`
}

// Status constants
const (
	StatusPending = "pending"
	StatusSent    = "sent"
	StatusFailed  = "failed"
)

// NewEmail creates a new email
func NewEmail(recipient, subject, body string) *Email {
	return &Email{
		Recipient:  recipient,
		Subject:    subject,
		Body:       body,
		Status:     StatusPending,
		CreatedAt:  time.Now(),
		RetryCount: 0,
	}
}
