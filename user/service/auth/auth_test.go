package auth

import (
	"fmt"
	"os"
	"testing"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/titpetric/platform/pkg/require"
)

func getJwtSecret() string {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "test-usage"
	}
	return jwtSecret
}

func getJwtUserClaim(userID string) jwt.MapClaims {
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	return claims
}

func getJwt(claims jwt.MapClaims, secret string) (string, error) {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return at.SignedString([]byte(secret))
}

func TestAuth(t *testing.T) {
	t.Parallel()

	uid := os.Getenv("JWT_USER")
	if uid == "" {
		uid = "test"
	}

	// generate tokens with test code and pkg code
	tokens := []func() (string, error){
		func() (string, error) {
			jwtSecret := getJwtSecret()
			jwtClaims := getJwtUserClaim(uid)
			return getJwt(jwtClaims, jwtSecret)
		},
		func() (string, error) {
			return NewJWT(getJwtSecret()).Create(uid, time.Hour)
		},
	}

	for idx, tokFn := range tokens {
		t.Run(fmt.Sprintf("token: %d", idx), func(t *testing.T) {
			token, err := tokFn()
			require.NoError(t, err)

			validator := NewJWT(getJwtSecret())
			require.True(t, validator.IsUser(token, uid))

			user, err := validator.Claims(token)
			require.NoError(t, err)

			t.Logf("Generated JWT: %s", token)
			t.Logf("Claims: %d", len(user.MapClaims))
			for idx, claim := range user.MapClaims {
				t.Logf(" - %s: %v (%T)", idx, claim, claim)
			}
		})
	}
}
