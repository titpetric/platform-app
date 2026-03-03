package model

import (
	"testing"

	"github.com/titpetric/platform/pkg/require"
)

func TestWebAuthnUser(t *testing.T) {
	user := &User{
		ID:       "user-123",
		Username: "testuser",
		FullName: "Test User",
	}

	passkeys := []UserPasskey{
		{
			ID:              "pk-1",
			CredentialID:    []byte("cred-1"),
			PublicKey:       []byte("key-1"),
			AttestationType: "none",
			Transport:       `["internal"]`,
			SignCount:       1,
		},
		{
			ID:              "pk-2",
			CredentialID:    []byte("cred-2"),
			PublicKey:       []byte("key-2"),
			AttestationType: "none",
			Transport:       `["usb"]`,
			SignCount:       2,
		},
	}

	waUser := &WebAuthnUser{
		User:     user,
		Passkeys: passkeys,
	}

	require.Equal(t, []byte("user-123"), waUser.WebAuthnID())
	require.Equal(t, "testuser", waUser.WebAuthnName())
	require.Equal(t, "Test User", waUser.WebAuthnDisplayName())

	creds := waUser.WebAuthnCredentials()
	require.Equal(t, 2, len(creds))
	require.Equal(t, []byte("cred-1"), creds[0].ID)
	require.Equal(t, []byte("cred-2"), creds[1].ID)
}

func TestWebAuthnUserNoPasskeys(t *testing.T) {
	user := &User{
		ID:       "user-456",
		Username: "nopasskeys",
		FullName: "No Passkeys",
	}

	waUser := &WebAuthnUser{
		User:     user,
		Passkeys: nil,
	}

	creds := waUser.WebAuthnCredentials()
	require.Equal(t, 0, len(creds))
}
