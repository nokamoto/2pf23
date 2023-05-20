package v1alpha

import (
	"context"
	"errors"
	"testing"

	"github.com/bufbuild/connect-go"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/go-cmp/cmp"
	"github.com/nokamoto/2pf23/internal/app"
	helperapi "github.com/nokamoto/2pf23/internal/server/helper"
	mockv1alpha "github.com/nokamoto/2pf23/internal/servergen/generated/ke/v1alpha/mock"
	"github.com/nokamoto/2pf23/internal/util/helper"
	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
	kev1alpha "github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
	"go.uber.org/zap/zaptest"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type testcase[T1 any, T2 any] struct {
	name     string
	req      *T1
	mock     func(*mockv1alpha.Mockruntime)
	expected *T2
	code     connect.Code
}

func run[T1 any, T2 any](t *testing.T, f func(*service, context.Context, *connect.Request[T1]) (*connect.Response[T2], error), tt []testcase[T1, T2]) {
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			rt := mockv1alpha.NewMockruntime(ctrl)
			s := NewService(zaptest.NewLogger(t), rt)

			if tc.mock != nil {
				tc.mock(rt)
			}

			res, err := f(s, context.TODO(), connect.NewRequest(tc.req))

			if code := helperapi.CodeOf(err); code != tc.code {
				t.Errorf("expected %v, got %v", tc.code, code)
			}

			if diff := cmp.Diff(helperapi.GetResponseMsg(res), tc.expected, protocmp.Transform()); diff != "" {
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
			code: connect.CodeInvalidArgument,
		},
		{
			name: "not found",
			mock: func(rt *mockv1alpha.Mockruntime) {
				rt.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, app.ErrNotFound)
			},
			code: connect.CodeNotFound,
		},
		{
			name: "unknown error",
			mock: func(rt *mockv1alpha.Mockruntime) {
				rt.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, errors.New("unknown"))
			},
			code: connect.CodeUnknown,
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
			code: connect.CodeInvalidArgument,
		},
		{
			name: "not found",
			mock: func(rt *mockv1alpha.Mockruntime) {
				rt.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, app.ErrNotFound)
			},
			code: connect.CodeNotFound,
		},
		{
			name: "unknown error",
			mock: func(rt *mockv1alpha.Mockruntime) {
				rt.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, errors.New("unknown"))
			},
			code: connect.CodeUnknown,
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
			code: connect.CodeInvalidArgument,
		},
		{
			name: "not found",
			mock: func(rt *mockv1alpha.Mockruntime) {
				rt.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil, app.ErrNotFound)
			},
			code: connect.CodeNotFound,
		},
		{
			name: "unknown error",
			mock: func(rt *mockv1alpha.Mockruntime) {
				rt.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil, errors.New("unknown"))
			},
			code: connect.CodeUnknown,
		},
	}
	run(t, (*service).DeleteCluster, testcases)
}

func Test_ListCluster(t *testing.T) {
	pageSize := int32(30)
	page := &v1.Pagination{
		Cursor: 100,
	}
	token, err := helperapi.PageToken(page)
	if err != nil {
		t.Fatal(err)
	}
	testcases := []testcase[kev1alpha.ListClusterRequest, kev1alpha.ListClusterResponse]{
		{
			name: "ok",
			req: &kev1alpha.ListClusterRequest{
				PageSize: pageSize,
			},
			mock: func(rt *mockv1alpha.Mockruntime) {
				rt.EXPECT().List(gomock.Any(), pageSize, nil).Return([]*kev1alpha.Cluster{
					{
						Name:        "foo",
						DisplayName: "test",
					},
				}, nil, nil)
			},
			expected: &kev1alpha.ListClusterResponse{
				Clusters: []*kev1alpha.Cluster{
					{
						Name:        "foo",
						DisplayName: "test",
					},
				},
			},
		},
		{
			name: "ok with page token",
			req: &kev1alpha.ListClusterRequest{
				PageSize:  pageSize,
				PageToken: token,
			},
			mock: func(rt *mockv1alpha.Mockruntime) {
				rt.EXPECT().List(gomock.Any(), pageSize, helper.ProtoEqual(page)).Return([]*kev1alpha.Cluster{
					{
						Name:        "foo",
						DisplayName: "test",
					},
				}, nil, nil)
			},
			expected: &kev1alpha.ListClusterResponse{
				Clusters: []*kev1alpha.Cluster{
					{
						Name:        "foo",
						DisplayName: "test",
					},
				},
			},
		},
		{
			name: "invalid page token",
			req: &kev1alpha.ListClusterRequest{
				PageToken: "invalid",
			},
			code: connect.CodeInvalidArgument,
		},
		{
			name: "invalid argument",
			mock: func(rt *mockv1alpha.Mockruntime) {
				rt.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil, app.ErrInvalidArgument)
			},
			code: connect.CodeInvalidArgument,
		},
		{
			name: "not found",
			mock: func(rt *mockv1alpha.Mockruntime) {
				rt.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil, app.ErrNotFound)
			},
			code: connect.CodeNotFound,
		},
		{
			name: "unknown error",
			mock: func(rt *mockv1alpha.Mockruntime) {
				rt.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil, errors.New("unknown"))
			},
			code: connect.CodeUnknown,
		},
	}
	run(t, (*service).ListCluster, testcases)
}

func Test_UpdateCluster(t *testing.T) {
	tests := []testcase[kev1alpha.UpdateClusterRequest, kev1alpha.Cluster]{
		{
			name: "ok",
			req: &kev1alpha.UpdateClusterRequest{
				Cluster: &kev1alpha.Cluster{
					Name:        "foo",
					DisplayName: "test",
				},
				UpdateMask: &fieldmaskpb.FieldMask{
					Paths: []string{"display_name"},
				},
			},
			mock: func(rt *mockv1alpha.Mockruntime) {
				rt.EXPECT().Update(gomock.Any(), helper.ProtoEqual(&kev1alpha.Cluster{
					Name:        "foo",
					DisplayName: "test",
				}), helper.ProtoEqual(&fieldmaskpb.FieldMask{
					Paths: []string{"display_name"},
				})).Return(&kev1alpha.Cluster{
					Name:        "foo",
					DisplayName: "test",
					NumNodes:    3,
				}, nil)
			},
			expected: &kev1alpha.Cluster{
				Name:        "foo",
				DisplayName: "test",
				NumNodes:    3,
			},
		},
		{
			name: "invalid argument",
			mock: func(rt *mockv1alpha.Mockruntime) {
				rt.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, app.ErrInvalidArgument)
			},
			code: connect.CodeInvalidArgument,
		},
		{
			name: "not found",
			mock: func(rt *mockv1alpha.Mockruntime) {
				rt.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, app.ErrNotFound)
			},
			code: connect.CodeNotFound,
		},
	}
	run(t, (*service).UpdateCluster, tests)
}
