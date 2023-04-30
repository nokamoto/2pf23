package servergen

import (
	"embed"
	"io"
	"path"
	"sort"
	"text/template"

	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
)

//go:embed templates
var f embed.FS

type Printer struct {
	main *template.Template
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

type mainArgs struct {
	Service *v1.Service
	Imports []*v1.ImportPath
}

func (p *Printer) PrintService(out io.Writer, svc *v1.Service) error {
	if err := initTemplate(&p.main, "templates/main.go.tmpl"); err != nil {
		return err
	}

	var imports []*v1.ImportPath
	if svc.ApiImportPath != nil {
		imports = append(imports, svc.ApiImportPath)
	}

	static := []string{
		"go.uber.org/zap",
		"google.golang.org/grpc/codes",
		"google.golang.org/grpc/status",
	}
	for _, s := range static {
		imports = append(imports, &v1.ImportPath{
			Path: s,
		})
	}

	sort.Slice(imports, func(i, j int) bool {
		return imports[i].Path < imports[j].Path
	})

	return p.main.Execute(out, mainArgs{
		Service: svc,
		Imports: imports,
	})
}
