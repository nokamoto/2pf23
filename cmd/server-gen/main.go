package main

import (
	"fmt"
	"os"

	"github.com/nokamoto/2pf23/internal/servergen"
	"github.com/spf13/cobra"
)

func main() {
	var enableMock bool
	cmd := cobra.Command{
		Use:          "server-gen input-directory output-directory",
		Short:        "server-gen is a server generator.",
		Long:         `server-gen is a server generator.`,
		Example:      `server-gen build/server github.com/nokamoto/2pf23/internal/service/generated`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				cmd.Usage()
				cmd.PrintErrln()
				return fmt.Errorf("invalid arguments: %v", args)
			}
			walk, err := servergen.NewWalk(args[0], args[1], enableMock)
			if err != nil {
				return fmt.Errorf("failed to create walk: %w", err)
			}
			if err := walk.Walk(); err != nil {
				return fmt.Errorf("failed to walk: %w", err)
			}
			return nil
		},
	}
	cmd.PersistentFlags().BoolVar(&enableMock, "mock", false, "enable mock generation")
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
