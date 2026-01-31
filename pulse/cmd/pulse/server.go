package main

import (
	"context"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/titpetric/cli"
	"github.com/titpetric/platform"

	"github.com/titpetric/platform-app/pulse"
	"github.com/titpetric/platform-app/user"
)

// cmdServer constructs the "server" command.
func cmdServer() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Title: "Run the server process",
		Run: func(ctx context.Context, args []string) error {
			opts := platform.NewOptions()
			svc := platform.New(opts)

			svc.Use(middleware.Logger)
			svc.Register(user.NewModule())
			svc.Register(pulse.NewModule())

			if err := svc.Start(ctx); err != nil {
				return err
			}

			svc.Wait()

			return nil
		},
	}
}
