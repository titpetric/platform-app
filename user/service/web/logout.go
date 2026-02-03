package web

import (
	"net/http"

	"github.com/titpetric/platform/pkg/telemetry"
)

// Logout deletes the session cookie and optionally the session in storage.
func (h *Service) Logout(w http.ResponseWriter, r *http.Request) {
	r, span := telemetry.StartRequest(r, "user.service.Logout")
	defer span.End()

	ctx := r.Context()

	cookie, err := r.Cookie("session_id")

	if err == nil && cookie.Value != "" {
		// Delete the session from the database
		_ = h.sessionStorage.Delete(ctx, cookie.Value)

		// Clear cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1,
		})
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
