package ke

import (
	"context"
	"fmt"

	"github.com/nokamoto/2pf23/internal/ent"
	"github.com/nokamoto/2pf23/internal/ent/cluster"
	"github.com/nokamoto/2pf23/internal/infra"
	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
	"github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
)

type Cluster struct {
	client *ent.ClusterClient
}

func NewCluster(client *ent.Client) *Cluster {
	return &Cluster{
		client: client.Cluster,
	}
}

func (c *Cluster) proto(x *ent.Cluster) *kev1alpha.Cluster {
	return &kev1alpha.Cluster{
		Name:        x.Name,
		DisplayName: x.DisplayName,
		NumNodes:    x.NumNodes,
	}
}

func (c *Cluster) Create(ctx context.Context, cluster *kev1alpha.Cluster) error {
	_, err := c.client.Create().
		SetName(cluster.GetName()).
		SetDisplayName(cluster.GetDisplayName()).
		SetNumNodes(cluster.GetNumNodes()).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed saving cluster: %w", err)
	}
	return nil
}

// Get returns a cluster by name.
// If the cluster does not exist, it returns infra.ErrNotFound.
func (c *Cluster) Get(ctx context.Context, name string) (*kev1alpha.Cluster, error) {
	res, err := c.client.Query().Where(cluster.Name(name)).Only(ctx)
	if ent.IsNotFound(err) {
		return nil, infra.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return c.proto(res), nil
}

// Delete deletes a cluster by name.
// If the cluster does not exist, it returns infra.ErrNotFound.
func (c *Cluster) Delete(ctx context.Context, name string) error {
	res, err := c.client.Delete().Where(cluster.Name(name)).Exec(ctx)
	if res == 0 || ent.IsNotFound(err) {
		return infra.ErrNotFound
	}
	if err != nil {
		return err
	}
	return nil
}

// List returns a list of clusters and next page.
//
// If the page is nil, it returns the first page.
// If next page does not exist, it returns nil.
func (c *Cluster) List(ctx context.Context, pageSize int32, page *v1.Pagination) ([]*kev1alpha.Cluster, *v1.Pagination, error) {
	res, err := c.client.Query().Limit(int(pageSize + 1)).Where(cluster.IDGTE(page.GetCursor())).All(ctx)
	if err != nil {
		return nil, nil, err
	}
	var next *v1.Pagination
	if len(res) == int(pageSize)+1 {
		next = &v1.Pagination{
			Cursor: res[pageSize].ID,
		}
		res = res[:pageSize]
	}
	var clusters []*kev1alpha.Cluster
	for _, x := range res {
		clusters = append(clusters, c.proto(x))
	}
	return clusters, next, nil
}
