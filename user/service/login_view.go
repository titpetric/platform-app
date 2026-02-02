package service

import (
	"net/http"

	"github.com/titpetric/platform/pkg/telemetry"
)

// LoginView renders login.tpl when no valid session exists,
// or logout.tpl with the full user model when a valid session is found.
func (h *Service) LoginView(w http.ResponseWriter, r *http.Request) {
	r, span := telemetry.StartRequest(r, "user.service.LoginView")
	defer span.End()

	ctx := r.Context()

	cookie, err := r.Cookie("session_id")
	if err == nil && cookie.Value != "" {
		if session, err := h.SessionStorage.Get(ctx, cookie.Value); err == nil {
			if user, err := h.UserStorage.Get(ctx, session.UserID); err == nil {
				if err := h.view.Logout(LogoutData{
					SessionUser: user,
					Links: Links{
						Login:    "/login",
						Logout:   "/logout",
						Register: "/register",
					},
				}).Render(ctx, w); err != nil {
					telemetry.CaptureError(ctx, err)
					h.Error(r, "Failed to render logout page", err)
				}
				return
			} else {
				telemetry.CaptureError(ctx, err)
			}
		} else {
			telemetry.CaptureError(ctx, err)
		}
	}

	if err := h.view.Login(LoginData{
		ErrorMessage: h.GetError(r),
		Email:        r.FormValue("email"),
		Links: Links{
			Login:    "/login",
			Logout:   "/logout",
			Register: "/register",
		},
	}).Render(ctx, w); err != nil {
		telemetry.CaptureError(ctx, err)
		h.Error(r, "Failed to render login page", err)
	}
}
