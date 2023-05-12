package cligen

import (
	"embed"
	"fmt"
	"io"
	"path"
	"sort"
	"strings"
	"text/template"

	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

//go:embed templates
var f embed.FS

type Printer struct {
	root *template.Template
	main *template.Template
}

type runtimeArg struct {
	Name string
	Type string
}

func newRuntimeArg() runtimeArg {
	return runtimeArg{
		Name: "rt",
		Type: "runtime.Runtime",
	}
}

func newImports(imports ...*v1.ImportPath) []*v1.ImportPath {
	imports = append(
		imports,
		&v1.ImportPath{
			Path: "github.com/nokamoto/2pf23/internal/cli/runtime",
		},
		&v1.ImportPath{
			Path: "github.com/spf13/cobra",
		},
	)
	sort.Slice(imports, func(i, j int) bool {
		return imports[i].Path < imports[j].Path
	})
	return imports
}

type commandArg struct {
	Runtime    runtimeArg
	Imports    []*v1.ImportPath
	Command    *v1.Command
	ExactArgs0 bool
	ExactArgs1 bool
}

type rootArg struct {
	Runtime runtimeArg
	Imports []*v1.ImportPath
	Package *v1.Package
}

func printFields(indent int, req *v1.RequestMessage) string {
	var s string
	tabs := strings.Repeat("\t", indent)

	max := 0
	for _, field := range req.GetFields() {
		if len(field.GetName()) > max {
			max = len(field.GetName())
		}
	}
	for _, field := range req.GetFields() {
		spaces := strings.Repeat(" ", 1+max-len(field.GetName()))
		s += fmt.Sprintf("%s%s:%s%s,\n", tabs, field.GetName(), spaces, field.GetValue())
	}

	for _, child := range req.GetChildren() {
		s += fmt.Sprintf("%s%s: &%s{\n", tabs, child.GetName(), child.GetType())
		s += printFields(indent+1, child)
		s += fmt.Sprintf("%s},\n", tabs)
	}
	return s
}

func toValue(indent int, req *v1.RequestMessage) string {
	s := fmt.Sprintf("&%s{", req.GetType())
	if len(req.GetFields()) == 0 && len(req.GetChildren()) == 0 {
		return s + "}"
	}
	s += "\n"
	s += printFields(indent+1, req)
	s += strings.Repeat("\t", indent)
	s += "}"
	return s
}

func initTemplate(v **template.Template, name string) error {
	if *v == nil {
		fm := template.FuncMap{
			"ToTitle":  cases.Title(language.English, cases.NoLower).String,
			"ToFields": printFields,
			"ToValue":  toValue,
		}
		t, err := template.New(path.Base(name)).Funcs(fm).ParseFS(f, name)
		if err != nil {
			return err
		}
		*v = t
	}
	return nil
}

func (p *Printer) PrintRoot(out io.Writer, pkg *v1.Package, currentPackage string) error {
	if err := initTemplate(&p.root, "templates/root.go.tmpl"); err != nil {
		return err
	}
	var imports []*v1.ImportPath
	for _, sub := range pkg.GetSubPackages() {
		imports = append(imports, &v1.ImportPath{
			Path: path.Join(currentPackage, sub.GetPackage()),
		})
	}
	return p.root.Execute(out, rootArg{
		Runtime: newRuntimeArg(),
		Imports: newImports(imports...),
		Package: pkg,
	})
}

func (p *Printer) PrintCommand(out io.Writer, cmd *v1.Command) error {
	if err := initTemplate(&p.main, "templates/main.go.tmpl"); err != nil {
		return err
	}
	var imports []*v1.ImportPath
	if cmd.ApiImportPath != nil {
		imports = append(imports, cmd.ApiImportPath)
	}
	imports = append(imports, &v1.ImportPath{
		Path: "google.golang.org/protobuf/encoding/protojson",
	})

	if cmd.GetMethodType() == v1.MethodType_METHOD_TYPE_LIST {
		imports = append(imports, &v1.ImportPath{
			Alias: "helper",
			Path:  "github.com/nokamoto/2pf23/internal/cli/helper",
		})
	}
	if cmd.GetMethodType() == v1.MethodType_METHOD_TYPE_UPDATE {
		imports = append(imports, &v1.ImportPath{
			Path: "google.golang.org/protobuf/types/known/fieldmaskpb",
		})
	}

	var arg0 bool
	if typ := cmd.GetMethodType(); typ == v1.MethodType_METHOD_TYPE_LIST || typ == v1.MethodType_METHOD_TYPE_CREATE {
		arg0 = true
	}

	return p.main.Execute(out, commandArg{
		Runtime:    newRuntimeArg(),
		Imports:    newImports(imports...),
		Command:    cmd,
		ExactArgs0: arg0,
		ExactArgs1: !arg0,
	})
}
