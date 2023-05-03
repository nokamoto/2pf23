package ke

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nokamoto/2pf23/internal/ent/enttest"
	"github.com/nokamoto/2pf23/internal/infra"
	"github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestCluster_Create(t *testing.T) {
	ctx := context.TODO()

	// init
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
	defer client.Close()

	c := &Cluster{
		client: client.Cluster,
	}

	want := &kev1alpha.Cluster{
		Name:        "test name",
		DisplayName: "test display name",
		NumNodes:    3,
	}

	// not found
	_, err := c.Get(ctx, want.Name)
	if !errors.Is(err, infra.ErrNotFound) {
		t.Fatalf("expected ErrNotFound but got %v", err)
	}

	// create
	err = c.Create(ctx, want)
	if err != nil {
		t.Fatalf("failed creating cluster: %v", err)
	}

	// get
	got, err := c.Get(ctx, want.Name)
	if err != nil {
		t.Fatalf("failed getting cluster: %v", err)
	}
	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Errorf("(-want, +got)\n%s", diff)
	}
}
