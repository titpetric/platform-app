package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/titpetric/cli"
	"github.com/titpetric/platform"

	"github.com/titpetric/platform-app/pulse"
	"github.com/titpetric/platform-app/user"
)

func main() {
	app := cli.NewApp("example")

	app.AddCommand("server", "Run the server process", cmdServer)
	app.AddCommand("client", "Run the client process", cmdClient)

	if err := app.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// cmdServer constructs the "server" command.
func cmdServer() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Title: "Run the server process",
		Run: func(ctx context.Context, args []string) error {
			opts := platform.NewOptions()
			svc := platform.New(opts)

			svc.Use(middleware.Logger)
			svc.Register(user.NewHandler())
			svc.Register(pulse.NewModule())

			if err := svc.Start(ctx); err != nil {
				return err
			}

			svc.Wait()

			return nil
		},
	}
}

// cmdClient constructs the "client" command.
func cmdClient() *cli.Command {
	var (
		addr    string
		timeout time.Duration
	)

	return &cli.Command{
		Name:  "client",
		Title: "Run the client process",
		Bind: func(fs *cli.FlagSet) {
			cli.StringVarP(&addr, "addr", "a", "http://localhost:8080", "server address")
			cli.DurationVar(&timeout, "timeout", 5*time.Second, "request timeout")
		},
		Run: func(context.Context, []string) error {
			// implementation intentionally omitted
			return nil
		},
	}
}
