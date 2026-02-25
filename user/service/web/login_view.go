package web

import (
	"net/http"

	"github.com/titpetric/platform/pkg/telemetry"
)

// LoginView renders login.tpl when no valid session exists,
// or logout.tpl with the full user model when a valid session is found.
func (h *Handlers) LoginView(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.loginView(w, r))
}

func (h *Handlers) loginView(w http.ResponseWriter, r *http.Request) error {
	r, span := telemetry.StartRequest(r, "user.service.LoginView")
	defer span.End()

	ctx := r.Context()

	cookie, err := r.Cookie("session_id")
	if err == nil && cookie.Value != "" {
		if session, err := h.sessionStorage.Get(ctx, cookie.Value); err == nil {
			if user, err := h.userStorage.Get(ctx, session.UserID); err == nil {
				return h.view.Logout(LogoutData{
					SessionUser: user,
					Links: Links{
						Login:    "/login",
						Logout:   "/logout",
						Register: "/register",
					},
				}).Render(ctx, w)
			} else {
				telemetry.CaptureError(ctx, err)
			}
		} else {
			telemetry.CaptureError(ctx, err)
		}
	}

	return h.view.Login(LoginData{
		ErrorMessage: h.GetError(r),
		Email:        r.FormValue("email"),
		Links: Links{
			Login:    "/login",
			Logout:   "/logout",
			Register: "/register",
		},
	}).Render(ctx, w)
}
