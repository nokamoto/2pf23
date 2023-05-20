package helper

import (
	"context"

	"github.com/bufbuild/connect-go"
	v1alpha "github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
	"google.golang.org/protobuf/types/known/emptypb"
)

// workaround for generating gomock with generics
// failed parsing arguments: don't know how to parse type *ast.IndexExpr
type (
	CreateClusterRequest = connect.Request[v1alpha.CreateClusterRequest]
	Cluster              = connect.Response[v1alpha.Cluster]
	GetClusterRequest    = connect.Request[v1alpha.GetClusterRequest]
	DeleteClusterRequest = connect.Request[v1alpha.DeleteClusterRequest]
	Empty                = connect.Response[emptypb.Empty]
	ListClusterRequest   = connect.Request[v1alpha.ListClusterRequest]
	ListClusterResponse  = connect.Response[v1alpha.ListClusterResponse]
	UpdateClusterRequest = connect.Request[v1alpha.UpdateClusterRequest]
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -source=$GOFILE -package=mock$GOPACKAGE -destination=mock/$GOFILE
type KeServiceClient interface {
	CreateCluster(context.Context, *CreateClusterRequest) (*Cluster, error)
	GetCluster(context.Context, *GetClusterRequest) (*Cluster, error)
	DeleteCluster(context.Context, *DeleteClusterRequest) (*Empty, error)
	ListCluster(context.Context, *ListClusterRequest) (*ListClusterResponse, error)
	UpdateCluster(context.Context, *UpdateClusterRequest) (*Cluster, error)
}
