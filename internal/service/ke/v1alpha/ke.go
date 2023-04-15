package kev1alpha

import (
	"context"

	kev1alpha "github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type service struct {
	kev1alpha.UnimplementedKeServiceServer
	logger *zap.Logger
}

// NewService returns a new KeServiceServer.
func NewService(logger *zap.Logger) *service {
	return &service{
		logger: logger.Named("ke.v1alpha"),
	}
}

func (s *service) CreateCluster(ctx context.Context, req *kev1alpha.CreateClusterRequest) (*kev1alpha.Cluster, error) {
	logger := s.logger.With(zap.String("method", "CreateCluster"), zap.Any("request", req))
	logger.Debug("request received")
	return nil, status.Errorf(codes.Unimplemented, "method CreateCluster not implemented")
}
