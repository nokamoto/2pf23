package cluster

import (
	"context"
	"errors"
	"testing"

	"github.com/bufbuild/connect-go"
	"github.com/golang/mock/gomock"
	"github.com/nokamoto/2pf23/internal/cli/runtime/mock"
	"github.com/nokamoto/2pf23/internal/util/helper"
	"github.com/nokamoto/2pf23/internal/util/helper/mock"
	kev1alpha "github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
	"google.golang.org/protobuf/encoding/protojson"
)

func Test_newGet(t *testing.T) {
	clientErr := errors.New("client error")
	rpcErr := connect.NewError(connect.CodeInternal, errors.New("rpc error"))

	testcases := []testcase[kev1alpha.Cluster]{
		{
			name: "ok",
			args: "foo",
			mock: func(rt *mockruntime.MockRuntime, c *mockhelper.MockKeServiceClient) {
				gomock.InOrder(
					rt.EXPECT().Context(gomock.Any()).Return(context.TODO()),
					rt.EXPECT().KeV1alpha(gomock.Any()).Return(c, nil),
					c.EXPECT().GetCluster(context.TODO(), helper.ConnectEqual(&kev1alpha.GetClusterRequest{
						Name: "foo",
					})).Return(connect.NewResponse(&kev1alpha.Cluster{
						Name: "foo",
					}), nil),
				)
			},
			expected: &kev1alpha.Cluster{
				Name: "foo",
			},
		},
		{
			name: "failed to get a client for ke.v1alpha",
			args: "foo",
			mock: func(rt *mockruntime.MockRuntime, c *mockhelper.MockKeServiceClient) {
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
			mock: func(rt *mockruntime.MockRuntime, c *mockhelper.MockKeServiceClient) {
				gomock.InOrder(
					rt.EXPECT().Context(gomock.Any()).Return(context.TODO()),
					rt.EXPECT().KeV1alpha(gomock.Any()).Return(c, nil),
					c.EXPECT().GetCluster(context.TODO(), gomock.Any()).Return(nil, rpcErr),
				)
			},
			err: rpcErr,
		},
	}

	run(t, testcases, newGetCluster, func(b []byte) (*kev1alpha.Cluster, error) {
		var cluster kev1alpha.Cluster
		if err := protojson.Unmarshal(b, &cluster); err != nil {
			return nil, err
		}
		return &cluster, nil
	})
}
