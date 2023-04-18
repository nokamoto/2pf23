package cligen

import (
	"bytes"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
)

func TestPrinter_PrintRoot(t *testing.T) {
	var buf bytes.Buffer
	p := Printer{}
	pkg := &v1.Package{
		Package: "cligen",
		Use:     "testdata",
		Short:   "short",
		Long:    "long",
		SubCommands: []*v1.Command{
			{
				Method: "Create",
			},
		},
	}
	if err := p.PrintRoot(&buf, pkg); err != nil {
		t.Fatal(err)
	}

	expected, err := os.ReadFile("../../testdata/cligen/root.go")
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(string(expected), buf.String()); diff != "" {
		t.Fatal(diff)
	}
}

func TestPrinter_PrintCommand(t *testing.T) {
	var buf bytes.Buffer
	p := Printer{}
	cmd := v1.Command{
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
	if err := p.PrintCommand(&buf, &cmd); err != nil {
		t.Fatal(err)
	}

	expected, err := os.ReadFile("../../testdata/cligen/main.go")
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(string(expected), buf.String()); diff != "" {
		t.Fatal(diff)
	}
}
