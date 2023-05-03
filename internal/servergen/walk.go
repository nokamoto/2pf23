package servergen

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

// Walk is a server generator.
// It generates a server from protogen result.
type Walk struct {
	services      []*v1.Service
	rootDirectory string
	enableMock    bool
}

// NewWalk creates a new Walk.
//
// inputDirectory is a directory where the input files are placed. The input files contain protojson encoded v1.Package.
//
// rootDirectory is a root directory where the generated files are placed.
//
// enableMock enables mock generation. see servergen.Printer.EnableMock.
func NewWalk(inputDirectory string, rootDirectory string, enableMock bool) (*Walk, error) {
	walk := Walk{
		rootDirectory: rootDirectory,
		enableMock:    enableMock,
	}
	err := filepath.WalkDir(inputDirectory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		b, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file: %s: %w", path, err)
		}

		var svc v1.Service
		if err := protojson.Unmarshal(b, &svc); err != nil {
			return fmt.Errorf("failed to unmarshal: %s: %w", path, err)
		}

		walk.services = append(walk.services, &svc)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk %s: %w", inputDirectory, err)
	}
	return &walk, nil
}

func (w *Walk) walk(p *Printer, svc *v1.Service) error {
	dir := filepath.Join(w.rootDirectory, svc.GetName(), svc.GetApiVersion())
	if err := os.MkdirAll(dir, 0o777); err != nil {
		return fmt.Errorf("failed to mkdir %s: %w", dir, err)
	}

	var b bytes.Buffer
	if err := p.PrintService(&b, svc); err != nil {
		return fmt.Errorf("failed to print %s: %w", dir, err)
	}

	file := filepath.Join(dir, "service.go")
	if err := os.WriteFile(file, b.Bytes(), 0o644); err != nil {
		return fmt.Errorf("failed to write %s: %w", file, err)
	}

	return nil
}

func (w *Walk) Walk() error {
	p := &Printer{}
	p.EnableMock(w.enableMock)
	for _, svc := range w.services {
		if err := w.walk(p, svc); err != nil {
			return err
		}
	}
	return nil
}
