package service

import (
	"github.com/titpetric/platform"

	"github.com/titpetric/platform-app/user/storage"
)

// Service encapsulates what we need to get from the handler.
type Service struct {
	UserStorage    *storage.UserStorage
	SessionStorage *storage.SessionStorage

	view *Renderer
}

// NewService takes in required dependencies to support the MVC framework.
func NewService(u *storage.UserStorage, s *storage.SessionStorage) *Service {
	svc := &Service{
		UserStorage:    u,
		SessionStorage: s,
		view:           NewRenderer(nil),
	}
	return svc
}

// Mount registers login, logout, and register routes.
func (s *Service) Mount(r platform.Router) {
	r.Get("/login", s.LoginView)
	r.Post("/login", s.Login)
	r.Get("/logout", s.LogoutView)
	r.Post("/logout", s.Logout)
	r.Get("/register", s.RegisterView)
	r.Post("/register", s.Register)
}
