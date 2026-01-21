package user

import (
	"context"
	"net/http"

	"github.com/titpetric/platform/pkg/telemetry"

	"github.com/titpetric/platform-app/user/model"
	"github.com/titpetric/platform-app/user/storage"
)

// Middleware will populate context for IsLoggedIn, GetSessionUser.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil || cookie.Value == "" {
			return
		}

		ctx := r.Context()

		db, err := storage.DB(ctx)

		userStorage := storage.NewUserStorage(db)
		sessionStorage := storage.NewSessionStorage(db)

		session, err := sessionStorage.Get(ctx, cookie.Value)
		if err != nil {
			telemetry.CaptureError(ctx, err)
			return
		}

		user, err := userStorage.Get(ctx, session.UserID)
		if err != nil {
			telemetry.CaptureError(ctx, err)
			return
		}
		if !user.IsActive() {
			// Clear cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "session_id",
				Value:    "",
				Path:     "/",
				HttpOnly: true,
				MaxAge:   -1,
			})
			return
		}

		sessionIDContext.Set(r, cookie.Value)
		sessionContext.Set(r, session)
		userContext.Set(r, user)

		next.ServeHTTP(w, r)
	})
}

// GetSessionUser will return the user bound to the session.
// If there's no active user bound, the return is nil, false.
func GetSessionUser(ctx context.Context) (*model.User, bool) {
	userdata := userContext.GetContext(ctx)
	return userdata, userdata.IsActive()
}

// SetSessionUser is here to aid testing.
func SetSessionUser(ctx context.Context, u *model.User) context.Context {
	return userContext.SetContext(ctx, u)
}
