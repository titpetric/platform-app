package service

import (
	"context"
	"time"
)

// Options is passed from user package scope. Every field has a defensible
// zero value so a bare Options{SigningKey: ...} keeps the historical
// behaviour. Toggles can be set explicitly to opt into newer flows.
type Options struct {
	SigningKey string

	// TokenTTL is the lifetime for issued access JWTs. If zero,
	// DefaultTokenTTL is used.
	TokenTTL time.Duration

	// EmailActivationEnabled gates account activation behind an email
	// confirmation step. When false (the default), users are activated
	// on creation; the email may still be confirmed via a separate
	// flow added later. When true, registration leaves the user
	// pending until they exchange an activation token sent via the
	// configured EmailSender, and login is refused for pending users.
	EmailActivationEnabled bool

	// EmailSender delivers transactional mail. The email module's
	// *email.Handler satisfies this interface. When
	// EmailActivationEnabled is true and EmailSender is nil, the user
	// module fails to start so the misconfiguration is loud.
	EmailSender EmailSender

	// ActivationURLFormat is a Sprintf-style template that the user
	// module substitutes the activation token into when composing the
	// activation email. Example:
	//   "https://example.com/activate?token=%s"
	// When empty, the activation email contains the bare token and the
	// caller is expected to provide UX instructions out-of-band.
	ActivationURLFormat string

	// ActivationSubject overrides the subject line of activation
	// emails. When empty, DefaultActivationSubject is used.
	ActivationSubject string
}

// EmailSender is the minimal contract the user module needs to deliver
// transactional mail. It deliberately mirrors the signature of the
// email module's Handler.Send so a *email.Handler can be passed
// directly without an adapter.
type EmailSender interface {
	Send(ctx context.Context, recipient, subject, body string) error
}

// Default values applied when the corresponding Options fields are zero.
// Kept exported so callers building custom configurations can compose
// against the same defaults.
const (
	// DefaultTokenTTL preserves the previously-hardcoded value of 30 days.
	DefaultTokenTTL = 30 * 24 * time.Hour

	// DefaultActivationSubject is used when Options.ActivationSubject is empty.
	DefaultActivationSubject = "Confirm your account"
)
