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
	"github.com/nokamoto/2pf23/internal/util/helper/mock"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/testing/protocmp"
)

type testcase[T any] struct {
	name     string
	args     string
	mock     func(*mockruntime.MockRuntime, *mockhelper.MockKeServiceClient)
	expected *T
	err      error
}

func run[T any](t *testing.T, tt []testcase[T], f func(runtime.Runtime) *cobra.Command, unmarshal func([]byte) (*T, error)) {
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			rt := mockruntime.NewMockRuntime(ctrl)
			client := mockhelper.NewMockKeServiceClient(ctrl)
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

			actual, err := unmarshal(stdout.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tc.expected, actual, protocmp.Transform()); diff != "" {
				t.Errorf("diff: %s", diff)
			}
		})
	}
}
