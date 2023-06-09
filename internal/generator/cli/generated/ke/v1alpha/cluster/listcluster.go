// Code generated by cli-gen. DO NOT EDIT.
package cluster

import (
	"fmt"
)

import (
	helper "github.com/nokamoto/2pf23/internal/cli/helper"
	"github.com/nokamoto/2pf23/internal/cli/runtime"
	v1alpha "github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

func newListCluster(rt runtime.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list is a command to list all Clusters",
		Long:  `list is a command to list all Cluster`,
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			ctx := rt.Context(cmd)
			c, err := rt.KeV1alpha(cmd)
			if err != nil {
				return fmt.Errorf("failed to create a client for ke.v1alpha: %w", err)
			}
			var res v1alpha.ListClusterResponse
			setter := func(token string) *v1alpha.ListClusterRequest {
				return &v1alpha.ListClusterRequest{
					PageToken: token,
				}
			}
			getter := func(v *v1alpha.ListClusterResponse) {
				res.Clusters = append(res.Clusters, v.Clusters...)
			}
			err = helper.ListAll(ctx, c.ListCluster, setter, getter, (*v1alpha.ListClusterResponse).GetNextPageToken)
			if err != nil {
				return fmt.Errorf("ke.v1alpha: failed to ListCluster: %w", err)
			}
			message, err := protojson.Marshal(&res)
			if err != nil {
				return err
			}
			cmd.Println(string(message))
			return nil
		},
	}
	return cmd
}
