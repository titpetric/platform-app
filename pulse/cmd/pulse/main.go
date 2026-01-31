package main

import (
	"fmt"
	"os"

	"github.com/titpetric/cli"
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
