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
	return service.NewUserModule(service.Options{})
}

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

func AuthHeader() MiddlewareOption {
	return func(mw *Middleware) {
		mw.options.Header = true
		mw.options.HeaderName = "Authorization"
		mw.options.HeaderSigningKey = os.Getenv("USER_JWT_SIGNING_KEY")
		if mw.options.HeaderSigningKey == "" {
			mw.options.HeaderSigningKey = "test-usage"
		}
	}
}

func AuthCookie() MiddlewareOption {
	return func(mw *Middleware) {
		mw.options.Cookie = true
		mw.options.CookieName = os.Getenv("USER_SESSION_COOKIE_NAME")
		if mw.options.CookieName == "" {
			mw.options.CookieName = "session_id"
		}
	}
}

// GetSessionUser will return the user bound to the session.
// If there's no active user bound, the return is nil, false.
func GetSessionUser(ctx context.Context) (*model.User, bool) {
	userdata := userContext.GetContext(ctx)
	return userdata, userdata.Ok()
}

// SetSessionUser is here to aid testing, for internal use.
func SetSessionUser(ctx context.Context, u *model.User) context.Context {
	return userContext.SetContext(ctx, u)
}
