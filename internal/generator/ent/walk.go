package ent

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/nokamoto/2pf23/internal/util/gen"
	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

// Walk is a ent query generator.
// It generates a ent query from protogen result.
type Walk struct {
	ents      []*v1.Ent
	directory string
}

// NewWalk creates a new Walk.
//
// indir is a directory where the protogen result files are placed.
// It contains protojson encoded v1.Ent files.
//
// outdir is a directory where the generated files are placed.
func NewWalk(indir string, outdir string) (*Walk, error) {
	w := &Walk{
		directory: outdir,
	}
	err := filepath.WalkDir(indir, func(path string, d fs.DirEntry, err error) error {
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

		var ent v1.Ent
		if err := protojson.Unmarshal(b, &ent); err != nil {
			return fmt.Errorf("failed to unmarshal: %s: %w", path, err)
		}

		w.ents = append(w.ents, &ent)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk %s: %w", indir, err)
	}
	return w, nil
}

func (w *Walk) Walk() error {
	p := Printer{}
	for _, ent := range w.ents {
		base := filepath.Base(w.directory)

		var b bytes.Buffer
		if err := p.PrintQuery(&b, ent, base); err != nil {
			return fmt.Errorf("failed to print query: %w", err)
		}

		file := filepath.Join(w.directory, fmt.Sprintf("%s.go", strings.ToLower(ent.GetName())))
		if err := gen.WriteFormattedGo(file, b.Bytes()); err != nil {
			return fmt.Errorf("failed to write file: %s: %w", file, err)
		}
	}
	return nil
}
