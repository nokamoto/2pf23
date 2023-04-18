package cligen

import (
	"bytes"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
)

func TestPrinter_PrintCommand(t *testing.T) {
	var buf bytes.Buffer
	p := &Printer{out: &buf}
	cmd := v1.Command{
		Package: "cligen",
		Use:     "use",
		Short:   "short",
		Long:    "long",
		Method:  "Create",
	}
	if err := p.PrintCommand(&cmd); err != nil {
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
