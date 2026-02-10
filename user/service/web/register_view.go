package web

import (
	"net/http"

	"github.com/titpetric/platform/pkg/telemetry"
)

// RegisterView renders the registration page.
func (h *Service) RegisterView(w http.ResponseWriter, r *http.Request) {
	r, span := telemetry.StartRequest(r, "user.service.RegisterView")
	defer span.End()

	err := h.view.Register(RegisterData{
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
	if err != nil {
		telemetry.CaptureError(r.Context(), err)
		h.Error(r, "Failed to render register page", err)
	}
}
