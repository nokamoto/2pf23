package cluster

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nokamoto/2pf23/internal/cli/runtime/mock"
	"github.com/nokamoto/2pf23/internal/mock/pkg/api/ke/v1alpha"
	"github.com/nokamoto/2pf23/internal/util/helper"
	kev1alpha "github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_newGet(t *testing.T) {
	clientErr := errors.New("client error")
	rpcErr := status.Errorf(codes.Unavailable, "rpc error")

	testcases := testcases{
		{
			name: "ok",
			args: "foo",
			mock: func(rt *mockruntime.MockRuntime, c *mock_kev1alpha.MockKeServiceClient) {
				gomock.InOrder(
					rt.EXPECT().Context(gomock.Any()).Return(context.TODO()),
					rt.EXPECT().KeV1alpha(gomock.Any()).Return(c, nil),
					c.EXPECT().GetCluster(context.TODO(), helper.ProtoEqual(&kev1alpha.GetClusterRequest{
						Name: "foo",
					})).Return(&kev1alpha.Cluster{
						Name: "foo",
					}, nil),
				)
			},
			expected: &kev1alpha.Cluster{
				Name: "foo",
			},
		},
		{
			name: "failed to get a client for ke.v1alpha",
			args: "foo",
			mock: func(rt *mockruntime.MockRuntime, c *mock_kev1alpha.MockKeServiceClient) {
				gomock.InOrder(
					rt.EXPECT().Context(gomock.Any()).Return(context.TODO()),
					rt.EXPECT().KeV1alpha(gomock.Any()).Return(nil, clientErr),
				)
			},
			err: clientErr,
		},
		{
			name: "failed to GetCluster",
			args: "foo",
			mock: func(rt *mockruntime.MockRuntime, c *mock_kev1alpha.MockKeServiceClient) {
				gomock.InOrder(
					rt.EXPECT().Context(gomock.Any()).Return(context.TODO()),
					rt.EXPECT().KeV1alpha(gomock.Any()).Return(c, nil),
					c.EXPECT().GetCluster(context.TODO(), gomock.Any()).Return(nil, rpcErr),
				)
			},
			err: rpcErr,
		},
	}

	testcases.run(t, newGetCluster)
}
