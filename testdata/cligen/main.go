// Code generated by cli-gen. DO NOT EDIT.
package generated

import (
	"fmt"
)

import (
	"github.com/nokamoto/2pf23/internal/cli/runtime"
	v1alpha "github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

func newCreateCluster(rt runtime.Runtime) *cobra.Command {
	var stringFlag string
	cmd := &cobra.Command{
		Use:          "use",
		Short:        "short",
		Long:         `long`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// standard create method
			ctx := rt.Context(cmd)
			c, err := rt.KeV1alpha(cmd)
			if err != nil {
				return fmt.Errorf("failed to create a client for ke.v1alpha: %w", err)
			}
			res, err := c.CreateCluster(ctx, &v1alpha.CreateClusterRequest{
				Cluster: &v1alpha.Cluster{
					DisplayName: stringFlag,
				},
			})
			if err != nil {
				return fmt.Errorf("ke.v1alpha: failed to CreateCluster: %w", err)
			}
			message, err := protojson.Marshal(res)
			if err != nil {
				return err
			}
			cmd.Println(string(message))
			return nil
		},
	}
	cmd.Flags().StringVar(&stringFlag, "string-flag", "value", "usage")
	return cmd
}
