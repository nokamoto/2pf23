package cli

import (
	"github.com/nokamoto/2pf23/internal/cli/generated"
	"github.com/nokamoto/2pf23/internal/cli/runtime"
	"github.com/spf13/cobra"
)

func NewCLI() *cobra.Command {
	rt := runtime.NewRuntime()
	return generated.NewRoot(rt)
}
