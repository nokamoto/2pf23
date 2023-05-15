package ke

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/go-cmp/cmp"
	"github.com/nokamoto/2pf23/internal/app"
	"github.com/nokamoto/2pf23/internal/app/ke/mock"
	"github.com/nokamoto/2pf23/internal/infra"
	"github.com/nokamoto/2pf23/internal/util/helper"
	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
	"github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
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
			sut := NewCluster(rt)
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
					})).Return(&kev1alpha.Cluster{
						Name:     "projects/unspecified/clusters/new-id",
						NumNodes: 3,
					}, nil),
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
					rt.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, runtimeError),
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

func TestCluster_Delete(t *testing.T) {
	runtimeError := errors.New("runtime error")
	name := "projects/unspecified/clusters/cluster-id"
	tests := []testcase[string, *empty.Empty]{
		{
			name: "ok",
			req:  name,
			mock: func(rt *mockke.Mockruntime) {
				rt.EXPECT().Delete(gomock.Any(), name).Return(nil)
			},
			want: &empty.Empty{},
		},
		{
			name: "not found",
			req:  name,
			mock: func(rt *mockke.Mockruntime) {
				rt.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(infra.ErrNotFound)
			},
			err: app.ErrNotFound,
		},
		{
			name: "runtime error",
			req:  name,
			mock: func(rt *mockke.Mockruntime) {
				rt.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(runtimeError)
			},
			err: runtimeError,
		},
	}
	run(t, (*Cluster).Delete, tests)
}

func TestCluster_List(t *testing.T) {
	type testcase struct {
		name     string
		pageSize int32
		mock     func(*mockke.Mockruntime)
		want     []*kev1alpha.Cluster
		wantPage *v1.Pagination
		err      error
	}

	testPageSize := func(pageSize int32, want int32) testcase {
		any1, any2 := []*kev1alpha.Cluster{}, &v1.Pagination{}
		return testcase{
			name:     fmt.Sprintf("page size %d", pageSize),
			pageSize: pageSize,
			mock: func(rt *mockke.Mockruntime) {
				rt.EXPECT().List(gomock.Any(), want, gomock.Any()).Return(any1, any2, nil)
			},
			want:     any1,
			wantPage: any2,
		}
	}

	runtimeError := errors.New("runtime error")
	page := &v1.Pagination{}
	tests := []testcase{
		{
			name:     "ok",
			pageSize: 10,
			mock: func(rt *mockke.Mockruntime) {
				rt.EXPECT().List(gomock.Any(), int32(10), page).Return(
					[]*kev1alpha.Cluster{
						{
							Name: "projects/unspecified/clusters/cluster-id",
						},
					},
					&v1.Pagination{
						Cursor: 11,
					},
					nil,
				)
			},
			want: []*kev1alpha.Cluster{
				{
					Name: "projects/unspecified/clusters/cluster-id",
				},
			},
			wantPage: &v1.Pagination{
				Cursor: 11,
			},
		},
		testPageSize(0, 30),
		testPageSize(29, 29),
		testPageSize(30, 30),
		testPageSize(31, 30),
		{
			name:     "page size is negative",
			pageSize: -1,
			err:      app.ErrInvalidArgument,
		},
		{
			name: "runtime error",
			mock: func(rt *mockke.Mockruntime) {
				rt.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil, runtimeError)
			},
			err: runtimeError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			rt := mockke.NewMockruntime(ctrl)
			if tt.mock != nil {
				tt.mock(rt)
			}
			sut := NewCluster(rt)

			got, gotPage, err := sut.List(context.TODO(), tt.pageSize, page)
			if !errors.Is(err, tt.err) {
				t.Errorf("error = %v, wantErr %v", err, tt.err)
			}
			if diff := cmp.Diff(got, tt.want, protocmp.Transform()); diff != "" {
				t.Error(diff)
			}
			if diff := cmp.Diff(gotPage, tt.wantPage, protocmp.Transform()); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestCluster_Update(t *testing.T) {
	type testcase struct {
		name string
		req  *kev1alpha.Cluster
		mask *fieldmaskpb.FieldMask
		mock func(*mockke.Mockruntime)
		want *kev1alpha.Cluster
		err  error
	}

	runtimeError := errors.New("runtime error")

	tests := []testcase{
		{
			name: "ok",
			req: &kev1alpha.Cluster{
				Name:        "projects/unspecified/clusters/cluster-id",
				DisplayName: "cluster-name",
			},
			mask: &fieldmaskpb.FieldMask{
				Paths: []string{"display_name"},
			},
			mock: func(rt *mockke.Mockruntime) {
				rt.EXPECT().Update(gomock.Any(), helper.ProtoEqual(
					&kev1alpha.Cluster{
						Name:        "projects/unspecified/clusters/cluster-id",
						DisplayName: "cluster-name",
					},
				), helper.ProtoEqual(&fieldmaskpb.FieldMask{
					Paths: []string{"display_name"},
				})).Return(&kev1alpha.Cluster{
					Name:        "projects/unspecified/clusters/cluster-id",
					DisplayName: "cluster-name",
					NumNodes:    3,
				}, nil)
			},
			want: &kev1alpha.Cluster{
				Name:        "projects/unspecified/clusters/cluster-id",
				DisplayName: "cluster-name",
				NumNodes:    3,
			},
		},
		{
			name: "invalid mask",
			mask: &fieldmaskpb.FieldMask{
				Paths: []string{"invalid"},
			},
			err: app.ErrInvalidArgument,
		},
		{
			name: "invalid argument if num nodes is zero",
			req:  &kev1alpha.Cluster{},
			mask: &fieldmaskpb.FieldMask{
				Paths: []string{"num_nodes"},
			},
			err: app.ErrInvalidArgument,
		},
		{
			name: "invalid argument if num nodes is greater than 5",
			req: &kev1alpha.Cluster{
				NumNodes: 6,
			},
			mask: &fieldmaskpb.FieldMask{
				Paths: []string{"num_nodes"},
			},
			err: app.ErrInvalidArgument,
		},
		{
			name: "empty field mask",
			req: &kev1alpha.Cluster{
				Name:        "projects/unspecified/clusters/cluster-id",
				DisplayName: "cluster-name",
			},
			mask: &fieldmaskpb.FieldMask{},
			err:  app.ErrInvalidArgument,
		},
		{
			name: "not found",
			req: &kev1alpha.Cluster{
				Name:        "projects/unspecified/clusters/cluster-id",
				DisplayName: "cluster-name",
			},
			mask: &fieldmaskpb.FieldMask{
				Paths: []string{"display_name"},
			},
			mock: func(rt *mockke.Mockruntime) {
				rt.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, infra.ErrNotFound)
			},
			err: app.ErrNotFound,
		},
		{
			name: "runtime error",
			req: &kev1alpha.Cluster{
				Name:        "projects/unspecified/clusters/cluster-id",
				DisplayName: "cluster-name",
			},
			mask: &fieldmaskpb.FieldMask{
				Paths: []string{"display_name"},
			},
			mock: func(rt *mockke.Mockruntime) {
				rt.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, runtimeError)
			},
			err: runtimeError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			rt := mockke.NewMockruntime(ctrl)
			if tt.mock != nil {
				tt.mock(rt)
			}
			sut := NewCluster(rt)

			got, err := sut.Update(context.TODO(), tt.req, tt.mask)
			if !errors.Is(err, tt.err) {
				t.Errorf("error = %v, wantErr %v", err, tt.err)
			}
			if diff := cmp.Diff(got, tt.want, protocmp.Transform()); diff != "" {
				t.Error(diff)
			}
		})
	}
}
