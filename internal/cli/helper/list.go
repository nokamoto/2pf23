package helper

import (
	"context"

	"google.golang.org/grpc"
)

type response interface {
	GetNextPageToken() string
}

func ListAll[T1 any, T2 response](ctx context.Context, list func(context.Context, T1, ...grpc.CallOption) (T2, error), f1 func(string) T1, f2 func(T2)) error {
	token := ""
	for {
		v, err := list(ctx, f1(token))
		if err != nil {
			return err
		}
		f2(v)
		if token = v.GetNextPageToken(); token == "" {
			break
		}
	}
	return nil
}
