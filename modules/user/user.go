package user

import (
	"net/http"

	"github.com/titpetric/platform/pkg/httpcontext"
	"github.com/titpetric/platform/pkg/telemetry"

	"github.com/titpetric/platform-app/modules/user/model"
	"github.com/titpetric/platform-app/modules/user/storage"
)

type userKey struct{}

var userContext = httpcontext.NewValue[*model.User](userKey{})

// IsLoggedIn returns true or false. Any errors are swallowed, returning false.
func IsLoggedIn(r *http.Request) bool {
	r, span := telemetry.StartRequest(r, "user.IsLoggedIn")
	defer span.End()

	if user := userContext.Get(r); user != nil {
		return user.IsActive()
	}

	user, err := GetSessionUser(r)
	if user == nil || err != nil {
		return false
	}

	return user.IsActive()
}

// GetSessionUser will return the user bound to the session. If no user is bound to
// the session or there is no session, the function will return nil, nil.
func GetSessionUser(r *http.Request) (*model.User, error) {
	r, span := telemetry.StartRequest(r, "user.GetSessionUser")
	defer span.End()

	ctx := r.Context()

	if user := userContext.GetContext(ctx); user != nil {
		return user, nil
	}

	db, err := storage.DB(ctx)
	if err != nil {
		return nil, err
	}

	cookie, err := r.Cookie("session_id")
	if err == nil && cookie.Value != "" {
		userStorage := storage.NewUserStorage(db)
		sessionStorage := storage.NewSessionStorage(db)

		session, err := sessionStorage.Get(ctx, cookie.Value)
		if err != nil {
			return nil, err
		}

		user, err := userStorage.Get(ctx, session.UserID)
		if err != nil {
			return nil, err
		}

		userContext.Set(r, user)

		return user, nil
	}

	return nil, nil
}
