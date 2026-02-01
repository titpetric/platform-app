package main

import (
	"fmt"
	"os"

	"github.com/titpetric/cli"
)

func main() {
	app := cli.NewApp("example")
	app.AddCommand("server", "Run the server process", cmd)
	app.DefaultCommand = "server"

	if err := app.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
