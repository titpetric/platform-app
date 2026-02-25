package web

import (
	"net/http"

	"github.com/titpetric/platform/pkg/telemetry"
)

// RegisterView renders the registration page.
func (h *Handlers) RegisterView(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.registerView(w, r))
}

func (h *Handlers) registerView(w http.ResponseWriter, r *http.Request) error {
	r, span := telemetry.StartRequest(r, "user.service.RegisterView")
	defer span.End()

	return h.view.Register(RegisterData{
		ErrorMessage: h.GetError(r),
		FullName:     r.FormValue("full_name"),
		Email:        r.FormValue("email"),
		Username:     r.FormValue("username"),
		Links: Links{
			Login:    "/login",
			Logout:   "/logout",
			Register: "/register",
		},
	}).Render(r.Context(), w)
}
