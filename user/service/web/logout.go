package web

import (
	"net/http"

	"github.com/titpetric/platform/pkg/telemetry"
)

// Logout deletes the session cookie and optionally the session in storage.
func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.logout(w, r))
}

func (h *Handlers) logout(w http.ResponseWriter, r *http.Request) error {
	r, span := telemetry.StartRequest(r, "user.service.Logout")
	defer span.End()

	ctx := r.Context()

	cookie, err := r.Cookie("session_id")

	if err == nil && cookie.Value != "" {
		_ = h.sessionStorage.Delete(ctx, cookie.Value)

		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1,
		})
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return nil
}
