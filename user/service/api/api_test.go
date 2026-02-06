package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/titpetric/platform/pkg/require"

	"github.com/titpetric/platform-app/user/service/auth"
)

func getTestSigningKey() string {
	return "test-signing-key"
}

func TestRefreshToken(t *testing.T) {
	t.Parallel()

	t.Run("valid token", func(t *testing.T) {
		userID := "user-456"
		signingKey := getTestSigningKey()

		jwtAuth := auth.NewJWT(signingKey)
		originalToken, err := jwtAuth.Create(userID, time.Hour)
		require.NoError(t, err)

		svc := &Service{
			opts: Options{SigningKey: signingKey},
		}

		req := httptest.NewRequest(http.MethodPost, "/api/user/token/refresh", nil)
		req.Header.Set("Authorization", "Bearer "+originalToken)
		w := httptest.NewRecorder()

		svc.RefreshToken(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Token     string `json:"token"`
			ExpiresAt int64  `json:"expires_at"`
		}
		err = json.NewDecoder(w.Body).Decode(&resp)
		require.NoError(t, err)
		require.True(t, resp.Token != "")
		require.True(t, resp.ExpiresAt > time.Now().Unix())

		extractedUserID, err := jwtAuth.UserID(resp.Token)
		require.NoError(t, err)
		require.Equal(t, userID, extractedUserID)
	})

	t.Run("missing authorization header", func(t *testing.T) {
		svc := &Service{
			opts: Options{SigningKey: getTestSigningKey()},
		}

		req := httptest.NewRequest(http.MethodPost, "/api/user/token/refresh", nil)
		w := httptest.NewRecorder()

		svc.RefreshToken(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid token", func(t *testing.T) {
		svc := &Service{
			opts: Options{SigningKey: getTestSigningKey()},
		}

		req := httptest.NewRequest(http.MethodPost, "/api/user/token/refresh", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		w := httptest.NewRecorder()

		svc.RefreshToken(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestCreateTokenInvalidBody(t *testing.T) {
	t.Parallel()

	svc := &Service{
		opts: Options{SigningKey: getTestSigningKey()},
	}

	body := `invalid json`
	req := httptest.NewRequest(http.MethodPost, "/api/user/token/create", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	svc.CreateToken(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegisterInvalidBody(t *testing.T) {
	t.Parallel()

	svc := &Service{
		opts: Options{SigningKey: getTestSigningKey()},
	}

	body := `invalid json`
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	svc.Register(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegisterMissingFields(t *testing.T) {
	t.Parallel()

	svc := &Service{
		opts: Options{SigningKey: getTestSigningKey()},
	}

	body := `{"first_name": "John", "last_name": "Doe"}`
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	svc.Register(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}
