package kev1alpha

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
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
		expected *kev1alpha.Cluster
		code     codes.Code
	}{
		{
			name: "unimplemented",
			code: codes.Unimplemented,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewService(zaptest.NewLogger(t))
			actual, err := s.CreateCluster(context.TODO(), tc.req)
			if code := status.Code(err); code != tc.code {
				t.Errorf("status.Code(%v) = %v, want %v", err, code, tc.code)
			}

			if diff := cmp.Diff(tc.expected, actual, protocmp.Transform()); diff != "" {
				t.Errorf("diff: (-want +got)\n%s", diff)
			}
		})
	}
}
