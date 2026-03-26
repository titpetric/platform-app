// Package server implements the blog HTTP server command.
package server

import (
	"context"
	"fmt"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/titpetric/cli"
	"github.com/titpetric/platform"
	"github.com/titpetric/vuego"

	"github.com/titpetric/platform-app/blog"
	"github.com/titpetric/platform-app/blog/config"
	"github.com/titpetric/platform-app/user"
)

// Name is the command title.
const Name = "Run the server process"

// NewCommand creates a new server command.
func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Title: Name,
		Run: func(ctx context.Context, args []string) error {
			return Run(ctx)
		},
	}
}

// Run starts the blog HTTP server.
func Run(ctx context.Context) error {
	platformOpts := platform.NewOptions()
	platformOpts.ConfigFS = vuego.NewOverlayFS(config.ConfigFS())

	svc := platform.New(platformOpts)

	svc.Use(middleware.Logger)
	svc.Register(user.NewModule())
	svc.Register(blog.NewModule())

	if err := svc.Start(ctx); err != nil {
		return fmt.Errorf("exit error: %w", err)
	}

	svc.Wait()
	return nil
}
