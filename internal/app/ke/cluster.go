//go:generate mockgen -source=$GOFILE -package=mock$GOPACKAGE -destination=mock/$GOFILE
package ke

import (
	"context"
	"errors"
	"fmt"

	"github.com/nokamoto/2pf23/internal/app"
	"github.com/nokamoto/2pf23/internal/infra"
	"github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
)

type runtime interface {
	NewID() string
	Create(context.Context, *kev1alpha.Cluster) error
	Get(context.Context, string) (*kev1alpha.Cluster, error)
}

const (
	todoProject     = "unspecified"
	defaultNumNodes = 3
)

// Cluster is an application for managing Kubernetes clusters.
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
