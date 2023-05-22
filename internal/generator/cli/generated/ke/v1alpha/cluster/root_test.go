package cluster

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
			Name: "create",
			Args: "create --help",
		},
		{
			Name: "get",
			Args: "get --help",
		},
	}.Run(t, NewRoot)
}
