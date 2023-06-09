package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nokamoto/2pf23/internal/protogen/server"
)

func main() {
	if err := protogen.NewPlugin().Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", filepath.Base(os.Args[0]), err)
		os.Exit(1)
	}
}
