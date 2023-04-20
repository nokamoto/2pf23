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
				Method: "CreateCluster",
			},
		},
		SubPackages: []*v1.Package{
			{
				Package:    "sub",
				ImportPath: "github.com/nokamoto/2pf23/testdata/cligen/sub",
				Use:        "sub",
				Short:      "short",
				Long:       "long",
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
		Api:        "ke",
		ApiVersion: "v1alpha",
		ApiImportPath: &v1.ImportPath{
			Alias: "v1alpha",
			Path:  "github.com/nokamoto/2pf23/pkg/api/ke/v1alpha",
		},
		Package:          "generated",
		Use:              "use",
		Short:            "short",
		Long:             "long",
		Method:           "CreateCluster",
		MethodType:       v1.MethodType_METHOD_TYPE_CREATE,
		CreateResourceId: "Cluster",
		CreateResource: &v1.Resource{
			Type: "v1alpha.Cluster",
			Fields: []*v1.ResourceField{
				{
					Id:       "DisplayName",
					FlagName: "stringFlag",
				},
			},
		},
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
