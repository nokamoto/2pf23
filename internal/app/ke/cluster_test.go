package ke

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/nokamoto/2pf23/internal/app"
	"github.com/nokamoto/2pf23/internal/app/ke/mock"
	"github.com/nokamoto/2pf23/internal/util/helper"
	"github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestCluster_Create(t *testing.T) {
	runtimeError := errors.New("runtime error")

	tests := []struct {
		name string
		req  *kev1alpha.Cluster
		mock func(*mockke.Mockruntime)
		want *kev1alpha.Cluster
		err  error
	}{
		{
			name: "ok",
			req:  &kev1alpha.Cluster{},
			mock: func(rt *mockke.Mockruntime) {
				gomock.InOrder(
					rt.EXPECT().NewID().Return("new-id"),
					rt.EXPECT().Create(gomock.Any(), helper.ProtoEqual(&kev1alpha.Cluster{
						Name:     "projects/unspecified/clusters/new-id",
						NumNodes: 3,
					})).Return(nil),
				)
			},
			want: &kev1alpha.Cluster{
				Name:     "projects/unspecified/clusters/new-id",
				NumNodes: 3,
			},
		},
		{
			name: "invalid num nodes",
			req: &kev1alpha.Cluster{
				NumNodes: 6,
			},
			err: app.ErrInvalidArgument,
		},
		{
			name: "runtime error",
			req:  &kev1alpha.Cluster{},
			mock: func(rt *mockke.Mockruntime) {
				gomock.InOrder(
					rt.EXPECT().NewID().Return("new-id"),
					rt.EXPECT().Create(gomock.Any(), gomock.Any()).Return(runtimeError),
				)
			},
			err: runtimeError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			rt := mockke.NewMockruntime(ctrl)
			sut := &Cluster{
				rt: rt,
			}
			if tt.mock != nil {
				tt.mock(rt)
			}

			got, err := sut.Create(context.TODO(), tt.req)
			if !errors.Is(err, tt.err) {
				t.Errorf("Cluster.Create() error = %v, wantErr %v", err, tt.err)
			}

			if diff := cmp.Diff(got, tt.want, protocmp.Transform()); diff != "" {
				t.Errorf("Cluster.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
