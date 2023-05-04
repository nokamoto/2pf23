package ke

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/nokamoto/2pf23/internal/app"
	"github.com/nokamoto/2pf23/internal/app/ke/mock"
	"github.com/nokamoto/2pf23/internal/infra"
	"github.com/nokamoto/2pf23/internal/util/helper"
	"github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
	"google.golang.org/protobuf/testing/protocmp"
)

type testcase[T1 any, T2 any] struct {
	name string
	req  T1
	mock func(*mockke.Mockruntime)
	want T2
	err  error
}

func run[T1 any, T2 any](t *testing.T, f func(*Cluster, context.Context, T1) (T2, error), tests []testcase[T1, T2]) {
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

			got, err := f(sut, context.Background(), tt.req)
			if !errors.Is(err, tt.err) {
				t.Errorf("error = %v, wantErr %v", err, tt.err)
			}

			if diff := cmp.Diff(got, tt.want, protocmp.Transform()); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestCluster_Create(t *testing.T) {
	runtimeError := errors.New("runtime error")

	tests := []testcase[*kev1alpha.Cluster, *kev1alpha.Cluster]{
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

	run(t, (*Cluster).Create, tests)
}

func TestCluster_Get(t *testing.T) {
	runtimeError := errors.New("runtime error")
	name := "projects/unspecified/clusters/cluster-id"

	tests := []testcase[string, *kev1alpha.Cluster]{
		{
			name: "ok",
			req:  name,
			mock: func(rt *mockke.Mockruntime) {
				rt.EXPECT().Get(gomock.Any(), name).Return(&kev1alpha.Cluster{
					Name:     name,
					NumNodes: 3,
				}, nil)
			},
			want: &kev1alpha.Cluster{
				Name:     name,
				NumNodes: 3,
			},
		},
		{
			name: "not found",
			req:  name,
			mock: func(rt *mockke.Mockruntime) {
				rt.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, infra.ErrNotFound)
			},
			err: app.ErrNotFound,
		},
		{
			name: "runtime error",
			req:  name,
			mock: func(rt *mockke.Mockruntime) {
				rt.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, runtimeError)
			},
			err: runtimeError,
		},
	}

	run(t, (*Cluster).Get, tests)
}
