package ke

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
			Name: "v1alpha",
			Args: "v1alpha --help",
		},
	}.Run(t, NewRoot)
}
