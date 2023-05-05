// Code generated by server-gen. DO NOT EDIT.
package v1alpha

import (
	"context"
)

import (
	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/nokamoto/2pf23/internal/server/helper"
	v1alpha "github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
	"go.uber.org/zap"
)

type runtime interface {
	Create(ctx context.Context, resource *v1alpha.Cluster) (*v1alpha.Cluster, error)
	Get(ctx context.Context, name string) (*v1alpha.Cluster, error)
	Delete(ctx context.Context, name string) (*empty.Empty, error)
}

type service struct {
	v1alpha.UnimplementedKeServiceServer
	logger *zap.Logger
	rt     runtime
}

func NewService(logger *zap.Logger, rt runtime) *service {
	return &service{
		logger: logger.Named("ke.v1alpha"),
		rt:     rt,
	}
}

func (s *service) CreateCluster(ctx context.Context, req *v1alpha.CreateClusterRequest) (*v1alpha.Cluster, error) {
	logger := s.logger.With(zap.String("method", "CreateCluster"), zap.Any("request", req))
	logger.Debug("request received")
	res, err := s.rt.Create(ctx, req.GetCluster())
	return helper.ErrorOr(logger, res, err)
}

func (s *service) GetCluster(ctx context.Context, req *v1alpha.GetClusterRequest) (*v1alpha.Cluster, error) {
	logger := s.logger.With(zap.String("method", "GetCluster"), zap.Any("request", req))
	logger.Debug("request received")
	res, err := s.rt.Get(ctx, req.GetName())
	return helper.ErrorOr(logger, res, err)
}

func (s *service) DeleteCluster(ctx context.Context, req *v1alpha.DeleteClusterRequest) (*empty.Empty, error) {
	logger := s.logger.With(zap.String("method", "DeleteCluster"), zap.Any("request", req))
	logger.Debug("request received")
	res, err := s.rt.Delete(ctx, req.GetName())
	return helper.ErrorOr(logger, res, err)
}
