package model

import (
	"testing"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/titpetric/platform/pkg/require"
)

func TestUserPasskeyToCredential(t *testing.T) {
	passkey := &UserPasskey{
		ID:              "pk-1",
		UserID:          "user-1",
		CredentialID:    []byte("cred-id"),
		PublicKey:       []byte("pub-key"),
		AttestationType: "none",
		Transport:       `["internal"]`,
		SignCount:       5,
	}

	cred := passkey.ToCredential()
	require.Equal(t, passkey.CredentialID, cred.ID)
	require.Equal(t, passkey.PublicKey, cred.PublicKey)
	require.Equal(t, passkey.AttestationType, cred.AttestationType)
	require.Equal(t, uint32(5), cred.Authenticator.SignCount)
}

func TestUserPasskeyToCredentialInvalidJSON(t *testing.T) {
	passkey := &UserPasskey{
		Transport: "invalid-json",
	}

	// Should not panic, just return empty transports
	cred := passkey.ToCredential()
	require.Equal(t, 0, len(cred.Transport))
}

func TestTransportJSON(t *testing.T) {
	transports := []protocol.AuthenticatorTransport{
		protocol.Internal,
		protocol.USB,
	}

	json := TransportJSON(transports)
	require.Equal(t, `["internal","usb"]`, json)
}

func TestTransportJSONEmpty(t *testing.T) {
	json := TransportJSON(nil)
	require.Equal(t, "null", json)
}
