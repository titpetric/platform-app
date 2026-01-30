package service

import (
	"context"

	_ "github.com/titpetric/platform/pkg/drivers"

	"github.com/titpetric/platform"

	"github.com/titpetric/platform-app/pulse/schema"
	"github.com/titpetric/platform-app/pulse/storage"
)

// Name is the service module name.
const Name = "pulse"

type PulseModule struct {
	platform.UnimplementedModule
	Options

	storage *storage.Storage
}

// Options is a placeholder for options.
type Options struct{}

func NewPulseModule(opts Options) *PulseModule {
	return &PulseModule{
		Options: opts,
	}
}

func (p *PulseModule) Name() string {
	return Name
}

func (p *PulseModule) Start(ctx context.Context) error {
	db, err := storage.DB(ctx)
	if err != nil {
		return err
	}

	if err := storage.Migrate(ctx, db, schema.Migrations); err != nil {
		return err
	}

	p.storage = storage.NewStorage(db)
	return nil
}

func (p *PulseModule) Mount(ctx context.Context, r platform.Router) error {
	handlers := NewHandlers(p.storage)
	handlers.Mount(r)
	return nil
}
