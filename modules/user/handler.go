package user

import (
	"context"

	"github.com/titpetric/platform"

	"github.com/titpetric/platform-app/modules/theme"
	"github.com/titpetric/platform-app/modules/user/service"
	"github.com/titpetric/platform-app/modules/user/storage"
)

// Handler implements a module contract.
type Handler struct {
	Service *service.Service
}

// Verify contract.
var _ platform.Module = (*Handler)(nil)

// NewHandler sets up dependencies and produces a handler.
func NewHandler() *Handler {
	return &Handler{}
}

// Start will initialize the service to handle requests.
func (h *Handler) Start() error {
	db, err := storage.DB(context.Background())
	if err != nil {
		return err
	}

	themeFS := theme.TemplateFS
	userStorage := storage.NewUserStorage(db)
	sessionStorage := storage.NewSessionStorage(db)

	options := &service.Options{
		UserStorage:    userStorage,
		SessionStorage: sessionStorage,
		ThemeFS:        themeFS,
		ModuleFS:       TemplateFS,
	}

	svc, err := service.NewService(options)
	if err != nil {
		return err
	}

	h.Service = svc
	return nil
}

// Name returns the name of the containing package.
func (h *Handler) Name() string {
	return "user"
}

// Mount registers login, logout, and register routes.
func (h *Handler) Mount(r platform.Router) error {
	r.Get("/login", h.Service.LoginView)
	r.Post("/login", h.Service.Login)
	r.Get("/logout", h.Service.LogoutView)
	r.Post("/logout", h.Service.Logout)
	r.Get("/register", h.Service.RegisterView)
	r.Post("/register", h.Service.Register)
	return nil
}

// Stop implements a closer for graceful shutdown.
func (h *Handler) Stop() error {
	h.Service.Close()
	return nil
}
