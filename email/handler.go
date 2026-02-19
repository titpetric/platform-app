package email

import (
	"context"

	"github.com/titpetric/platform"

	"github.com/titpetric/platform-app/email/schema"
	"github.com/titpetric/platform-app/email/storage"
)

// Handler implements a module contract.
type Handler struct {
	platform.UnimplementedModule

	service *Service
}

// Verify contract.
var _ platform.Module = (*Handler)(nil)

// NewModule sets up dependencies and produces a handler.
func NewModule() *Handler {
	return &Handler{}
}

// Start will initialize the service to handle email requests.
func (h *Handler) Start(ctx context.Context) error {
	db, err := storage.DB(ctx)
	if err != nil {
		return err
	}

	if err := storage.Migrate(ctx, db, schema.Migrations()); err != nil {
		return err
	}

	emailStorage := storage.NewEmailStorage(db)

	// Pass nil for config to use environment variables automatically
	h.service = NewService(emailStorage, nil)
	h.service.Start()

	return nil
}

// Stop stops the email service when the module shuts down
func (h *Handler) Stop(ctx context.Context) error {
	if h.service != nil {
		h.service.Stop()
	}
	return nil
}

// Name returns the name of the containing package.
func (h *Handler) Name() string {
	return "email"
}

// Mount registers email routes (if any).
func (h *Handler) Mount(_ context.Context, r platform.Router) error {
	return nil
}
