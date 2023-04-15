//go:generate mockgen -source=$GOFILE -package=mock$GOPACKAGE -destination=mock/$GOFILE
package runtime

import (
	"context"

	kev1alpha "github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Runtime interface {
	Context(*cobra.Command) context.Context
	KeV1Alpha(*cobra.Command) (kev1alpha.KeServiceClient, error)
}

type runtime struct{}

func NewRuntime() Runtime {
	return runtime{}
}

func run[T any](cobra *cobra.Command, f func(grpc.ClientConnInterface) T, empty T) (T, error) {
	var target string
	var opts []grpc.DialOption
	target = "localhost:9000"
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(target, opts...)
	if err != nil {
		return empty, err
	}
	return f(conn), nil
}

func (runtime) Context(cobra *cobra.Command) context.Context {
	return context.Background()
}

func (runtime) KeV1Alpha(cobra *cobra.Command) (kev1alpha.KeServiceClient, error) {
	return run(cobra, kev1alpha.NewKeServiceClient, nil)
}
