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
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func initClusters(t *testing.T, c *ent.Client, init ...*kev1alpha.Cluster) {
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

type testcase[T1 any, T2 any] struct {
	name string
	init []*kev1alpha.Cluster
	req  T1
	want T2
	err  error
}

func run[T1 any, T2 any](t *testing.T, tests []testcase[T1, T2], f func(*Cluster, context.Context, T1) (T2, error)) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()

			client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
			defer client.Close()

			c := NewCluster(client)
			initClusters(t, client, tt.init...)

			got, err := f(c, ctx, tt.req)
			if !errors.Is(err, tt.err) {
				t.Errorf("want %v, got %v", tt.err, err)
			}
			if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
		})
	}
}

func TestCluster_Create(t *testing.T) {
	tests := []testcase[*kev1alpha.Cluster, *kev1alpha.Cluster]{
		{
			name: "ok",
			req: &kev1alpha.Cluster{
				Name:        "test name",
				DisplayName: "test display name",
				NumNodes:    3,
			},
			want: &kev1alpha.Cluster{
				Name:        "test name",
				DisplayName: "test display name",
				NumNodes:    3,
			},
		},
		{
			name: "already exists",
			init: []*kev1alpha.Cluster{
				{
					Name: "test name",
				},
			},
			req: &kev1alpha.Cluster{
				Name:        "test name",
				DisplayName: "test display name",
				NumNodes:    3,
			},
			err: infra.ErrAlreadyExists,
		},
	}

	run(t, tests, (*Cluster).Create)
}

func TestCluster_Get(t *testing.T) {
	tests := []testcase[string, *kev1alpha.Cluster]{
		{
			name: "ok",
			init: []*kev1alpha.Cluster{
				{
					Name:        "test name",
					DisplayName: "test display name",
					NumNodes:    3,
				},
			},
			req: "test name",
			want: &kev1alpha.Cluster{
				Name:        "test name",
				DisplayName: "test display name",
				NumNodes:    3,
			},
		},
		{
			name: "not found",
			req:  "test name",
			err:  infra.ErrNotFound,
		},
	}

	run(t, tests, (*Cluster).Get)
}

func TestCluster_Delete(t *testing.T) {
	tests := []struct {
		name string
		init *kev1alpha.Cluster
		req  string
		want *kev1alpha.Cluster
		err  error
	}{
		{
			name: "ok",
			init: &kev1alpha.Cluster{
				Name: "test name",
			},
			req: "test name",
		},
		{
			name: "not found",
			req:  "test name",
			err:  infra.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()

			client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
			defer client.Close()

			c := NewCluster(client)
			if tt.init != nil {
				initClusters(t, client, tt.init)
			}

			err := c.Delete(ctx, tt.req)
			if !errors.Is(err, tt.err) {
				t.Errorf("want %v, got %v", tt.err, err)
			}

			got, _ := c.Get(ctx, tt.req)
			if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
		})
	}
}

func TestCluster_List(t *testing.T) {
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

func TestCluster_Update(t *testing.T) {
	tests := []struct {
		name string
		init []*kev1alpha.Cluster
		req  *kev1alpha.Cluster
		mask *fieldmaskpb.FieldMask
		want *kev1alpha.Cluster
		err  error
	}{
		{
			name: "ok",
			init: []*kev1alpha.Cluster{
				{
					Name:        "test name",
					DisplayName: "test display name",
					NumNodes:    1,
				},
			},
			req: &kev1alpha.Cluster{
				Name:        "test name",
				DisplayName: "test display name 2",
				NumNodes:    2,
			},
			mask: &fieldmaskpb.FieldMask{
				Paths: []string{"display_name", "num_nodes"},
			},
			want: &kev1alpha.Cluster{
				Name:        "test name",
				DisplayName: "test display name 2",
				NumNodes:    2,
			},
		},
		{
			name: "update part of fields",
			init: []*kev1alpha.Cluster{
				{
					Name:        "test name",
					DisplayName: "test display name",
					NumNodes:    1,
				},
			},
			req: &kev1alpha.Cluster{
				Name:        "test name",
				DisplayName: "test display name 2",
				NumNodes:    2,
			},
			mask: &fieldmaskpb.FieldMask{
				Paths: []string{"display_name"},
			},
			want: &kev1alpha.Cluster{
				Name:        "test name",
				DisplayName: "test display name 2",
				NumNodes:    1,
			},
		},
		{
			name: "not found",
			req: &kev1alpha.Cluster{
				Name:        "test name",
				DisplayName: "test display name 2",
				NumNodes:    2,
			},
			mask: &fieldmaskpb.FieldMask{
				Paths: []string{"display_name", "num_nodes"},
			},
			err: infra.ErrNotFound,
		},
		{
			name: "invalid mask",
			mask: &fieldmaskpb.FieldMask{
				Paths: []string{"invalid"},
			},
			err: infra.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()

			client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
			defer client.Close()

			c := NewCluster(client)
			initClusters(t, client, tt.init...)

			got, err := c.Update(ctx, tt.req, tt.mask)
			if !errors.Is(err, tt.err) {
				t.Errorf("want %v, got %v", tt.err, err)
			}

			t.Log(client.Cluster.Query().All(ctx))

			if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
		})
	}
}
