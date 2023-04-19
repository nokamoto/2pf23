package main

import (
	"bytes"
	"os"

	"github.com/nokamoto/2pf23/internal/cligen"
	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
)

// cligen-testdata generates internal/cligen/generated/command.go.
func main() {
	p := cligen.Printer{}

	cmd := &v1.Command{
		Api:        "ke",
		ApiVersion: "v1alpha",
		Package:    "cligen",
		Use:        "use",
		Short:      "short",
		Long:       "long",
		Method:     "Create",
		MethodType: v1.MethodType_METHOD_TYPE_CREATE,
		StringFlags: []*v1.Flag{
			{
				Name:        "stringFlag",
				DisplayName: "string-flag",
				Value:       "value",
				Usage:       "usage",
			},
		},
	}

	var buf bytes.Buffer
	if err := p.PrintCommand(&buf, cmd); err != nil {
		panic(err)
	}
	if err := os.WriteFile("internal/cligen/generated/command.go", buf.Bytes(), 0o644); err != nil {
		panic(err)
	}
}
