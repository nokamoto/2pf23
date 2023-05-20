package cluster

import (
	"context"
	"errors"
	"testing"

	"github.com/bufbuild/connect-go"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/nokamoto/2pf23/internal/cli/runtime/mock"
	"github.com/nokamoto/2pf23/internal/util/helper"
	"github.com/nokamoto/2pf23/internal/util/helper/mock"
	kev1alpha "github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
	"google.golang.org/protobuf/encoding/protojson"
)

func Test_newDelete(t *testing.T) {
	clientErr := errors.New("client error")
	rpcErr := connect.NewError(connect.CodeInternal, errors.New("rpc error"))

	testcases := []testcase[empty.Empty]{
		{
			name: "ok",
			args: "foo",
			mock: func(rt *mockruntime.MockRuntime, c *mockhelper.MockKeServiceClient) {
				gomock.InOrder(
					rt.EXPECT().Context(gomock.Any()).Return(context.TODO()),
					rt.EXPECT().KeV1alpha(gomock.Any()).Return(c, nil),
					c.EXPECT().DeleteCluster(context.TODO(), helper.ConnectEqual(&kev1alpha.DeleteClusterRequest{
						Name: "foo",
					})).Return(connect.NewResponse(&empty.Empty{}), nil),
				)
			},
			expected: &empty.Empty{},
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
			name: "failed to DeleteCluster",
			args: "foo",
			mock: func(rt *mockruntime.MockRuntime, c *mockhelper.MockKeServiceClient) {
				gomock.InOrder(
					rt.EXPECT().Context(gomock.Any()).Return(context.TODO()),
					rt.EXPECT().KeV1alpha(gomock.Any()).Return(c, nil),
					c.EXPECT().DeleteCluster(context.TODO(), gomock.Any()).Return(nil, rpcErr),
				)
			},
			err: rpcErr,
		},
	}

	run(t, testcases, newDeleteCluster, func(b []byte) (*empty.Empty, error) {
		var v empty.Empty
		if err := protojson.Unmarshal(b, &v); err != nil {
			return nil, err
		}
		return &v, nil
	})
}
