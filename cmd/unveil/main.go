package main

import (
	"fmt"
	"os"

	"github.com/gi8lino/unveil/internal/app"
)

var (
	Version = "dev"
	Commit  = "none"
)

// main sets up the application context and runs the main loop.
func main() {
	if err := app.Run(
		Version,
		Commit,
		os.Args[1:],
		os.Stdout,
	); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
