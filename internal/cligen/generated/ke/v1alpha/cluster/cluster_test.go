package cluster

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/nokamoto/2pf23/internal/cli/runtime"
	"github.com/nokamoto/2pf23/internal/cli/runtime/mock"
	"github.com/nokamoto/2pf23/internal/mock/pkg/api/ke/v1alpha"
	kev1alpha "github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/testing/protocmp"
)

type testcase struct {
	name     string
	args     string
	mock     func(*mockruntime.MockRuntime, *mock_kev1alpha.MockKeServiceClient)
	expected *kev1alpha.Cluster
	err      error
}

type testcases []testcase

func (tt testcases) run(t *testing.T, f func(runtime.Runtime) *cobra.Command) {
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			rt := mockruntime.NewMockRuntime(ctrl)
			client := mock_kev1alpha.NewMockKeServiceClient(ctrl)
			if tc.mock != nil {
				tc.mock(rt, client)
			}

			var args []string
			for _, arg := range strings.Split(tc.args, " ") {
				if arg != "" {
					args = append(args, arg)
				}
			}

			cmd := f(rt)
			cmd.SetArgs(args)
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
