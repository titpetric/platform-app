package model

import (
	"github.com/go-webauthn/webauthn/webauthn"
)

// WebAuthnUser wraps a User with passkey credentials for the webauthn.User interface.
type WebAuthnUser struct {
	*User
	Passkeys []UserPasskey
}

// WebAuthnID returns the user handle as bytes.
func (u *WebAuthnUser) WebAuthnID() []byte {
	return []byte(u.ID)
}

// WebAuthnName returns the username.
func (u *WebAuthnUser) WebAuthnName() string {
	return u.Username
}

// WebAuthnDisplayName returns the full name.
func (u *WebAuthnUser) WebAuthnDisplayName() string {
	return u.FullName
}

// WebAuthnCredentials returns the user's stored passkey credentials.
func (u *WebAuthnUser) WebAuthnCredentials() []webauthn.Credential {
	creds := make([]webauthn.Credential, len(u.Passkeys))
	for i, p := range u.Passkeys {
		creds[i] = p.ToCredential()
	}
	return creds
}
