package main

import (
	"fmt"
	"os"

	"github.com/nokamoto/2pf23/internal/cligen"
	"github.com/spf13/cobra"
)

func main() {
	cmd := cobra.Command{
		Use:          "cli-gen input-file output-directory root-package",
		Short:        "cli-gen is a command line interface generator.",
		Long:         `cli-gen is a command line interface generator.`,
		Example:      `cli-gen build/cli/test.json internal/cli/generated github.com/nokamoto/2pf23/internal/cli/generated`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 3 {
				_ = cmd.Usage()
				cmd.PrintErrln()
				return fmt.Errorf("invalid arguments: %v", args)
			}
			walk, err := cligen.NewWalk(args[0], args[1], args[2])
			if err != nil {
				return fmt.Errorf("failed to create walk: %w", err)
			}
			if err := walk.Walk(); err != nil {
				return fmt.Errorf("failed to walk: %w", err)
			}
			return nil
		},
	}
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
