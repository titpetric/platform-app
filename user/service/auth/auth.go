package auth

import (
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/titpetric/platform/pkg/ulid"
)

// JWT and Claims provide JWT token creation and validation.
type (
	JWT struct {
		secret        string
		signingMethod *jwt.SigningMethodHMAC
	}

	Claims struct {
		UserID string `json:"user_id"`
		// JTI is the JWT identifier, used for revocation. Tokens issued
		// before JTI rollout will have an empty JTI; callers should treat
		// an empty JTI as "not revocable" rather than as an error.
		JTI string `json:"jti"`
		// ExpiresAt is the unix timestamp from the `exp` claim, if present.
		ExpiresAt int64 `json:"exp"`

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
			c := &Claims{
				MapClaims: claims,
				UserID:    userID,
			}
			if jti, ok := claims["jti"].(string); ok {
				c.JTI = jti
			}
			if exp, ok := claims["exp"].(float64); ok {
				c.ExpiresAt = int64(exp)
			}
			return c, nil
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

// Create generates a signed JWT token for the given userID with the specified
// TTL. A fresh JTI (JWT ID) is embedded so the token can be individually
// revoked. Use CreateWithJTI when the caller needs to know the JTI in
// advance (e.g. to insert it into a revocation table later).
func (u *JWT) Create(userID string, ttl time.Duration) (string, error) {
	token, _, err := u.CreateWithJTI(userID, ttl)
	return token, err
}

// CreateWithJTI generates a signed JWT and returns both the encoded token
// and its JTI claim. The JTI is a ULID, lex-sortable by issue time.
func (u *JWT) CreateWithJTI(userID string, ttl time.Duration) (string, string, error) {
	signingSecret := func() []byte {
		return []byte(u.secret)
	}

	jti := ulid.String()
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["jti"] = jti
	claims["exp"] = time.Now().Add(ttl).Unix()

	at := jwt.NewWithClaims(u.signingMethod, claims)
	signed, err := at.SignedString(signingSecret())
	if err != nil {
		return "", "", err
	}
	return signed, jti, nil
}
