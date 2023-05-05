package v1alpha

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/go-cmp/cmp"
	"github.com/nokamoto/2pf23/internal/app"
	mockv1alpha "github.com/nokamoto/2pf23/internal/servergen/generated/ke/v1alpha/mock"
	"github.com/nokamoto/2pf23/internal/util/helper"
	kev1alpha "github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"
)

type testcase[T1 any, T2 any] struct {
	name     string
	req      *T1
	mock     func(*mockv1alpha.Mockruntime)
	expected *T2
	code     codes.Code
}

func run[T1 any, T2 any](t *testing.T, f func(*service, context.Context, *T1) (*T2, error), tt []testcase[T1, T2]) {
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			rt := mockv1alpha.NewMockruntime(ctrl)
			s := NewService(zaptest.NewLogger(t), rt)

			if tc.mock != nil {
				tc.mock(rt)
			}

			res, err := f(s, context.TODO(), tc.req)
			if code := status.Code(err); code != tc.code {
				t.Errorf("expected %v, got %v", tc.code, code)
			}

			if diff := cmp.Diff(res, tc.expected, protocmp.Transform()); diff != "" {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
		})
	}
}

func Test_CreateCluster(t *testing.T) {
	testcases := []testcase[kev1alpha.CreateClusterRequest, kev1alpha.Cluster]{
		{
			name: "ok",
			req: &kev1alpha.CreateClusterRequest{
				Cluster: &kev1alpha.Cluster{
					DisplayName: "test",
				},
			},
			mock: func(rt *mockv1alpha.Mockruntime) {
				rt.EXPECT().Create(gomock.Any(), helper.ProtoEqual(&kev1alpha.Cluster{
					DisplayName: "test",
				})).Return(&kev1alpha.Cluster{
					Name:        "foo",
					DisplayName: "test",
				}, nil)
			},
			expected: &kev1alpha.Cluster{
				Name:        "foo",
				DisplayName: "test",
			},
		},
		{
			name: "invalid argument",
			mock: func(rt *mockv1alpha.Mockruntime) {
				rt.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, app.ErrInvalidArgument)
			},
			code: codes.InvalidArgument,
		},
		{
			name: "not found",
			mock: func(rt *mockv1alpha.Mockruntime) {
				rt.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, app.ErrNotFound)
			},
			code: codes.NotFound,
		},
		{
			name: "unknown error",
			mock: func(rt *mockv1alpha.Mockruntime) {
				rt.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, errors.New("unknown"))
			},
			code: codes.Unknown,
		},
	}

	run(t, (*service).CreateCluster, testcases)
}

func Test_GetCluster(t *testing.T) {
	testcases := []testcase[kev1alpha.GetClusterRequest, kev1alpha.Cluster]{
		{
			name: "ok",
			req: &kev1alpha.GetClusterRequest{
				Name: "foo",
			},
			mock: func(rt *mockv1alpha.Mockruntime) {
				rt.EXPECT().Get(gomock.Any(), "foo").Return(&kev1alpha.Cluster{
					Name:        "foo",
					DisplayName: "test",
				}, nil)
			},
			expected: &kev1alpha.Cluster{
				Name:        "foo",
				DisplayName: "test",
			},
		},
		{
			name: "invalid argument",
			mock: func(rt *mockv1alpha.Mockruntime) {
				rt.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, app.ErrInvalidArgument)
			},
			code: codes.InvalidArgument,
		},
		{
			name: "not found",
			mock: func(rt *mockv1alpha.Mockruntime) {
				rt.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, app.ErrNotFound)
			},
			code: codes.NotFound,
		},
		{
			name: "unknown error",
			mock: func(rt *mockv1alpha.Mockruntime) {
				rt.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, errors.New("unknown"))
			},
			code: codes.Unknown,
		},
	}
	run(t, (*service).GetCluster, testcases)
}

func Test_DeleteCluster(t *testing.T) {
	testcases := []testcase[kev1alpha.DeleteClusterRequest, empty.Empty]{
		{
			name: "ok",
			req: &kev1alpha.DeleteClusterRequest{
				Name: "foo",
			},
			mock: func(rt *mockv1alpha.Mockruntime) {
				rt.EXPECT().Delete(gomock.Any(), "foo").Return(&empty.Empty{}, nil)
			},
			expected: &empty.Empty{},
		},
		{
			name: "invalid argument",
			mock: func(rt *mockv1alpha.Mockruntime) {
				rt.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil, app.ErrInvalidArgument)
			},
			code: codes.InvalidArgument,
		},
		{
			name: "not found",
			mock: func(rt *mockv1alpha.Mockruntime) {
				rt.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil, app.ErrNotFound)
			},
			code: codes.NotFound,
		},
		{
			name: "unknown error",
			mock: func(rt *mockv1alpha.Mockruntime) {
				rt.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil, errors.New("unknown"))
			},
			code: codes.Unknown,
		},
	}
	run(t, (*service).DeleteCluster, testcases)
}
