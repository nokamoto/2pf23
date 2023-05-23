package ke

import (
	"context"
	"fmt"

	"github.com/mennanov/fmutils"
	"github.com/nokamoto/2pf23/internal/ent"
	entcluster "github.com/nokamoto/2pf23/internal/ent/cluster"
	entproto "github.com/nokamoto/2pf23/internal/ent/proto"
	"github.com/nokamoto/2pf23/internal/infra"
	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
	"github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type Cluster struct {
	client *ent.Client
}

func NewCluster(client *ent.Client) *Cluster {
	return &Cluster{
		client: client,
	}
}

// Create creates a cluster.
// If the cluster already exists, it returns infra.ErrAlreadyExists.
func (c *Cluster) Create(ctx context.Context, cluster *kev1alpha.Cluster) (*kev1alpha.Cluster, error) {
	res, err := entproto.ClusterCreateQuery(c.client.Cluster.Create(), cluster).Save(ctx)
	if ent.IsConstraintError(err) {
		return nil, fmt.Errorf("%w: %s", infra.ErrAlreadyExists, cluster.GetName())
	}
	if err != nil {
		return nil, fmt.Errorf("failed saving cluster: %w", err)
	}
	return entproto.ClusterProto(res), nil
}

// Get returns a cluster by name.
// If the cluster does not exist, it returns infra.ErrNotFound.
func (c *Cluster) Get(ctx context.Context, name string) (*kev1alpha.Cluster, error) {
	res, err := c.client.Cluster.Query().Where(entcluster.Name(name)).Only(ctx)
	if ent.IsNotFound(err) {
		return nil, infra.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return entproto.ClusterProto(res), nil
}

// Delete deletes a cluster by name.
// If the cluster does not exist, it returns infra.ErrNotFound.
func (c *Cluster) Delete(ctx context.Context, name string) error {
	res, err := c.client.Cluster.Delete().Where(entcluster.Name(name)).Exec(ctx)
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
	res, err := c.client.Cluster.Query().Limit(int(pageSize + 1)).Where(entcluster.IDGTE(page.GetCursor())).All(ctx)
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
		clusters = append(clusters, entproto.ClusterProto(x))
	}
	return clusters, next, nil
}

// Update updates a cluster.
//
// If the cluster does not exist, it returns infra.ErrNotFound.
// If the mask is invalid, it returns infra.ErrInvalidArgument.
func (c *Cluster) Update(ctx context.Context, cluster *kev1alpha.Cluster, mask *fieldmaskpb.FieldMask) (*kev1alpha.Cluster, error) {
	name := cluster.GetName()
	rollback := func(tx *ent.Tx, err error) error {
		if rerr := tx.Rollback(); rerr != nil {
			err = fmt.Errorf("%w: with rollback error: %v", err, rerr)
		}
		return err
	}
	mask.Normalize()
	if !mask.IsValid(cluster) {
		return nil, fmt.Errorf("%w: invalid mask: %v", infra.ErrInvalidArgument, mask)
	}
	fmutils.Filter(cluster, mask.GetPaths())
	tx, err := c.client.Tx(ctx)
	if err != nil {
		return nil, err
	}
	got, err := tx.Cluster.Query().Where(entcluster.Name(name)).Only(ctx)
	if ent.IsNotFound(err) {
		return nil, rollback(tx, fmt.Errorf("%w: %s", infra.ErrNotFound, cluster.GetName()))
	}
	if err != nil {
		return nil, rollback(tx, err)
	}
	updated := entproto.ClusterProto(got)
	proto.Merge(updated, cluster)
	got, err = entproto.ClusterUpdateOneQuery(tx.Cluster.UpdateOneID(got.ID), updated).Save(ctx)
	if err != nil {
		return nil, rollback(tx, err)
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return entproto.ClusterProto(got), nil
}
