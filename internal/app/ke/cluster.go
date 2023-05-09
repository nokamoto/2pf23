//go:generate go run -mod=mod github.com/golang/mock/mockgen -source=$GOFILE -package=mock$GOPACKAGE -destination=mock/$GOFILE
package ke

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/nokamoto/2pf23/internal/app"
	"github.com/nokamoto/2pf23/internal/infra"
	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
	"github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
)

type runtime interface {
	NewID() string
	Create(context.Context, *kev1alpha.Cluster) error
	Get(context.Context, string) (*kev1alpha.Cluster, error)
	Delete(context.Context, string) error
	List(context.Context, int32, *v1.Pagination) ([]*kev1alpha.Cluster, *v1.Pagination, error)
}

const (
	todoProject     = "unspecified"
	defaultNumNodes = 3
)

// Cluster is an application for managing Kubernetes clusters.
//
// All methods of Cluster returns errors defined in internal/app/errors.go, such as app.ErrNotFound. Or it returns errors as is if unknown.
type Cluster struct {
	rt runtime
}

func NewCluster(rt runtime) *Cluster {
	return &Cluster{rt: rt}
}

func (c *Cluster) generateName() string {
	return fmt.Sprintf("projects/%s/clusters/%s", todoProject, c.rt.NewID())
}

// Create creates a new cluster.
//
// The unique name of the cluster is generated automatically by the application. The name is returned in the response.
func (c *Cluster) Create(ctx context.Context, cluster *kev1alpha.Cluster) (*kev1alpha.Cluster, error) {
	if cluster.GetNumNodes() == 0 {
		cluster.NumNodes = defaultNumNodes
	}
	if cluster.GetNumNodes() > 5 {
		return nil, fmt.Errorf("%w: num_nodes must be less than or equal to 5", app.ErrInvalidArgument)
	}

	cluster.Name = c.generateName()

	err := c.rt.Create(ctx, cluster)
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

// Get returns a cluster by name.
func (c *Cluster) Get(ctx context.Context, name string) (*kev1alpha.Cluster, error) {
	res, err := c.rt.Get(ctx, name)
	if errors.Is(err, infra.ErrNotFound) {
		return nil, fmt.Errorf("%w: %s", app.ErrNotFound, name)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Delete deletes a cluster by name.
//
// If the cluster does not exist, it returns app.ErrNotFound.
func (c *Cluster) Delete(ctx context.Context, name string) (*empty.Empty, error) {
	err := c.rt.Delete(ctx, name)
	if errors.Is(err, infra.ErrNotFound) {
		return nil, fmt.Errorf("%w: %s", app.ErrNotFound, name)
	}
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

// List returns a list of clusters and next page.
//
// If pageSize is 0, it returns 30 clusters.
func (c *Cluster) List(ctx context.Context, pageSize int32, page *v1.Pagination) ([]*kev1alpha.Cluster, *v1.Pagination, error) {
	if pageSize == 0 {
		pageSize = 30
	}
	if pageSize < 0 {
		return nil, nil, fmt.Errorf("%w: pageSize must be greater than or equal to 0", app.ErrInvalidArgument)
	}

	res, next, err := c.rt.List(ctx, pageSize, page)
	if err != nil {
		return nil, nil, err
	}
	return res, next, nil
}
