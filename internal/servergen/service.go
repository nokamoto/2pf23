package servergen

import (
	"embed"
	"io"
	"text/template"
)

//go:embed templates
var f embed.FS

type Printer struct {
	main *template.Template
}

func PrintService(out io.Writer) {}
