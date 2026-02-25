package web

import (
	"io/fs"

	"github.com/titpetric/platform"

	"github.com/titpetric/platform-app/user/storage"
)

// Handlers encapsulates what we need to get from the handler.
type Handlers struct {
	userStorage    *storage.UserStorage
	sessionStorage *storage.SessionStorage

	view *Renderer
}

// NewHandlers takes in required dependencies to support the MVC framework.
// Context should be passed from Start() to access platform options.
func NewHandlers(u *storage.UserStorage, s *storage.SessionStorage, viewFS fs.FS) *Handlers {
	svc := &Handlers{
		userStorage:    u,
		sessionStorage: s,
		view:           NewRenderer(viewFS, nil),
	}
	return svc
}

// Mount registers login, logout, and register routes.
func (s *Handlers) Mount(r platform.Router) {
	r.Get("/login", s.LoginView)
	r.Post("/login", s.Login)
	r.Get("/logout", s.LogoutView)
	r.Post("/logout", s.Logout)
	r.Get("/register", s.RegisterView)
	r.Post("/register", s.Register)
}
