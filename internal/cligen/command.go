package cligen

import (
	"embed"
	"io"
	"path"
	"sort"
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

type importArg struct {
	Alias string
	Path  string
}

func newImports(imports []importArg) []importArg {
	imports = append(
		imports,
		importArg{
			Path: "github.com/nokamoto/2pf23/internal/cli/runtime",
		},
		importArg{
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
	Imports []importArg
	Command *v1.Command
}

type rootArg struct {
	Runtime runtimeArg
	Imports []importArg
	Package *v1.Package
}

func initTemplate(v **template.Template, name string) error {
	if *v == nil {
		fm := template.FuncMap{
			"ToTitle": cases.Title(language.English, cases.NoLower).String,
		}
		t, err := template.New(path.Base(name)).Funcs(fm).ParseFS(f, name)
		if err != nil {
			return err
		}
		*v = t
	}
	return nil
}

func (p *Printer) PrintRoot(out io.Writer, pkg *v1.Package) error {
	if err := initTemplate(&p.root, "templates/root.go.tmpl"); err != nil {
		return err
	}
	return p.root.Execute(out, rootArg{
		Runtime: newRuntimeArg(),
		Imports: newImports(nil),
		Package: pkg,
	})
}

func (p *Printer) PrintCommand(out io.Writer, cmd *v1.Command) error {
	if err := initTemplate(&p.main, "templates/main.go.tmpl"); err != nil {
		return err
	}
	return p.main.Execute(out, commandArg{
		Runtime: newRuntimeArg(),
		Imports: newImports(nil),
		Command: cmd,
	})
}
