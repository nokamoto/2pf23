package helper

import (
	"bytes"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nokamoto/2pf23/internal/cli/runtime"
	mockruntime "github.com/nokamoto/2pf23/internal/cli/runtime/mock"
	"github.com/spf13/cobra"
)

type RootTest struct {
	Name string
	Args string
}

type RootTests []RootTest

func (ts RootTests) Run(t *testing.T, f func(rt runtime.Runtime) *cobra.Command) {
	for _, tc := range ts {
		t.Run(tc.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			rt := mockruntime.NewMockRuntime(ctrl)
			cmd := f(rt)
			cmd.SetArgs(strings.Split(tc.Args, " "))
			var stdout bytes.Buffer
			cmd.SetOut(&stdout)

			err := cmd.Execute()
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
