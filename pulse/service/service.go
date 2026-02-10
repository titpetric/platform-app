// Package service implements the pulse HTTP service module.
package service

import (
	"context"

	"github.com/titpetric/platform"

	"github.com/titpetric/platform-app/pulse/schema"
	"github.com/titpetric/platform-app/pulse/storage"
	userstorage "github.com/titpetric/platform-app/user/storage"
)

// Name is the service module name.
const Name = "pulse"

// PulseModule is the pulse platform module.
type PulseModule struct {
	platform.UnimplementedModule

	storage *storage.Storage

	userStorage *userstorage.UserStorage
}

// NewPulseModule creates a new pulse module.
func NewPulseModule() *PulseModule {
	return &PulseModule{}
}

// Name returns the module name.
func (p *PulseModule) Name() string {
	return Name
}

func (p *PulseModule) setupStorage(ctx context.Context) error {
	if err := p.setupUserStorage(ctx); err != nil {
		return err
	}
	if err := p.setupPulseStorage(ctx); err != nil {
		return err
	}
	return nil
}

func (p *PulseModule) setupUserStorage(ctx context.Context) error {
	db, err := userstorage.DB(ctx)
	if err != nil {
		return err
	}

	p.userStorage = userstorage.NewUserStorage(db)
	return nil
}

func (p *PulseModule) setupPulseStorage(ctx context.Context) error {
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

// Start initializes module storage.
func (p *PulseModule) Start(ctx context.Context) error {
	if err := p.setupStorage(ctx); err != nil {
		return err
	}
	return nil
}

// Mount registers module HTTP handlers.
func (p *PulseModule) Mount(ctx context.Context, r platform.Router) error {
	handlers := NewHandlers(p.storage)
	handlers.Mount(r)
	return nil
}
