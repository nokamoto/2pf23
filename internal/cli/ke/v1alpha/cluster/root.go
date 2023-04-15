package kev1alphacluster

import (
	"github.com/nokamoto/2pf23/internal/cli/runtime"
	"github.com/spf13/cobra"
)

func NewCluster(rt runtime.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "cluster is a CLI for managing the kubernetes cluster.",
		Long:  `cluster is a CLI for managing the kubernetes cluster.`,
	}
	cmd.AddCommand(newCreate(rt))
	return cmd
}
