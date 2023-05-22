package main

import (
	"fmt"
	"os"

	"github.com/nokamoto/2pf23/internal/generator/ent"
	"github.com/spf13/cobra"
)

func entCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ent input-directory output-directory packge",
		Example: "ent build/ent internal/ent/proto",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 3 {
				_ = cmd.Usage()
				cmd.PrintErrln()
				return fmt.Errorf("invalid arguments: %v", args)
			}
			w, err := ent.NewWalk(args[0], args[1], args[2])
			if err != nil {
				return fmt.Errorf("failed to create walk: %w", err)
			}
			if err := w.Walk(); err != nil {
				return fmt.Errorf("failed to walk: %w", err)
			}
			return nil
		},
	}
	return cmd
}

func main() {
	cmd := cobra.Command{
		Use:          "2pf23-tools",
		SilenceUsage: true,
	}
	cmd.AddCommand(entCommand())
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
