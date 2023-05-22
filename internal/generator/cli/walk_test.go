package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWalk_Walk(t *testing.T) {
	temp, err := os.MkdirTemp("", "cligen")
	if err != nil {
		t.Fatal(err)
	}

	w, err := NewWalk("../../../testdata/generator/cli/test.json", temp, "github.com/nokamoto/2pf23/internal/generator/cli/generated")
	if err != nil {
		t.Fatal(err)
	}

	if err := w.Walk(); err != nil {
		t.Fatal(err)
	}

	err = filepath.WalkDir("generated", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		target := filepath.Join(temp, strings.TrimPrefix(path, "generated"))

		if d.IsDir() && strings.Contains(d.Name(), "helper") {
			return filepath.SkipDir
		}
		if d.IsDir() || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		expected, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}

		actual, err := os.ReadFile(target)
		if err != nil {
			return fmt.Errorf("%s: %w", target, err)
		}

		if diff := cmp.Diff(string(expected), string(actual)); diff != "" {
			return fmt.Errorf("%s %s (-want +got)\n%s", path, target, diff)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
