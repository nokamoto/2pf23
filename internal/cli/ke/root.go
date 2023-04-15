package ke

import (
	"github.com/nokamoto/2pf23/internal/cli/ke/v1alpha"
	"github.com/nokamoto/2pf23/internal/cli/runtime"
	"github.com/spf13/cobra"
)

func NewKe(rt runtime.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ke",
		Short: "ke is a CLI for managing the kubernetes service.",
		Long:  `ke is a CLI for managing the kubernetes service.`,
	}
	cmd.AddCommand(kev1alpha.NewV1Alpha(rt))
	cmd.AddCommand(kev1alpha.NewV1AlphaSubCommands(rt)...)
	return cmd
}
