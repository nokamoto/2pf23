package ke

import (
	"context"
	"fmt"

	"github.com/nokamoto/2pf23/internal/ent"
	"github.com/nokamoto/2pf23/internal/ent/cluster"
	"github.com/nokamoto/2pf23/internal/infra"
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

func (c *Cluster) Create(ctx context.Context, cluster *kev1alpha.Cluster) error {
	_, err := c.client.Create().
		SetName(cluster.GetName()).
		SetDisplayName(cluster.GetDisplayName()).
		SetNumNodes(int(cluster.GetNumNodes())).
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
	return &kev1alpha.Cluster{
		Name:        res.Name,
		DisplayName: res.DisplayName,
		NumNodes:    int32(res.NumNodes),
	}, nil
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
