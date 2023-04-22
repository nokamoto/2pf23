package cligen

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func diff(t *testing.T, in, out string) {
	t.Helper()

	expected, err := os.ReadFile(in)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(string(expected), string(actual)); diff != "" {
		t.Errorf("%s %s (-want +got)\n%s", in, out, diff)
	}
}

func TestWalk_Walk(t *testing.T) {
	temp, err := os.MkdirTemp("", "cligen")
	if err != nil {
		t.Fatal(err)
	}

	w, err := NewWalk("../../testdata/cligen/generated.json", temp, "github.com/nokamoto/2pf23/internal/cligen/generated")
	if err != nil {
		t.Fatal(err)
	}

	if err := w.Walk(); err != nil {
		t.Fatal(err)
	}

	// diff(t, "generated/createcluster.go", fmt.Sprintf("%s/createcluster.go", temp))
	diff(t, "generated/root.go", fmt.Sprintf("%s/root.go", temp))
}
