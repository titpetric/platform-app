package main

import (
	"fmt"
	"os"

	_ "github.com/titpetric/platform/pkg/drivers"

	"github.com/titpetric/cli"

	"github.com/titpetric/platform-app/pulse/cmd/pulse/version"
)

func main() {
	app := cli.NewApp("example")
	app.AddCommand("server", "Run the server process", cmd)
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
