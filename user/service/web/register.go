package web

import (
	"net/http"

	"github.com/titpetric/platform/pkg/telemetry"

	"github.com/titpetric/platform-app/user/model"
)

// Register handles creating a new user and starting a session via HTML form submission.
func (h *Service) Register(w http.ResponseWriter, r *http.Request) {
	r, span := telemetry.StartRequest(r, "user.service.Register")
	defer span.End()

	ctx := r.Context()

	req := &model.UserCreateRequest{
		FullName: r.FormValue("full_name"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	if !req.Valid() {
		h.Error(r, "All fields are required", nil)
		return
	}

	createdUser, err := h.userStorage.Create(ctx, req)
	if err != nil {
		h.Error(r, "Failed to create user", err)
		h.RegisterView(w, r)
		return
	}

	session, err := h.sessionStorage.Create(ctx, createdUser.ID)
	if err != nil {
		h.Error(r, "Failed to create session", err)
		h.RegisterView(w, r)
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Expires:  *session.ExpiresAt,
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
