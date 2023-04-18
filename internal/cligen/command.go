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
	out  io.Writer
	main *template.Template
}

func (p *Printer) PrintCommand(cmd *v1.Command) error {
	if p.main == nil {
		main, err := template.ParseFS(f, "templates/main.go.tmpl")
		if err != nil {
			return err
		}
		p.main = main
	}

	return p.main.Execute(p.out, cmd)
}
