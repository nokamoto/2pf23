package ke

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nokamoto/2pf23/internal/ent"
	"github.com/nokamoto/2pf23/internal/ent/enttest"
	"github.com/nokamoto/2pf23/internal/infra"
	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
	"github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestCluster_Create(t *testing.T) {
	ctx := context.TODO()

	// init
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
	defer client.Close()

	c := NewCluster(client)

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

	// delete
	err = c.Delete(ctx, want.Name)
	if err != nil {
		t.Fatalf("failed deleting cluster: %v", err)
	}

	// already deleted
	err = c.Delete(ctx, want.Name)
	if !errors.Is(err, infra.ErrNotFound) {
		t.Fatalf("expected ErrNotFound but got %v", err)
	}
}

func TestCluster_List(t *testing.T) {
	initClusters := func(t *testing.T, c *ent.Client, init ...*kev1alpha.Cluster) {
		t.Helper()
		var builders []*ent.ClusterCreate
		for _, x := range init {
			builders = append(builders, c.Cluster.Create().SetName(x.Name).SetDisplayName(x.DisplayName).SetNumNodes(x.NumNodes))
		}
		_, err := c.Cluster.CreateBulk(builders...).Save(context.Background())
		if err != nil {
			t.Fatalf("failed initializing clusters: %v", err)
		}
	}

	cluster := func(i int) *kev1alpha.Cluster {
		return &kev1alpha.Cluster{
			Name:        fmt.Sprintf("test name %d", i),
			DisplayName: fmt.Sprintf("test display name %d", i),
			NumNodes:    int32(i),
		}
	}

	list := func(xs ...*kev1alpha.Cluster) []*kev1alpha.Cluster {
		return xs
	}

	tests := []struct {
		name     string
		init     []*kev1alpha.Cluster
		pageSize int32
		page     *v1.Pagination
		want     []*kev1alpha.Cluster
		wantNext *v1.Pagination
	}{
		{
			name:     "empty",
			pageSize: 1,
		},
		{
			name:     "first page",
			init:     list(cluster(0), cluster(1), cluster(2)),
			pageSize: 1,
			want:     list(cluster(0)),
			wantNext: &v1.Pagination{
				Cursor: 2,
			},
		},
		{
			name:     "second page",
			init:     list(cluster(0), cluster(1), cluster(2)),
			pageSize: 1,
			page: &v1.Pagination{
				Cursor: 2,
			},
			want: list(cluster(1)),
			wantNext: &v1.Pagination{
				Cursor: 3,
			},
		},
		{
			name:     "last page",
			init:     list(cluster(0), cluster(1), cluster(2)),
			pageSize: 2,
			page: &v1.Pagination{
				Cursor: 2,
			},
			want: list(cluster(1), cluster(2)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()

			client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
			defer client.Close()

			c := NewCluster(client)
			initClusters(t, client, tt.init...)

			got, next, err := c.List(ctx, tt.pageSize, tt.page)
			if err != nil {
				t.Fatalf("failed listing clusters: %v", err)
			}
			if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantNext, next, protocmp.Transform()); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
		})
	}
}
