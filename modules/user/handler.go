package user

import (
	"context"

	"github.com/titpetric/platform"

	"github.com/titpetric/platform-app/modules/user/service"
	"github.com/titpetric/platform-app/modules/user/storage"
)

// Handler implements a module contract.
type Handler struct {
	platform.UnimplementedModule

	*service.Service
}

// Verify contract.
var _ platform.Module = (*Handler)(nil)

// NewHandler sets up dependencies and produces a handler.
func NewHandler() *Handler {
	return &Handler{}
}

// Start will initialize the service to handle requests.
func (h *Handler) Start(ctx context.Context) error {
	db, err := storage.DB(ctx)
	if err != nil {
		return err
	}

	userStorage := storage.NewUserStorage(db)
	sessionStorage := storage.NewSessionStorage(db)

	h.Service = service.NewService(userStorage, sessionStorage)
	return nil
}

// Name returns the name of the containing package.
func (h *Handler) Name() string {
	return "user"
}

// Mount registers login, logout, and register routes.
func (h *Handler) Mount(r platform.Router) error {
	r.Get("/login", h.LoginView)
	r.Post("/login", h.Login)
	r.Get("/logout", h.LogoutView)
	r.Post("/logout", h.Logout)
	r.Get("/register", h.RegisterView)
	r.Post("/register", h.Register)
	return nil
}
