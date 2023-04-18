package main

import (
	"bytes"
	"os"

	"github.com/nokamoto/2pf23/internal/cligen"
	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
)

// cligen-testdata generates testdata/cligen/main.go.
func main() {
	var buf bytes.Buffer

	cmd := &v1.Command{
		Package: "cligen",
		Use:     "use",
		Short:   "short",
		Long:    "long",
		Method:  "Create",
		StringFlags: []*v1.Flag{
			{
				Name:        "stringFlag",
				DisplayName: "string-flag",
				Value:       "value",
				Usage:       "usage",
			},
		},
	}
	p := cligen.Printer{}
	if err := p.PrintCommand(&buf, cmd); err != nil {
		panic(err)
	}
	if err := os.WriteFile("testdata/cligen/main.go", buf.Bytes(), 0o644); err != nil {
		panic(err)
	}

	pkg := &v1.Package{
		Package:     "cligen",
		Use:         "testdata",
		Short:       "short",
		Long:        "long",
		SubCommands: []*v1.Command{cmd},
	}
	buf.Reset()
	if err := p.PrintRoot(&buf, pkg); err != nil {
		panic(err)
	}
	if err := os.WriteFile("testdata/cligen/root.go", buf.Bytes(), 0o644); err != nil {
		panic(err)
	}
}
