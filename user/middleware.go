package user

import (
	"context"
	"net/http"
	"sync"

	"github.com/titpetric/platform/pkg/telemetry"

	"github.com/titpetric/platform-app/user/model"
	"github.com/titpetric/platform-app/user/service/auth"
	"github.com/titpetric/platform-app/user/storage"
)

// Middleware requires a user login for the registered route.
// If an Authorize header is provided, it will decode the JWT
// and use the user_id claim to authenticate API requests.
//
// How do you generate the JWT you ask? Good question. TBD.
type Middleware struct {
	nextHandler http.Handler

	once    sync.Once
	options model.AuthOptions

	userStorage    *storage.UserStorage
	sessionStorage *storage.SessionStorage
}

func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.init(r.Context())

	if err := m.serveHTTP(w, r); err != nil {
		telemetry.CaptureError(r.Context(), err)
		return
	}

	m.nextHandler.ServeHTTP(w, r)
}

func (m *Middleware) serveHTTP(w http.ResponseWriter, r *http.Request) error {
	if m.options.Header {
		err := m.authorizeJWT(w, r)
		if !m.options.Cookie {
			return err
		}
	}

	if m.options.Cookie {
		return m.authorizeCookie(w, r)
	}
	return ErrLoginRequired
}

func (m *Middleware) authorizeCookie(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie(m.options.CookieName)
	if err != nil || cookie.Value == "" {
		return ErrLoginRequired
	}

	ctx := r.Context()

	session, err := m.sessionStorage.Get(ctx, cookie.Value)
	if err != nil {
		telemetry.CaptureError(ctx, err)
		return ErrLoginRequired
	}

	user, err := m.authorizeUser(w, r, session.UserID)
	if err != nil {
		return err
	}

	if !user.Ok() {
		// Clear cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1,
		})
		return ErrLoginRequired
	}

	sessionIDContext.Set(r, cookie.Value)
	sessionContext.Set(r, session)
	return nil
}

func (m *Middleware) authorizeJWT(w http.ResponseWriter, r *http.Request) error {
	token := r.Header.Get(m.options.HeaderName)
	if token == "" {
		return ErrLoginRequired
	}

	userID, err := auth.NewJWT(m.options.HeaderSigningKey).UserID(token)
	if err != nil {
		return err
	}

	if _, err := m.authorizeUser(w, r, userID); err != nil {
		return err
	}
	return nil
}

func (m *Middleware) authorizeUser(w http.ResponseWriter, r *http.Request, userID string) (*model.User, error) {
	ctx := r.Context()
	user, err := m.userStorage.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	userContext.Set(r, user)
	return user, nil
}

func (m *Middleware) init(ctx context.Context) error {
	var resultErr error
	m.once.Do(func() {
		db, err := storage.DB(ctx)
		if err != nil {
			resultErr = err
			return
		}

		m.userStorage = storage.NewUserStorage(db)
		m.sessionStorage = storage.NewSessionStorage(db)
	})
	return resultErr
}
