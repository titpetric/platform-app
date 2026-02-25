package web

import (
	"net/http"

	"github.com/titpetric/platform/pkg/telemetry"
)

// LogoutView just wraps LoginView. The view changes based on if
// the user is logged in already, allowing them to log in or out.
func (h *Handlers) LogoutView(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.logoutView(w, r))
}

func (h *Handlers) logoutView(w http.ResponseWriter, r *http.Request) error {
	r, span := telemetry.StartRequest(r, "user.service.LogoutView")
	defer span.End()

	return h.loginView(w, r)
}
