// Code generated by cli-gen. DO NOT EDIT.
package cligen

import (
	"github.com/spf13/cobra"
)

func newCreate() *cobra.Command {
	var stringFlag string

	cmd := &cobra.Command{
		Use:          "use",
		Short:        "short",
		Long:         `long`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
            return nil
		},
	}
	cmd.Flags().StringVar(&stringFlag, "string-flag", "value", "usage")

	return cmd
}