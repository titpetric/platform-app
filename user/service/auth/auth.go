package auth

import (
	"strings"

	jwt "github.com/golang-jwt/jwt/v5"
)

type (
	JWT struct {
		secret string
	}

	Claims struct {
		UserID string `json:"user_id"`

		jwt.MapClaims
	}
)

func NewJWT(secret string) *JWT {
	return &JWT{
		secret: secret,
	}
}

// UserID retrieves the `user_id` claim from the JWT token
func (u *JWT) UserID(token string) (string, error) {
	claims, err := u.Claims(token)
	if err != nil {
		return "", err
	}
	return string(claims.UserID), nil
}

// Claims returns the complete JWT claims object
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

	token, err := jwt.Parse(tokenString, signingSecret, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
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

// Validate just checks if the JWT claims match an userID
func (u *JWT) Validate(token string, userID string) (bool, error) {
	uid, err := u.UserID(token)
	if err != nil {
		return false, err
	}
	return uid == userID, nil
}

// IsUser is a simpler version of Validate, throwing away the error
func (u *JWT) IsUser(token string, userID string) bool {
	isUser, _ := u.Validate(token, userID)
	return isUser
}
