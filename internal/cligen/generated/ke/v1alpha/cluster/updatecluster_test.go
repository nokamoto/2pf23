package cluster

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	mockruntime "github.com/nokamoto/2pf23/internal/cli/runtime/mock"
	mock_kev1alpha "github.com/nokamoto/2pf23/internal/mock/pkg/api/ke/v1alpha"
	"github.com/nokamoto/2pf23/internal/util/helper"
	kev1alpha "github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func Test_newUpdateCluster(t *testing.T) {
	mask := func(paths ...string) *fieldmaskpb.FieldMask {
		m, err := fieldmaskpb.New(&kev1alpha.Cluster{}, paths...)
		if err != nil {
			t.Fatalf("failed to create a field mask: %v", err)
		}
		return m
	}

	clientErr := errors.New("client error")
	rpcErr := errors.New("rpc error")

	tt := []testcase[kev1alpha.Cluster]{
		{
			name: "unchanged",
			args: "foo",
			mock: func(rt *mockruntime.MockRuntime, c *mock_kev1alpha.MockKeServiceClient) {
				gomock.InOrder(
					rt.EXPECT().Context(gomock.Any()).Return(context.TODO()),
					rt.EXPECT().KeV1alpha(gomock.Any()).Return(c, nil),
					c.EXPECT().UpdateCluster(context.TODO(), helper.ProtoEqual(&kev1alpha.UpdateClusterRequest{
						Cluster: &kev1alpha.Cluster{
							Name: "foo",
						},
						UpdateMask: mask(),
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
			name: "update fields",
			args: "foo --display-name bar --num-nodes 3",
			mock: func(rt *mockruntime.MockRuntime, c *mock_kev1alpha.MockKeServiceClient) {
				gomock.InOrder(
					rt.EXPECT().Context(gomock.Any()).Return(context.TODO()),
					rt.EXPECT().KeV1alpha(gomock.Any()).Return(c, nil),
					c.EXPECT().UpdateCluster(context.TODO(), helper.ProtoEqual(&kev1alpha.UpdateClusterRequest{
						Cluster: &kev1alpha.Cluster{
							Name:        "foo",
							DisplayName: "bar",
							NumNodes:    3,
						},
						UpdateMask: mask("display_name", "num_nodes"),
					})).Return(&kev1alpha.Cluster{
						Name:        "foo",
						DisplayName: "bar",
						NumNodes:    3,
					}, nil),
				)
			},
			expected: &kev1alpha.Cluster{
				Name:        "foo",
				DisplayName: "bar",
				NumNodes:    3,
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
			name: "failed to update a cluster",
			args: "foo",
			mock: func(rt *mockruntime.MockRuntime, c *mock_kev1alpha.MockKeServiceClient) {
				gomock.InOrder(
					rt.EXPECT().Context(gomock.Any()).Return(context.TODO()),
					rt.EXPECT().KeV1alpha(gomock.Any()).Return(c, nil),
					c.EXPECT().UpdateCluster(gomock.Any(), gomock.Any()).Return(nil, rpcErr),
				)
			},
			err: rpcErr,
		},
	}
	run(t, tt, newUpdateCluster, func(b []byte) (*kev1alpha.Cluster, error) {
		var v kev1alpha.Cluster
		err := protojson.Unmarshal(b, &v)
		return &v, err
	})
}
