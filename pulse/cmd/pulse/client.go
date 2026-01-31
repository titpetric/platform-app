package main

import (
	"context"
	"time"

	"github.com/titpetric/cli"
)

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
