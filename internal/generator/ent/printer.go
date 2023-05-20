package ent

import (
	"embed"
	"fmt"
	"io"
	"path"
	"text/template"

	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
)

//go:embed templates
var f embed.FS

type Printer struct {
	query *template.Template
}


func newSetters(ent *v1.Ent) []string{
	var setters []string
	size := len(ent.GetFields()) + len(ent.GetEnumFields())
	for i, f := range ent.GetFields() {
		var dot string
		if i != size - 1 {
			dot = "."
		}
		setters = append(setters, fmt.Sprintf("Set%s(v.Get%s())%s", f, f, dot))
	}
	for i, f := range ent.GetEnumFields() {
		var dot string
		if len(ent.GetFields())+i != size - 1 {
			dot = "."
		}
		setters = append(setters, fmt.Sprintf("Set%s(int32(v.Get%s()))%s", f, f, dot))
	}
	return setters
}

type args struct {
	Ent *v1.Ent
	Package string
	Imports []*v1.ImportPath
	Setters []string
}

func initTemplate(v **template.Template, name string) error {
	if *v == nil {
		fm := template.FuncMap{}
		t, err := template.New(path.Base(name)).Funcs(fm).ParseFS(f, name)
		if err != nil {
			return err
		}
		*v = t
	}
	return nil
}

func (p *Printer) PrintQuery(out io.Writer, ent *v1.Ent, pkg string) error {
	if err := initTemplate(&p.query, "templates/query.go.tmpl"); err != nil {
		return err
	}

	var imports []*v1.ImportPath
	imports = append(imports, ent.GetImportPath())
	imports = append(imports, &v1.ImportPath{
		Alias: "ent",
		Path: "github.com/nokamoto/2pf23/internal/ent",
	})

	return p.query.Execute(out, args{
		Ent: ent,
		Package: pkg,
		Imports: imports,
		Setters: newSetters(ent),
	})
}
