// Code generated by cli-gen. DO NOT EDIT.
package cluster

import (
	"fmt"
)

import (
	"github.com/nokamoto/2pf23/internal/cli/runtime"
	v1alpha "github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func newUpdateCluster(rt runtime.Runtime) *cobra.Command {
	var displayName string
	var numNodes int32
	var machineType v1alpha.MachineType
	cmd := &cobra.Command{
		Use:   "update cluster-name",
		Short: "update is a command to update the Cluster",
		Long:  `update is a command to update the Cluster`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			ctx := rt.Context(cmd)
			c, err := rt.KeV1alpha(cmd)
			if err != nil {
				return fmt.Errorf("failed to create a client for ke.v1alpha: %w", err)
			}
			var paths []string
			if cmd.Flags().Changed("display-name") {
				paths = append(paths, "display_name")
			}
			if cmd.Flags().Changed("num-nodes") {
				paths = append(paths, "num_nodes")
			}
			mask, err := fieldmaskpb.New(&v1alpha.Cluster{}, paths...)
			if err != nil {
				return fmt.Errorf("failed to create a field mask: %w", err)
			}
			res, err := c.UpdateCluster(ctx, &v1alpha.UpdateClusterRequest{
				UpdateMask: mask,
				Cluster: &v1alpha.Cluster{
					Name:        args[0],
					DisplayName: displayName,
					NumNodes:    numNodes,
					MachineType: machineType,
				},
			})
			if err != nil {
				return fmt.Errorf("ke.v1alpha: failed to UpdateCluster: %w", err)
			}
			message, err := protojson.Marshal(res)
			if err != nil {
				return err
			}
			cmd.Println(string(message))
			return nil
		},
	}
	cmd.Flags().StringVar(&displayName, "display-name", "", "The display name of the cluster.")
	cmd.Flags().Int32Var(&numNodes, "num-nodes", 0, "The number of worker nodes in the cluster.")
	return cmd
}
