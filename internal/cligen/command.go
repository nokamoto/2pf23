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
	Runtime runtimeArg
	Imports []*v1.ImportPath
	Command *v1.Command
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
		s += fmt.Sprintf("%s},", tabs)
	}
	return s
}

func initTemplate(v **template.Template, name string) error {
	if *v == nil {
		fm := template.FuncMap{
			"ToTitle":  cases.Title(language.English, cases.NoLower).String,
			"ToFields": printFields,
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
	return p.main.Execute(out, commandArg{
		Runtime: newRuntimeArg(),
		Imports: newImports(imports...),
		Command: cmd,
	})
}
