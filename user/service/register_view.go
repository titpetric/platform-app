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

	h.view.Register(view.RegisterData{
		ErrorMessage: h.GetError(r),
		FirstName:    r.FormValue("first_name"),
		LastName:     r.FormValue("last_name"),
		Email:        r.FormValue("email"),
	}).Render(r.Context(), w)
}
