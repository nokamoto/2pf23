package generated

import (
	"bytes"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	mockruntime "github.com/nokamoto/2pf23/internal/cli/runtime/mock"
)

func TestNewRoot(t *testing.T) {
	testcases := []struct {
		name string
		args string
	}{
		{
			name: "empty",
			args: "",
		},
		{
			name: "create",
			args: "create --help",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			rt := mockruntime.NewMockRuntime(ctrl)
			cmd := NewRoot(rt)
			cmd.SetArgs(strings.Split(tc.args, " "))
			var stdout bytes.Buffer
			cmd.SetOut(&stdout)

			err := cmd.Execute()
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
