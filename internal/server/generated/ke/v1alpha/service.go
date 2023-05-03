// Code generated by server-gen. DO NOT EDIT.
package v1alpha

import (
	"context"
	"errors"
)

import (
	"github.com/nokamoto/2pf23/internal/app"
	v1alpha "github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type runtime interface {
	Create(ctx context.Context, resource *v1alpha.Cluster) (*v1alpha.Cluster, error)
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
	// standard create method
	res, err := s.rt.Create(ctx, req.GetCluster())
	if errors.Is(err, app.ErrInvalidArgument) {
		logger.Error("invalid argument", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err != nil {
		logger.Error("unknown error", zap.Error(err))
		return nil, status.Error(codes.Unknown, "unknown error")
	}
	return res, nil
}
