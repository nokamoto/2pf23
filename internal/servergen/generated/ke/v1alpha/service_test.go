package v1alpha

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	mockv1alpha "github.com/nokamoto/2pf23/internal/servergen/generated/ke/v1alpha/mock"
	"github.com/nokamoto/2pf23/internal/util/helper"
	kev1alpha "github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"
)

func Test_CreateCluster(t *testing.T) {
	testcases := []struct {
		name     string
		req      *kev1alpha.CreateClusterRequest
		mock     func(*mockv1alpha.Mockruntime)
		expected *kev1alpha.Cluster
		code     codes.Code
	}{
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
			name: "unknown error",
			req:  &kev1alpha.CreateClusterRequest{},
			mock: func(rt *mockv1alpha.Mockruntime) {
				rt.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, errors.New("unknown"))
			},
			code: codes.Unknown,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			rt := mockv1alpha.NewMockruntime(ctrl)
			s := NewService(zaptest.NewLogger(t), rt)

			if tc.mock != nil {
				tc.mock(rt)
			}

			res, err := s.CreateCluster(context.TODO(), tc.req)
			if code := status.Code(err); code != tc.code {
				t.Errorf("expected %v, got %v", tc.code, code)
			}

			if diff := cmp.Diff(res, tc.expected, protocmp.Transform()); diff != "" {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
		})
	}
}
