package cli

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/nokamoto/2pf23/internal/util/gen"
	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

// Walk is a command line interface generator.
// It generates a command line interface from protogen result.
type Walk struct {
	pkg           *v1.Package
	rootDirectory string
	rootPackage   string
}

// NewWalk creates a new Walk.
//
// file is a protogen result file. It contains protojson encoded v1.Package.
//
// rootDirectory is a root directory where the generated files are placed.
// e.g. internal/cli/generated
//
// rootPackage is a root package name of the generated root file.
// e.g. github.com/nokamoto/2pf23/internal/cli/generated
func NewWalk(file string, rootDirectory string, rootPackage string) (*Walk, error) {
	var pkg v1.Package
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	if err := protojson.Unmarshal(b, &pkg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}
	pkg.Use = "pf"
	pkg.Package = path.Base(rootPackage)
	pkg.Long = "pf is a CLI for managing the platform."
	return &Walk{
		pkg:           &pkg,
		rootDirectory: rootDirectory,
		rootPackage:   rootPackage,
	}, nil
}

func (w *Walk) walk(p *Printer, dir string, pkg *v1.Package, currentPackage string) error {
	err := os.MkdirAll(dir, 0o777)
	if err != nil {
		return fmt.Errorf("failed to mkdir %s: %w", dir, err)
	}

	file := path.Join(dir, "root.go")
	var b bytes.Buffer
	if err := p.PrintRoot(&b, pkg, currentPackage); err != nil {
		return fmt.Errorf("failed to print %s: %w", file, err)
	}
	if err := gen.WriteFormattedGo(file, b.Bytes()); err != nil {
		return fmt.Errorf("failed to write %s: %w", file, err)
	}

	for _, cmd := range pkg.GetSubCommands() {
		file := path.Join(dir, strings.ToLower(cmd.GetMethod())+".go")
		var b bytes.Buffer
		if err := p.PrintCommand(&b, cmd); err != nil {
			return fmt.Errorf("failed to print %s: %w", file, err)
		}
		if err := gen.WriteFormattedGo(file, b.Bytes()); err != nil {
			return fmt.Errorf("failed to write %s: %w", file, err)
		}
	}

	for _, sub := range pkg.GetSubPackages() {
		subDir := path.Join(dir, sub.GetPackage())
		if err := w.walk(p, subDir, sub, path.Join(currentPackage, sub.GetPackage())); err != nil {
			return fmt.Errorf("failed to walk %s: %w", subDir, err)
		}
	}

	return nil
}

func (w *Walk) Walk() error {
	return w.walk(&Printer{}, w.rootDirectory, w.pkg, w.rootPackage)
}
