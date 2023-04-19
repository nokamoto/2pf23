package kev1alphacluster

import (
	"fmt"

	"github.com/nokamoto/2pf23/internal/cli/runtime"
	"github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

func newCreate(rt runtime.Runtime) *cobra.Command {
	var dislapyName string
	cmd := &cobra.Command{
		Use:          "create",
		Short:        "create is a CLI for creating the kubernetes cluster.",
		Long:         `create is a CLI for creating the kubernetes cluster.`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := rt.KeV1alpha(cmd)
			if err != nil {
				return fmt.Errorf("failed to create a client: %w", err)
			}
			ctx := rt.Context(cmd)
			res, err := c.CreateCluster(ctx, &kev1alpha.CreateClusterRequest{
				Cluster: &kev1alpha.Cluster{
					DisplayName: dislapyName,
				},
			})
			if err != nil {
				return fmt.Errorf("failed to create a cluster: %w", err)
			}
			message, err := protojson.Marshal(res)
			if err != nil {
				return err
			}
			cmd.Println(string(message))
			return nil
		},
	}
	cmd.PersistentFlags().StringVar(&dislapyName, "display-name", "", "Display name of the cluster.")
	return cmd
}
