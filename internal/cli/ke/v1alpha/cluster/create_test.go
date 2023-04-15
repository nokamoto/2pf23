package kev1alphacluster

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
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/testing/protocmp"
)

func Test_Create(t *testing.T) {
	testcases := []struct {
		name     string
		args     string
		mock     func(*mockruntime.MockRuntime, *mock_kev1alpha.MockKeServiceClient)
		expected *kev1alpha.Cluster
		err      error
	}{
		{
			name: "created",
			args: "",
			mock: func(rt *mockruntime.MockRuntime, c *mock_kev1alpha.MockKeServiceClient) {
				gomock.InOrder(
					rt.EXPECT().KeV1Alpha(gomock.Any()).Return(c, nil),
					rt.EXPECT().Context(gomock.Any()).Return(context.TODO()),
					c.EXPECT().CreateCluster(gomock.Any(), gomock.Any()).Return(&kev1alpha.Cluster{
						Name:        "foo",
						DisplayName: "bar",
					}, nil),
				)
			},
			expected: &kev1alpha.Cluster{
				Name:        "foo",
				DisplayName: "bar",
			},
		},
		{
			name: "with display-name",
			args: "--display-name bar",
			mock: func(rt *mockruntime.MockRuntime, c *mock_kev1alpha.MockKeServiceClient) {
				gomock.InOrder(
					rt.EXPECT().KeV1Alpha(gomock.Any()).Return(c, nil),
					rt.EXPECT().Context(gomock.Any()).Return(context.TODO()),
					c.EXPECT().CreateCluster(gomock.Any(), &kev1alpha.CreateClusterRequest{
						Cluster: &kev1alpha.Cluster{
							DisplayName: "bar",
						},
					}).Return(&kev1alpha.Cluster{
						Name:        "foo",
						DisplayName: "bar",
					}, nil),
				)
			},
			expected: &kev1alpha.Cluster{
				Name:        "foo",
				DisplayName: "bar",
			},
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

			cmd := newCreate(rt)
			cmd.SetArgs(strings.Split(tc.args, " "))
			var stdout bytes.Buffer
			cmd.SetOut(&stdout)

			err := cmd.Execute()
			if !errors.Is(err, tc.err) {
				t.Errorf("got %v, want %v", err, tc.err)
			}

			var actual kev1alpha.Cluster
			if err = protojson.Unmarshal(stdout.Bytes(), &actual); err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tc.expected, &actual, protocmp.Transform()); diff != "" {
				t.Errorf("diff: (-want +got)\n%s", diff)
			}
		})
	}
}
