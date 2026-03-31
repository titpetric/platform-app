package user

import (
	"context"
	"net/http"
	"os"

	"github.com/titpetric/platform-app/user/model"
	"github.com/titpetric/platform-app/user/service"
)

// NewModule will return the user module.
func NewModule() *service.UserModule {
	return service.NewUserModule(service.Options{
		SigningKey: SigningKey(),
	})
}

// MiddlewareOption configures the user authentication middleware.
type MiddlewareOption func(*Middleware)

// NewMiddleware will populate context for IsLoggedIn, GetSessionUser.
func NewMiddleware(opts ...MiddlewareOption) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		mw := &Middleware{
			nextHandler: next,
		}
		for _, opt := range opts {
			opt(mw)
		}
		return mw
	}
}

// AuthHeader enables JWT-based authentication via the Authorization header.
func AuthHeader() MiddlewareOption {
	return func(mw *Middleware) {
		mw.options.Header = true
		mw.options.HeaderName = "Authorization"
		mw.options.HeaderSigningKey = SigningKey()
	}
}

func SigningKey() string {
	if key := os.Getenv("USER_JWT_SIGNING_KEY"); key != "" {
		return key
	}
	return "test-usage"
}

// AuthCookie enables session-based authentication via a cookie.
func AuthCookie() MiddlewareOption {
	return func(mw *Middleware) {
		mw.options.Cookie = true
		mw.options.CookieName = os.Getenv("USER_SESSION_COOKIE_NAME")
		if mw.options.CookieName == "" {
			mw.options.CookieName = "session_id"
		}
	}
}

// AuthOptional makes the middleware non-blocking on auth failure.
// The request proceeds even if no valid session is found.
func AuthOptional() MiddlewareOption {
	return func(mw *Middleware) {
		mw.options.Optional = true
	}
}

// AuthQuery enables JWT-based authentication via URL query parameter.
// Useful for headless browsers that can't set Authorization headers.
func AuthQuery(paramName string) MiddlewareOption {
	return func(mw *Middleware) {
		mw.options.Query = true
		mw.options.QueryName = paramName
		mw.options.HeaderSigningKey = SigningKey()
	}
}

// GetSessionUser will return the user bound to the session.
// If there's no active user bound, the return is nil, false.
func GetSessionUser(ctx context.Context) (*model.User, bool) {
	userdata := userContext.GetContext(ctx)
	return userdata, userdata.Ok()
}

// IsLoggedIn returns true if there's an active user session.
func IsLoggedIn(ctx context.Context) bool {
	_, ok := GetSessionUser(ctx)
	return ok
}

// SetSessionUser is here to aid testing, for internal use.
func SetSessionUser(ctx context.Context, u *model.User) context.Context {
	return userContext.SetContext(ctx, u)
}
