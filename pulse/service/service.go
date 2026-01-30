package service

import (
	"context"

	_ "github.com/titpetric/platform/pkg/drivers"

	"github.com/titpetric/platform"
)

// Name is the service module name.
const Name = "pulse"

type PulseModule struct {
	platform.UnimplementedModule
	Options
}

type Options struct {
	Path string
}

func NewPulseModule(opts Options) *PulseModule {
	return &PulseModule{
		Options: opts,
	}
}

func (p *PulseModule) Name() string {
	return Name
}

func (p *PulseModule) Mount(ctx context.Context, r platform.Router) error {
	handlers := NewHandlers(p.Options)
	handlers.Mount(r)
	return nil
}
