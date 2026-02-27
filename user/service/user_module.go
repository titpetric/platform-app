package service

import (
	"context"
	"os"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/titpetric/platform"

	"github.com/titpetric/platform-app/user/schema"
	"github.com/titpetric/platform-app/user/service/api"
	"github.com/titpetric/platform-app/user/service/passkey"
	"github.com/titpetric/platform-app/user/service/web"
	"github.com/titpetric/platform-app/user/storage"
)

// Name is the service module name.
const Name = "user"

// UserModule implements a module contract.
type UserModule struct {
	platform.UnimplementedModule

	opts Options
	web  *web.Handlers
	api  *api.Handlers
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
	return Name
}

// Start will initialize the service to handle requests.
func (h *UserModule) Start(ctx context.Context) error {
	db, err := storage.DB(ctx)
	if err != nil {
		return err
	}

	if err := storage.Migrate(ctx, db, schema.Migrations()); err != nil {
		return err
	}

	userStorage := storage.NewUserStorage(db)
	sessionStorage := storage.NewSessionStorage(db)
	passkeyStorage := storage.NewPasskeyStorage(db)

	rpID := os.Getenv("WEBAUTHN_RP_ID")
	if rpID == "" {
		rpID = "localhost"
	}
	rpOrigin := os.Getenv("WEBAUTHN_RP_ORIGIN")
	if rpOrigin == "" {
		rpOrigin = "http://localhost:3000"
	}

	wa, err := webauthn.New(&webauthn.Config{
		RPID:          rpID,
		RPDisplayName: "Platform App",
		RPOrigins:     []string{rpOrigin},
	})
	if err != nil {
		return err
	}

	passkeySvc := passkey.New(wa, passkeyStorage, userStorage)

	h.web = web.NewHandlers(userStorage, sessionStorage, FS(ctx))
	h.api = api.NewHandlers(api.Options{
		SigningKey:     h.opts.SigningKey,
		UserStorage:    userStorage,
		SessionStorage: sessionStorage,
		PasskeyService: passkeySvc,
	})
	return nil
}

// Mount registers login, logout, and register routes.
func (h *UserModule) Mount(_ context.Context, r platform.Router) error {
	h.web.Mount(r)
	h.api.Mount(r)
	return nil
}
