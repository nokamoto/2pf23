//go:generate go run -mod=mod github.com/golang/mock/mockgen -source=$GOFILE -package=mock$GOPACKAGE -destination=mock/$GOFILE
package runtime

import (
	"context"
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/nokamoto/2pf23/pkg/api/ke/v1alpha/kev1alphaconnect"
	"github.com/spf13/cobra"
)

type Runtime interface {
	Context(*cobra.Command) context.Context
	KeV1alpha(*cobra.Command) (kev1alphaconnect.KeServiceClient, error)
}

type runtime struct{}

func NewRuntime() Runtime {
	return runtime{}
}

func (runtime) Context(cobra *cobra.Command) context.Context {
	return context.Background()
}

func (runtime) KeV1alpha(cobra *cobra.Command) (kev1alphaconnect.KeServiceClient, error) {
	return kev1alphaconnect.NewKeServiceClient(
		http.DefaultClient,
		"http://localhost:9000",
		connect.WithGRPC(),
	), nil
}
