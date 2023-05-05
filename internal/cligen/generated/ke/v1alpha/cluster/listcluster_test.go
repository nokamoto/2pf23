package cluster

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	mockruntime "github.com/nokamoto/2pf23/internal/cli/runtime/mock"
	mock_kev1alpha "github.com/nokamoto/2pf23/internal/mock/pkg/api/ke/v1alpha"
	"github.com/nokamoto/2pf23/internal/util/helper"
	kev1alpha "github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
	"google.golang.org/protobuf/encoding/protojson"
)

func Test_newListCluster(t *testing.T) {
	tt := []testcase[kev1alpha.ListClusterResponse]{
		{
			name: "ok",
			args: "",
			mock: func(rt *mockruntime.MockRuntime, c *mock_kev1alpha.MockKeServiceClient) {
				gomock.InOrder(
					rt.EXPECT().Context(gomock.Any()).Return(context.TODO()),
					rt.EXPECT().KeV1alpha(gomock.Any()).Return(c, nil),
					c.EXPECT().ListCluster(context.TODO(), helper.ProtoEqual(
						&kev1alpha.ListClusterRequest{},
					)).Return(&kev1alpha.ListClusterResponse{
						Clusters: []*kev1alpha.Cluster{
							{
								Name: "foo",
							},
						},
					}, nil),
				)
			},
			expected: &kev1alpha.ListClusterResponse{
				Clusters: []*kev1alpha.Cluster{
					{
						Name: "foo",
					},
				},
			},
		},
		{
			name: "call twice",
			args: "",
			mock: func(rt *mockruntime.MockRuntime, c *mock_kev1alpha.MockKeServiceClient) {
				gomock.InOrder(
					rt.EXPECT().Context(gomock.Any()).Return(context.TODO()),
					rt.EXPECT().KeV1alpha(gomock.Any()).Return(c, nil),
					c.EXPECT().ListCluster(context.TODO(), helper.ProtoEqual(
						&kev1alpha.ListClusterRequest{},
					)).Return(&kev1alpha.ListClusterResponse{
						Clusters: []*kev1alpha.Cluster{
							{
								Name: "foo",
							},
						},
						NextPageToken: "bar",
					}, nil),
					c.EXPECT().ListCluster(context.TODO(), helper.ProtoEqual(
						&kev1alpha.ListClusterRequest{
							PageToken: "bar",
						},
					)).Return(&kev1alpha.ListClusterResponse{
						Clusters: []*kev1alpha.Cluster{
							{
								Name: "baz",
							},
						},
					}, nil),
				)
			},
			expected: &kev1alpha.ListClusterResponse{
				Clusters: []*kev1alpha.Cluster{
					{
						Name: "foo",
					},
					{
						Name: "baz",
					},
				},
			},
		},
	}

	run(t, tt, newListCluster, func(b []byte) (*kev1alpha.ListClusterResponse, error) {
		var res kev1alpha.ListClusterResponse
		err := protojson.Unmarshal(b, &res)
		if err != nil {
			return nil, err
		}
		return &res, nil
	})
}
