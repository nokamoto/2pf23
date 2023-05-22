package main

import (
	"fmt"
	"os"

	"github.com/nokamoto/2pf23/internal/generator/cli"
	"github.com/nokamoto/2pf23/internal/generator/ent"
	"github.com/nokamoto/2pf23/internal/generator/server"
	"github.com/spf13/cobra"
)

func entCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ent input-directory output-directory packge",
		Short:   "Generate ent files from v1.Ent protojson files",
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

func cliCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "cli input-file output-directory root-package",
		Short:        "Generate cli files from the v1.Package protojson file",
		Example:      `cli build/cli/test.json internal/cli/generated github.com/nokamoto/2pf23/internal/cli/generated`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 3 {
				_ = cmd.Usage()
				cmd.PrintErrln()
				return fmt.Errorf("invalid arguments: %v", args)
			}
			walk, err := cli.NewWalk(args[0], args[1], args[2])
			if err != nil {
				return fmt.Errorf("failed to create walk: %w", err)
			}
			if err := walk.Walk(); err != nil {
				return fmt.Errorf("failed to walk: %w", err)
			}
			return nil
		},
	}
	return cmd
}

func serverCommand() *cobra.Command {
	var enableMock bool
	cmd := cobra.Command{
		Use:          "server input-directory output-directory",
		Short:        "Generate server files from v1.Service protojson files",
		Example:      `server build/server github.com/nokamoto/2pf23/internal/service/generated`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				_ = cmd.Usage()
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
	return &cmd
}

func main() {
	cmd := cobra.Command{
		Use: "2pf23-tools",
	}
	cmd.AddCommand(entCommand(), cliCommand(), serverCommand())
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
