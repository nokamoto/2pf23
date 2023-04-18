package cligen

import (
	"embed"
	"io"
	"text/template"

	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
)

//go:embed templates
var f embed.FS

type Printer struct {
	root *template.Template
	main *template.Template
}

func initTemplate(v **template.Template, name string) error {
	if *v == nil {
		t, err := template.ParseFS(f, name)
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
	return p.root.Execute(out, pkg)
}

func (p *Printer) PrintCommand(out io.Writer, cmd *v1.Command) error {
	if err := initTemplate(&p.main, "templates/main.go.tmpl"); err != nil {
		return err
	}
	return p.main.Execute(out, cmd)
}
