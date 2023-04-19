package generated

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nokamoto/2pf23/internal/cli/runtime/mock"
	"github.com/nokamoto/2pf23/internal/mock/pkg/api/ke/v1alpha"
)

func Test_newCreate(t *testing.T) {
	testcases := []struct {
		name string
		args string
		mock func(*mockruntime.MockRuntime, *mock_kev1alpha.MockKeServiceClient)
		err  error
	}{
		{
			name: "ok",
			mock: func(rt *mockruntime.MockRuntime, c *mock_kev1alpha.MockKeServiceClient) {
				gomock.InOrder(
					rt.EXPECT().KeV1alpha(gomock.Any()).Return(c, nil),
				)
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
		})
	}
}
