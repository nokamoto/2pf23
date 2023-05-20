// Code generated by cli-gen. DO NOT EDIT.
package cluster

import (
	"fmt"
	"strings"
)

import (
	"github.com/bufbuild/connect-go"
	helper "github.com/nokamoto/2pf23/internal/cli/helper"
	"github.com/nokamoto/2pf23/internal/cli/runtime"
	helperapi "github.com/nokamoto/2pf23/internal/util/helper"
	v1alpha "github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

func newCreateCluster(rt runtime.Runtime) *cobra.Command {
	var displayName string
	var numNodes int32
	machineType := helper.NewEnumFlag[v1alpha.MachineType](v1alpha.MachineType_name, v1alpha.MachineType_value)
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create is a command to create a new Cluster",
		Long:  `create is a command to create a new Cluster`,
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			ctx := rt.Context(cmd)
			c, err := rt.KeV1alpha(cmd)
			if err != nil {
				return fmt.Errorf("failed to create a client for ke.v1alpha: %w", err)
			}
			res, err := c.CreateCluster(ctx, connect.NewRequest(&v1alpha.CreateClusterRequest{
				Cluster: &v1alpha.Cluster{
					DisplayName: displayName,
					NumNodes:    numNodes,
					MachineType: machineType.Value(),
				},
			}))
			if err != nil {
				return fmt.Errorf("ke.v1alpha: failed to CreateCluster: %w", err)
			}
			message, err := protojson.Marshal(helperapi.GetResponseMsg(res))
			if err != nil {
				return err
			}
			cmd.Println(string(message))
			return nil
		},
	}
	cmd.Flags().StringVar(&displayName, "display-name", "", "The display name of the cluster.")
	cmd.Flags().Int32Var(&numNodes, "num-nodes", 0, "")
	cmd.Flags().Var(machineType, "machine-type", fmt.Sprintf(" [%s]", strings.Join(machineType.Names(), ", ")))
	cmd.RegisterFlagCompletionFunc("machine-type", machineType.CompletionFunc())
	return cmd
}
