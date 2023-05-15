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
	main       *template.Template
	enableMock bool
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
	Service    *v1.Service
	Imports    []*v1.ImportPath
	EnableMock bool
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
		"github.com/nokamoto/2pf23/internal/server/helper",
	}
	for _, s := range static {
		imports = append(imports, &v1.ImportPath{
			Path: s,
		})
	}

	var delete, list, update bool
	for _, call := range svc.GetCalls() {
		if call.GetMethodType() == v1.MethodType_METHOD_TYPE_DELETE {
			delete = true
		}
		if call.GetMethodType() == v1.MethodType_METHOD_TYPE_LIST {
			list = true
		}
		if call.GetMethodType() == v1.MethodType_METHOD_TYPE_UPDATE {
			update = true
		}
	}
	if delete {
		imports = append(imports, &v1.ImportPath{
			Alias: "empty",
			Path:  "github.com/golang/protobuf/ptypes/empty",
		})
	}
	if list {
		imports = append(imports, &v1.ImportPath{
			Alias: "v1",
			Path:  "github.com/nokamoto/2pf23/pkg/api/inhouse/v1",
		})
	}
	if update {
		imports = append(imports, &v1.ImportPath{
			Path: "google.golang.org/protobuf/types/known/fieldmaskpb",
		})
	}

	sort.Slice(imports, func(i, j int) bool {
		return imports[i].Path < imports[j].Path
	})

	return p.main.Execute(out, mainArgs{
		Service:    svc,
		Imports:    imports,
		EnableMock: p.enableMock,
	})
}

// EnableMock enables mock generation.
// If enable is true, "//go:generate mockgen" is inserted into the generated code.
func (p *Printer) EnableMock(enable bool) {
	p.enableMock = enable
}
