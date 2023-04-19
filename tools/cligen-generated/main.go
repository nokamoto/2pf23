package main

import (
	"bytes"
	"os"

	"github.com/nokamoto/2pf23/internal/cligen"
	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
)

func write(file string, v *v1.Command, p cligen.Printer) {
	var buf bytes.Buffer
	if err := p.PrintCommand(&buf, v); err != nil {
		panic(err)
	}
	if err := os.WriteFile(file, buf.Bytes(), 0o644); err != nil {
		panic(err)
	}
}

// cligen-testdata generates internal/cligen/generated/command.go.
func main() {
	p := cligen.Printer{}

	cmd := &v1.Command{
		Api:        "ke",
		ApiVersion: "v1alpha",
		Package:    "generated",
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
	write("internal/cligen/generated/create.go", cmd, p)
}