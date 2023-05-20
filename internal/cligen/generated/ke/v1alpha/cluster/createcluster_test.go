package cluster

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/bufbuild/connect-go"
	"github.com/golang/mock/gomock"
	"github.com/nokamoto/2pf23/internal/cli/runtime/mock"
	"github.com/nokamoto/2pf23/internal/util/helper"
	"github.com/nokamoto/2pf23/internal/util/helper/mock"
	kev1alpha "github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
	"google.golang.org/protobuf/encoding/protojson"
)

func Test_newCreate(t *testing.T) {
	clientErr := errors.New("client error")
	rpcErr := connect.NewError(connect.CodeInternal, errors.New("rpc error"))

	set := func(args string, cluster *kev1alpha.Cluster) testcase[kev1alpha.Cluster] {
		return testcase[kev1alpha.Cluster]{
			name: fmt.Sprintf("set %s", args),
			args: args,
			mock: func(rt *mockruntime.MockRuntime, c *mockhelper.MockKeServiceClient) {
				gomock.InOrder(
					rt.EXPECT().Context(gomock.Any()).Return(context.TODO()),
					rt.EXPECT().KeV1alpha(gomock.Any()).Return(c, nil),
					c.EXPECT().CreateCluster(context.TODO(), helper.ConnectEqual(&kev1alpha.CreateClusterRequest{
						Cluster: cluster,
					})).Return(connect.NewResponse(cluster), nil),
				)
			},
			expected: cluster,
		}
	}

	testcases := []testcase[kev1alpha.Cluster]{
		{
			name: "ok",
			mock: func(rt *mockruntime.MockRuntime, c *mockhelper.MockKeServiceClient) {
				gomock.InOrder(
					rt.EXPECT().Context(gomock.Any()).Return(context.TODO()),
					rt.EXPECT().KeV1alpha(gomock.Any()).Return(c, nil),
					c.EXPECT().CreateCluster(context.TODO(), gomock.Any()).Return(connect.NewResponse(&kev1alpha.Cluster{
						Name: "foo",
					}), nil),
				)
			},
			expected: &kev1alpha.Cluster{
				Name: "foo",
			},
		},
		set("--display-name bar", &kev1alpha.Cluster{
			DisplayName: "bar",
		}),
		set("--num-nodes 3", &kev1alpha.Cluster{
			NumNodes: 3,
		}),
		set("--machine-type MACHINE_TYPE_STANDARD", &kev1alpha.Cluster{
			MachineType: kev1alpha.MachineType_MACHINE_TYPE_STANDARD,
		}),
		{
			name: "failed to create a client for ke.v1alpha",
			mock: func(rt *mockruntime.MockRuntime, c *mockhelper.MockKeServiceClient) {
				gomock.InOrder(
					rt.EXPECT().Context(gomock.Any()).Return(context.TODO()),
					rt.EXPECT().KeV1alpha(gomock.Any()).Return(nil, clientErr),
				)
			},
			err: clientErr,
		},
		{
			name: "failed to CreateCluster",
			mock: func(rt *mockruntime.MockRuntime, c *mockhelper.MockKeServiceClient) {
				gomock.InOrder(
					rt.EXPECT().Context(gomock.Any()).Return(context.TODO()),
					rt.EXPECT().KeV1alpha(gomock.Any()).Return(c, nil),
					c.EXPECT().CreateCluster(context.TODO(), gomock.Any()).Return(nil, rpcErr),
				)
			},
			err: rpcErr,
		},
	}

	run(t, testcases, newCreateCluster, func(b []byte) (*kev1alpha.Cluster, error) {
		var cluster kev1alpha.Cluster
		if err := protojson.Unmarshal(b, &cluster); err != nil {
			return nil, err
		}
		return &cluster, nil
	})
}
