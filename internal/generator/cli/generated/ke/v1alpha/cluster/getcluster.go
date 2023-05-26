// Code generated by cli-gen. DO NOT EDIT.
package cluster

import (
	"fmt"
)

import (
	"github.com/bufbuild/connect-go"
	"github.com/nokamoto/2pf23/internal/cli/runtime"
	helperapi "github.com/nokamoto/2pf23/internal/util/helper"
	v1alpha "github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

func newGetCluster(rt runtime.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get cluster-name",
		Short: "get is a command to get the Cluster",
		Long:  `get is a command to get the Cluster`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			ctx := rt.Context(cmd)
			c, err := rt.KeV1alpha(cmd)
			if err != nil {
				return fmt.Errorf("failed to create a client for ke.v1alpha: %w", err)
			}
			res, err := c.GetCluster(ctx, connect.NewRequest(&v1alpha.GetClusterRequest{
				Name: args[0],
			}))
			if err != nil {
				return fmt.Errorf("ke.v1alpha: failed to GetCluster: %w", err)
			}
			message, err := protojson.Marshal(helperapi.GetResponseMsg(res))
			if err != nil {
				return err
			}
			cmd.Println(string(message))
			return nil
		},
	}
	return cmd
}