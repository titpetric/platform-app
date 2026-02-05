package main

import (
	"fmt"
	"os"
	"slices"

	_ "github.com/titpetric/platform/pkg/drivers"

	"github.com/titpetric/cli"
	"github.com/titpetric/platform"

	"github.com/titpetric/platform-app/pulse/cmd/pulse/login"
	"github.com/titpetric/platform-app/pulse/cmd/pulse/register"
	"github.com/titpetric/platform-app/pulse/cmd/pulse/version"
)

type ExtendedProvider interface {
	List() []string
	Register(string, string)
}

func main() {
	// Add default storage for pulse.
	if val, ok := platform.Database.(ExtendedProvider); ok {
		connectionList := val.List()
		if !slices.Contains(connectionList, "user") {
			val.Register("user", "sqlite://user.db")
		}
		if !slices.Contains(connectionList, "pulse") {
			val.Register("pulse", "sqlite://pulse.db")
		}
	}

	app := cli.NewApp("example")
	app.AddCommand("server", "Run the server process", cmd)
	app.AddCommand("login", login.Name, login.NewCommand)
	app.AddCommand("register", register.Name, register.NewCommand)
	app.AddCommand("version", version.Name, func() *cli.Command {
		return version.NewCommand(version.Info{
			Version:    Version,
			Commit:     Commit,
			CommitTime: CommitTime,
			Branch:     Branch,
		})
	})

	app.DefaultCommand = "server"

	if err := app.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
