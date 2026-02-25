package web

import (
	"net/http"

	"github.com/titpetric/platform/pkg/telemetry"

	"github.com/titpetric/platform-app/user/model"
)

// Login handles user authentication via HTML form submission.
func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.login(w, r))
}

func (h *Handlers) login(w http.ResponseWriter, r *http.Request) error {
	r, span := telemetry.StartRequest(r, "user.service.Login")
	defer span.End()

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		h.Error(r, "Email and Password are required", nil)
		h.LoginView(w, r)
		return nil
	}

	user, err := h.userStorage.Authenticate(r.Context(), model.UserAuth{
		Email:    email,
		Password: password,
	})
	if err != nil || !user.Ok() {
		h.Error(r, "Invalid credentials for login", err)
		h.LoginView(w, r)
		return nil
	}

	session, err := h.sessionStorage.Create(r.Context(), user.ID)
	if err != nil {
		h.Error(r, "Can't create session", err)
		h.LoginView(w, r)
		return nil
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  *session.ExpiresAt,
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return nil
}
