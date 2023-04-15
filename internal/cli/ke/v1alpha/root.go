package kev1alpha

import (
	"github.com/nokamoto/2pf23/internal/cli/ke/v1alpha/cluster"
	"github.com/nokamoto/2pf23/internal/cli/runtime"
	"github.com/spf13/cobra"
)

func NewV1Alpha(rt runtime.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "v1alpha",
		Short: "ke v1alpha is a CLI for managing the kubernetes service.",
		Long:  `ke v1alpha is a CLI for managing the kubernetes service.`,
	}
	cmd.AddCommand(NewV1AlphaSubCommands(rt)...)
	return cmd
}

func NewV1AlphaSubCommands(rt runtime.Runtime) []*cobra.Command {
	return []*cobra.Command{kev1alphacluster.NewCluster(rt)}
}
