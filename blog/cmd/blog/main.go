package main

import (
	"fmt"
	"os"
	"os/exec"
	"slices"

	_ "github.com/titpetric/platform/pkg/drivers"

	"github.com/titpetric/cli"
	"github.com/titpetric/platform"

	"github.com/titpetric/platform-app/blog/cmd/blog/server"
	"github.com/titpetric/platform-app/blog/cmd/blog/version"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// ExtendedProvider extends database providers with listing and registration.
type ExtendedProvider interface {
	List() []string
	Register(string, string)
}

func run() error {
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git is not installed or not found in PATH: the blog module requires git for content management")
	}

	// Add default storage for blog.
	if val, ok := platform.Database.(ExtendedProvider); ok {
		connectionList := val.List()
		if !slices.Contains(connectionList, "user") {
			val.Register("user", "sqlite://user.db")
		}
		if !slices.Contains(connectionList, "blog") {
			val.Register("blog", "sqlite://blog.db")
		}
	}

	app := cli.NewApp("blog")
	app.AddCommand("server", server.Name, server.NewCommand)
	app.AddCommand("version", version.Name, func() *cli.Command {
		return version.NewCommand(version.Info{
			Version:    Version,
			Commit:     Commit,
			CommitTime: CommitTime,
			Branch:     Branch,
		})
	})

	app.DefaultCommand = "server"

	return app.Run()
}
