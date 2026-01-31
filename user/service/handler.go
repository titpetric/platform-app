package service

import (
	"context"
	"io/fs"

	"github.com/titpetric/platform"

	"github.com/titpetric/platform-app/user/schema"
	"github.com/titpetric/platform-app/user/storage"
)

// UserModule implements a module contract.
type UserModule struct {
	platform.UnimplementedModule

	templateFS fs.FS
	svc        *Service
}

// Options is passed from user package scope to service.
type Options struct {
	TemplateFS fs.FS
}

// Verify contract.
var _ platform.Module = (*UserModule)(nil)

// NewUserModule sets up dependencies and produces a UserModule.
func NewUserModule(opts Options) *UserModule {
	return &UserModule{
		templateFS: opts.TemplateFS,
	}
}

// Name returns the name of the containing package.
func (h *UserModule) Name() string {
	return "user"
}

// Start will initialize the service to handle requests.
func (h *UserModule) Start(ctx context.Context) error {
	db, err := storage.DB(ctx)
	if err != nil {
		return err
	}

	if err := storage.Migrate(ctx, db, schema.Migrations); err != nil {
		return err
	}

	userStorage := storage.NewUserStorage(db)
	sessionStorage := storage.NewSessionStorage(db)

	h.svc = NewService(h.templateFS, userStorage, sessionStorage)
	return nil
}

// Mount registers login, logout, and register routes.
func (h *UserModule) Mount(_ context.Context, r platform.Router) error {
	h.svc.Mount(r)
	return nil
}
