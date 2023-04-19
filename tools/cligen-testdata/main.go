package main

import (
	"bytes"
	"io"
	"os"

	"github.com/nokamoto/2pf23/internal/cligen"
	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
)

func write[T any](file string, v T, f func(io.Writer, T) error) {
	var buf bytes.Buffer
	if err := f(&buf, v); err != nil {
		panic(err)
	}
	if err := os.WriteFile(file, buf.Bytes(), 0o644); err != nil {
		panic(err)
	}
}

// cligen-testdata generates testdata/cligen go files.
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
	write("testdata/cligen/main.go", cmd, p.PrintCommand)

	sub := &v1.Package{
		Package:    "sub",
		ImportPath: "github.com/nokamoto/2pf23/testdata/cligen/sub",
		Use:        "sub",
		Short:      "short",
		Long:       "long",
	}
	write("testdata/cligen/sub/root.go", sub, p.PrintRoot)

	pkg := &v1.Package{
		Package:     "cligen",
		Use:         "testdata",
		Short:       "short",
		Long:        "long",
		SubCommands: []*v1.Command{cmd},
		SubPackages: []*v1.Package{sub},
	}
	write("testdata/cligen/root.go", pkg, p.PrintRoot)
}
