// Package server implements the pulse HTTP server command.
package server

import (
	"context"
	"fmt"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/titpetric/cli"
	"github.com/titpetric/platform"

	"github.com/titpetric/platform-app/pulse"
	"github.com/titpetric/platform-app/pulse/view"
	"github.com/titpetric/platform-app/user"
)

// Name is the command title.
const Name = "Run the server process"

// Options holds server command configuration.
type Options struct{}

// Bind registers server flags with the flag set.
func (o *Options) Bind(flag *cli.FlagSet) {
	// Server has no CLI options
}

// NewCommand creates a new server command.
func NewCommand() *cli.Command {
	var opts Options

	return &cli.Command{
		Name:  "server",
		Title: Name,
		Bind: func(flag *cli.FlagSet) {
			opts.Bind(flag)
		},
		Run: func(ctx context.Context, args []string) error {
			return Run(ctx, opts)
		},
	}
}

// Run starts the pulse HTTP server.
func Run(ctx context.Context, opts Options) error {
	platformOpts := platform.NewOptions()
	platformOpts.ThemeFS = view.FS

	svc := platform.New(platformOpts)

	svc.Use(middleware.Logger)
	svc.Register(user.NewModule())
	svc.Register(pulse.NewModule())

	if err := svc.Start(ctx); err != nil {
		return fmt.Errorf("exit error: %w", err)
	}

	svc.Wait()
	return nil
}
