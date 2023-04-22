package cligen

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"

	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

// Walk is a command line interface generator.
// It generates a command line interface from protogen result.
type Walk struct {
	pkg           *v1.Package
	rootDirectory string
}

func NewWalk(file string, rootDirectory string) (*Walk, error) {
	var pkg v1.Package
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	if err := protojson.Unmarshal(b, &pkg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}
	return &Walk{
		pkg:           &pkg,
		rootDirectory: rootDirectory,
	}, nil
}

func (w *Walk) Walk() error {
	p := Printer{}
	file := path.Join(w.rootDirectory, "root.go")
	var b bytes.Buffer
	if err := p.PrintRoot(&b, w.pkg); err != nil {
		return fmt.Errorf("failed to print %s: %w", file, err)
	}
	if err := os.WriteFile(file, b.Bytes(), 0o644); err != nil {
		return fmt.Errorf("failed to write %s: %w", file, err)
	}

	for _, cmd := range w.pkg.SubCommands {
		file := path.Join(w.rootDirectory, strings.ToLower(cmd.GetMethod())+".go")
		var b bytes.Buffer
		if err := p.PrintCommand(&b, cmd); err != nil {
			return fmt.Errorf("failed to print %s: %w", file, err)
		}
		if err := os.WriteFile(file, b.Bytes(), 0o644); err != nil {
			return fmt.Errorf("failed to write %s: %w", file, err)
		}
	}

	return nil
}
