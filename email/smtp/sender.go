package smtp

import (
	"fmt"
	"net/smtp"
)

// Sender interface defines methods for sending emails
type Sender interface {
	Send(recipient, subject, body string) error
}

// SMTPSender implements the Sender interface using SMTP
type SMTPSender struct {
	config Config
}

// NewSMTPSender creates a new SMTP sender
func NewSMTPSender(config Config) *SMTPSender {
	return &SMTPSender{config: config}
}

// Send sends an email via SMTP
func (s *SMTPSender) Send(recipient, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	// Create message
	message := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", recipient, subject, body)

	// For mailhog, no authentication is needed
	var auth smtp.Auth
	if s.config.Username != "" && s.config.Password != "" {
		auth = smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
	}

	// Send the email
	err := smtp.SendMail(addr, auth, s.config.From, []string{recipient}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
