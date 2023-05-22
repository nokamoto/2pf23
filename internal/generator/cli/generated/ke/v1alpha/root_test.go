package v1alpha

import (
	"testing"

	"github.com/nokamoto/2pf23/internal/generator/cli/generated/helper"
)

func TestNewRoot(t *testing.T) {
	helper.RootTests{
		{
			Name: "empty",
			Args: "",
		},
		{
			Name: "cluster",
			Args: "cluster --help",
		},
	}.Run(t, NewRoot)
}
