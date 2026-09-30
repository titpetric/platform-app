package auth

import (
	"fmt"
	"os"
	"testing"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"

	"github.com/titpetric/platform-app/internal/assert"
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
			assert.NoError(t, err)

			validator := NewJWT(getJwtSecret())
			assert.True(t, validator.IsUser(token, uid))

			user, err := validator.Claims(token)
			assert.NoError(t, err)

			t.Logf("Generated JWT: %s", token)
			t.Logf("Claims: %d", len(user.MapClaims))
			for idx, claim := range user.MapClaims {
				t.Logf(" - %s: %v (%T)", idx, claim, claim)
			}
		})
	}
}

func TestCreateWithJTI(t *testing.T) {
	t.Parallel()

	j := NewJWT(getJwtSecret())

	token, jti, err := j.CreateWithJTI("user-1", time.Hour)
	assert.NoError(t, err)
	assert.True(t, token != "")
	assert.True(t, jti != "")

	claims, err := j.Claims(token)
	assert.NoError(t, err)
	assert.Equal(t, "user-1", claims.UserID)
	assert.Equal(t, jti, claims.JTI)
	assert.True(t, claims.ExpiresAt > time.Now().Unix())
}

func TestCreateAlwaysEmbedsJTI(t *testing.T) {
	t.Parallel()

	j := NewJWT(getJwtSecret())

	a, err := j.Create("user-a", time.Hour)
	assert.NoError(t, err)
	b, err := j.Create("user-a", time.Hour)
	assert.NoError(t, err)

	ca, err := j.Claims(a)
	assert.NoError(t, err)
	cb, err := j.Claims(b)
	assert.NoError(t, err)

	assert.True(t, ca.JTI != "")
	assert.True(t, cb.JTI != "")
	assert.True(t, ca.JTI != cb.JTI)
}

func TestClaimsBackCompatNoJTI(t *testing.T) {
	t.Parallel()

	// Token generated externally (no jti, mirrors pre-rollout tokens).
	claims := jwt.MapClaims{
		"user_id": "legacy",
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	tok, err := getJwt(claims, getJwtSecret())
	assert.NoError(t, err)

	parsed, err := NewJWT(getJwtSecret()).Claims(tok)
	assert.NoError(t, err)
	assert.Equal(t, "legacy", parsed.UserID)
	assert.Equal(t, "", parsed.JTI)
}
