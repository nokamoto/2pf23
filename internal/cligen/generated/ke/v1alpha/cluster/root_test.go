package cluster

import (
	"testing"

	"github.com/nokamoto/2pf23/internal/cligen/generated/helper"
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
	}.Run(t, NewRoot)
}
