package service

import (
	"net/http"

	"github.com/titpetric/platform/pkg/telemetry"

	"github.com/titpetric/platform-app/user/view"
)

// RegisterView renders the registration page.
func (h *Service) RegisterView(w http.ResponseWriter, r *http.Request) {
	r, span := telemetry.StartRequest(r, "user.service.RegisterView")
	defer span.End()

	err := h.view.Register(view.RegisterData{
		ErrorMessage: h.GetError(r),
		FirstName:    r.FormValue("first_name"),
		LastName:     r.FormValue("last_name"),
		Email:        r.FormValue("email"),
		Links: view.Links{
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
