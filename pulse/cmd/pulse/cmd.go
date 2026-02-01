package main

import (
	"context"
	"fmt"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/titpetric/cli"
	"github.com/titpetric/platform"

	"github.com/titpetric/platform-app/pulse"
	"github.com/titpetric/platform-app/pulse/service"
	"github.com/titpetric/platform-app/user"
)

// cmd constructs the "pulse" command.
func cmd() *cli.Command {
	var opts service.Options

	return &cli.Command{
		Name:  "server",
		Title: "Run the server process",
		Bind: func(flag *cli.FlagSet) {
			opts.Bind(flag)
		},
		Run: func(ctx context.Context, args []string) error {
			svc := platform.New(platform.NewOptions())

			svc.Use(middleware.Logger)
			svc.Register(user.NewModule())
			svc.Register(pulse.NewModule(opts))

			if err := svc.Start(ctx); err != nil {
				return fmt.Errorf("exit error: %w", err)
			}

			svc.Wait()
			return nil
		},
	}
}
