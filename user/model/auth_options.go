package model

// AuthOptions holds incoming parameters for the service.
// Best keep it flat, requires some field prefixing.
type AuthOptions struct {
	// Cookie if true will read in CookieName for a session ID.
	Cookie     bool
	CookieName string

	// Header if true will read in HeaderName for a request auth token.
	// Authorization already produces the user ID based on JWT claims.
	Header     bool
	HeaderName string

	// This is a JWT specific option.
	HeaderSigningKey string
}
