package main

import (
	"os"

	"github.com/nokamoto/2pf23/internal/cli"
)

func main() {
	if err := cli.NewCLI().Execute(); err != nil {
		os.Exit(1)
	}
}
