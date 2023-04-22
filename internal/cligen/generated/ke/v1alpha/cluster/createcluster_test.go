package cluster

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/nokamoto/2pf23/internal/cli/runtime/mock"
	"github.com/nokamoto/2pf23/internal/mock/pkg/api/ke/v1alpha"
	kev1alpha "github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/testing/protocmp"
)

func Test_newCreate(t *testing.T) {
	clientErr := errors.New("client error")
	rpcErr := status.Errorf(codes.Unavailable, "rpc error")
	testcases := []struct {
		name     string
		args     string
		mock     func(*mockruntime.MockRuntime, *mock_kev1alpha.MockKeServiceClient)
		expected *kev1alpha.Cluster
		err      error
	}{
		{
			name: "ok",
			mock: func(rt *mockruntime.MockRuntime, c *mock_kev1alpha.MockKeServiceClient) {
				gomock.InOrder(
					rt.EXPECT().Context(gomock.Any()).Return(context.TODO()),
					rt.EXPECT().KeV1alpha(gomock.Any()).Return(c, nil),
					c.EXPECT().CreateCluster(context.TODO(), gomock.Any()).Return(&kev1alpha.Cluster{
						Name: "foo",
					}, nil),
				)
			},
			expected: &kev1alpha.Cluster{
				Name: "foo",
			},
		},
		{
			name: "set --display-name to display_name",
			args: "--display-name bar",
			mock: func(rt *mockruntime.MockRuntime, c *mock_kev1alpha.MockKeServiceClient) {
				gomock.InOrder(
					rt.EXPECT().Context(gomock.Any()).Return(context.TODO()),
					rt.EXPECT().KeV1alpha(gomock.Any()).Return(c, nil),
					c.EXPECT().CreateCluster(context.TODO(), &kev1alpha.CreateClusterRequest{
						Cluster: &kev1alpha.Cluster{
							DisplayName: "bar",
						},
					}).Return(&kev1alpha.Cluster{}, nil),
				)
			},
			expected: &kev1alpha.Cluster{},
		},
		{
			name: "failed to create a client for ke.v1alpha",
			mock: func(rt *mockruntime.MockRuntime, c *mock_kev1alpha.MockKeServiceClient) {
				gomock.InOrder(
					rt.EXPECT().Context(gomock.Any()).Return(context.TODO()),
					rt.EXPECT().KeV1alpha(gomock.Any()).Return(nil, clientErr),
				)
			},
			err: clientErr,
		},
		{
			name: "failed to CreateCluster",
			mock: func(rt *mockruntime.MockRuntime, c *mock_kev1alpha.MockKeServiceClient) {
				gomock.InOrder(
					rt.EXPECT().Context(gomock.Any()).Return(context.TODO()),
					rt.EXPECT().KeV1alpha(gomock.Any()).Return(c, nil),
					c.EXPECT().CreateCluster(context.TODO(), gomock.Any()).Return(nil, rpcErr),
				)
			},
			err: rpcErr,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			rt := mockruntime.NewMockRuntime(ctrl)
			client := mock_kev1alpha.NewMockKeServiceClient(ctrl)
			if tc.mock != nil {
				tc.mock(rt, client)
			}

			cmd := newCreateCluster(rt)
			cmd.SetArgs(strings.Split(tc.args, " "))
			var stdout bytes.Buffer
			cmd.SetOut(&stdout)

			err := cmd.Execute()
			if !errors.Is(err, tc.err) {
				t.Fatalf("got %v, want %v", err, tc.err)
			}

			if tc.expected == nil && stdout.Len() == 0 {
				return
			}

			var actual kev1alpha.Cluster
			if err := protojson.Unmarshal(stdout.Bytes(), &actual); err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tc.expected, &actual, protocmp.Transform()); diff != "" {
				t.Errorf("diff: %s", diff)
			}
		})
	}
}