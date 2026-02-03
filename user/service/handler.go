package service

import (
	"context"
	"io/fs"

	"github.com/titpetric/platform"

	"github.com/titpetric/platform-app/user/schema"
	"github.com/titpetric/platform-app/user/service/api"
	"github.com/titpetric/platform-app/user/service/web"
	"github.com/titpetric/platform-app/user/storage"
)

// UserModule implements a module contract.
type UserModule struct {
	platform.UnimplementedModule

	templateFS fs.FS

	opts Options
	web  *web.Service
	api  *api.Service
}

// Options is passed from user package scope.
type Options struct {
	SigningKey string
}

// Verify contract.
var _ platform.Module = (*UserModule)(nil)

// NewUserModule sets up dependencies and produces a UserModule.
func NewUserModule(opts Options) *UserModule {
	return &UserModule{
		opts: opts,
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

	h.web = web.NewService(userStorage, sessionStorage)
	h.api = api.NewService(userStorage, api.Options{
		SigningKey: h.opts.SigningKey,
	})
	return nil
}

// Mount registers login, logout, and register routes.
func (h *UserModule) Mount(_ context.Context, r platform.Router) error {
	h.web.Mount(r)
	h.api.Mount(r)
	return nil
}
