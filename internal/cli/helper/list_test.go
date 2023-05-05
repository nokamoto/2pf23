package helper

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	mock_kev1alpha "github.com/nokamoto/2pf23/internal/mock/pkg/api/ke/v1alpha"
	"github.com/nokamoto/2pf23/internal/util/helper"
	kev1alpha "github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestListAll(t *testing.T) {
	size := int32(30)
	rpcErr := errors.New("rpc error")

	tests := []struct {
		name string
		mock func(*mock_kev1alpha.MockKeServiceClient)
		want []*kev1alpha.Cluster
		err  error
	}{
		{
			name: "empty",
			mock: func(m *mock_kev1alpha.MockKeServiceClient) {
				gomock.InOrder(
					m.EXPECT().ListCluster(gomock.Any(), helper.ProtoEqual(&kev1alpha.ListClusterRequest{
						PageSize: size,
					})).Return(&kev1alpha.ListClusterResponse{}, nil),
				)
			},
		},
		{
			name: "call twice",
			mock: func(m *mock_kev1alpha.MockKeServiceClient) {
				gomock.InOrder(
					m.EXPECT().ListCluster(gomock.Any(), helper.ProtoEqual(&kev1alpha.ListClusterRequest{
						PageSize:  size,
						PageToken: "",
					})).Return(&kev1alpha.ListClusterResponse{
						Clusters: []*kev1alpha.Cluster{
							{
								Name: "1",
							},
						},
						NextPageToken: "1",
					}, nil),
					m.EXPECT().ListCluster(gomock.Any(), helper.ProtoEqual(&kev1alpha.ListClusterRequest{
						PageSize:  size,
						PageToken: "1",
					})).Return(&kev1alpha.ListClusterResponse{
						Clusters: []*kev1alpha.Cluster{
							{
								Name: "2",
							},
						},
					}, nil),
				)
			},
			want: []*kev1alpha.Cluster{
				{
					Name: "1",
				},
				{
					Name: "2",
				},
			},
		},
		{
			name: "error",
			mock: func(m *mock_kev1alpha.MockKeServiceClient) {
				gomock.InOrder(
					m.EXPECT().ListCluster(gomock.Any(), gomock.Any()).Return(&kev1alpha.ListClusterResponse{
						Clusters: []*kev1alpha.Cluster{
							{
								Name: "1",
							},
						},
						NextPageToken: "1",
					}, nil),
					m.EXPECT().ListCluster(gomock.Any(), gomock.Any()).Return(nil, rpcErr),
				)
			},
			want: []*kev1alpha.Cluster{
				{
					Name: "1",
				},
			},
			err: rpcErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			m := mock_kev1alpha.NewMockKeServiceClient(ctrl)

			if tt.mock != nil {
				tt.mock(m)
			}

			var got kev1alpha.ListClusterResponse
			err := ListAll(
				context.TODO(),
				m.ListCluster,
				func(s string) *kev1alpha.ListClusterRequest {
					return &kev1alpha.ListClusterRequest{
						PageSize:  size,
						PageToken: s,
					}
				},
				func(v *kev1alpha.ListClusterResponse) {
					got.Clusters = append(got.Clusters, v.GetClusters()...)
				},
			)

			if !errors.Is(err, tt.err) {
				t.Errorf("ListAll() error = %v, wantErr %v", err, tt.err)
			}

			want := &kev1alpha.ListClusterResponse{
				Clusters: tt.want,
			}
			if diff := cmp.Diff(&got, want, protocmp.Transform()); diff != "" {
				t.Errorf("ListAll() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
