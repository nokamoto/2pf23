package cli

import (
	"github.com/nokamoto/2pf23/internal/cli/ke"
	"github.com/nokamoto/2pf23/internal/cli/runtime"
	"github.com/spf13/cobra"
)

func NewCLI() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "pf",
		Long: `pf is a CLI for managing the platform.`,
	}
	rt := runtime.NewRuntime()
	cmd.AddCommand(ke.NewKe(rt))
	return cmd
}
