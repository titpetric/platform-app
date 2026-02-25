package web_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/titpetric/platform-app/user/service/web"
)

func TestRegisterValidationRendersView(t *testing.T) {
	renderer := web.NewRenderer(newViewFS(), nil)

	t.Run("missing all fields renders register view with error", func(t *testing.T) {
		svc := web.NewService(nil, nil, newViewFS())
		_ = renderer

		form := url.Values{}
		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		svc.Register(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Create an account")
	})

	t.Run("missing username renders register view with username error", func(t *testing.T) {
		svc := web.NewService(nil, nil, newViewFS())

		form := url.Values{
			"full_name": {"John Doe"},
			"email":     {"john@example.com"},
			"password":  {"secret123"},
		}
		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		svc.Register(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		body := w.Body.String()
		assert.Contains(t, body, "Create an account")
		assert.Contains(t, body, "username is required")
	})

	t.Run("short username renders register view with error", func(t *testing.T) {
		svc := web.NewService(nil, nil, newViewFS())

		form := url.Values{
			"full_name": {"John Doe"},
			"username":  {"ab"},
			"email":     {"john@example.com"},
			"password":  {"secret123"},
		}
		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		svc.Register(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		body := w.Body.String()
		assert.Contains(t, body, "Create an account")
		assert.Contains(t, body, "username must be more than 3 characters")
	})

	t.Run("form preserves submitted values on error", func(t *testing.T) {
		svc := web.NewService(nil, nil, newViewFS())

		form := url.Values{
			"full_name": {"Jane Smith"},
			"email":     {"jane@example.com"},
			"username":  {"jn"},
			"password":  {"secret123"},
		}
		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		svc.Register(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		body := w.Body.String()
		assert.Contains(t, body, "Jane Smith")
		assert.Contains(t, body, "jane@example.com")
		assert.Contains(t, body, "jn")
	})
}
