package model

import (
	"encoding/json"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

// ToCredential converts a stored UserPasskey to a webauthn.Credential.
func (p *UserPasskey) ToCredential() webauthn.Credential {
	var transports []protocol.AuthenticatorTransport
	_ = json.Unmarshal([]byte(p.Transport), &transports)

	return webauthn.Credential{
		ID:              p.CredentialID,
		PublicKey:       p.PublicKey,
		AttestationType: p.AttestationType,
		Transport:       transports,
		Authenticator: webauthn.Authenticator{
			SignCount: uint32(p.SignCount),
		},
	}
}

// TransportJSON returns the JSON encoding of the given transports.
func TransportJSON(transports []protocol.AuthenticatorTransport) string {
	b, err := json.Marshal(transports)
	if err != nil {
		return "[]"
	}
	return string(b)
}
