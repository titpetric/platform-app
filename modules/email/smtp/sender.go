package smtp

import (
	"fmt"
	"net/smtp"
	"os"
	"strconv"
)

// Sender interface defines methods for sending emails
type Sender interface {
	Send(recipient, subject, body string) error
}

// Config holds SMTP configuration
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// ConfigFromEnv creates a Config from environment variables
func ConfigFromEnv() Config {
	return Config{
		Host:     getEnv("PLATFORM_EMAIL_HOST", "localhost"),
		Port:     getEnvInt("PLATFORM_EMAIL_PORT", 1025),
		Username: getEnv("PLATFORM_EMAIL_USERNAME", ""),
		Password: getEnv("PLATFORM_EMAIL_PASSWORD", ""),
		From:     getEnv("PLATFORM_EMAIL_FROM", "noreply@example.com"),
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvInt gets an integer environment variable with a default value
func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
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
