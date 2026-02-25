package auth

import (
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

// JWT and Claims provide JWT token creation and validation.
type (
	JWT struct {
		secret        string
		signingMethod *jwt.SigningMethodHMAC
	}

	Claims struct {
		UserID string `json:"user_id"`

		jwt.MapClaims
	}
)

// NewJWT creates a new JWT instance with the given secret.
func NewJWT(secret string) *JWT {
	return &JWT{
		secret:        secret,
		signingMethod: jwt.SigningMethodHS256,
	}
}

// UserID retrieves the `user_id` claim from the JWT token.
func (u *JWT) UserID(token string) (string, error) {
	claims, err := u.Claims(token)
	if err != nil {
		return "", err
	}
	return string(claims.UserID), nil
}

// Claims returns the complete JWT claims object.
func (u *JWT) Claims(tokenString string) (*Claims, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	if tokenString == "" {
		return nil, errEmptyToken
	}
	if u.secret == "" {
		return nil, errEmptySecret
	}

	signingSecret := func(token *jwt.Token) (any, error) {
		return []byte(u.secret), nil
	}

	token, err := jwt.Parse(tokenString, signingSecret, jwt.WithValidMethods([]string{u.signingMethod.Alg()}))
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if userID, ok := claims["user_id"].(string); ok && userID != "" {
			return &Claims{
				MapClaims: claims,
				UserID:    userID,
			}, nil
		}
	}

	return nil, errInvalidClaims
}

// Validate checks if the JWT claims match a userID.
func (u *JWT) Validate(token string, userID string) (bool, error) {
	uid, err := u.UserID(token)
	if err != nil {
		return false, err
	}
	return uid == userID, nil
}

// IsUser is a simpler version of Validate, discarding the error.
func (u *JWT) IsUser(token string, userID string) bool {
	isUser, _ := u.Validate(token, userID)
	return isUser
}

// Create generates a signed JWT token for the given userID with the specified TTL.
func (u *JWT) Create(userID string, ttl time.Duration) (string, error) {
	signingSecret := func() []byte {
		return []byte(u.secret)
	}

	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(ttl).Unix()

	at := jwt.NewWithClaims(u.signingMethod, claims)
	return at.SignedString(signingSecret())
}
