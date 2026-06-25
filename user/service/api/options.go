package api

import (
	"context"
	"time"

	"github.com/titpetric/platform-app/user/service/passkey"
	"github.com/titpetric/platform-app/user/storage"
)

// EmailSender is the minimal contract used by the activation flow.
// It mirrors service.EmailSender but is re-declared here so this
// package does not import its parent.
type EmailSender interface {
	Send(ctx context.Context, recipient, subject, body string) error
}

// Options is passed from user service scope.
type Options struct {
	SigningKey     string
	TokenTTL       time.Duration
	UserStorage    *storage.UserStorage
	SessionStorage *storage.SessionStorage
	RevokedStorage *storage.RevokedTokenStorage
	PasskeyService *passkey.Service

	// Activation configuration; see service.Options.
	EmailActivationEnabled bool
	EmailSender            EmailSender
	ActivationURLFormat    string
	ActivationSubject      string
}
