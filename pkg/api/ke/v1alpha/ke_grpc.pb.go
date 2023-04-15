// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: api/ke/v1alpha/ke.proto

package kev1alpha

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	KeService_CreateCluster_FullMethodName = "/api.ke.v1alpha.KeService/CreateCluster"
)

// KeServiceClient is the client API for KeService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KeServiceClient interface {
	CreateCluster(ctx context.Context, in *CreateClusterRequest, opts ...grpc.CallOption) (*Cluster, error)
}

type keServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewKeServiceClient(cc grpc.ClientConnInterface) KeServiceClient {
	return &keServiceClient{cc}
}

func (c *keServiceClient) CreateCluster(ctx context.Context, in *CreateClusterRequest, opts ...grpc.CallOption) (*Cluster, error) {
	out := new(Cluster)
	err := c.cc.Invoke(ctx, KeService_CreateCluster_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KeServiceServer is the server API for KeService service.
// All implementations must embed UnimplementedKeServiceServer
// for forward compatibility
type KeServiceServer interface {
	CreateCluster(context.Context, *CreateClusterRequest) (*Cluster, error)
	mustEmbedUnimplementedKeServiceServer()
}

// UnimplementedKeServiceServer must be embedded to have forward compatible implementations.
type UnimplementedKeServiceServer struct {
}

func (UnimplementedKeServiceServer) CreateCluster(context.Context, *CreateClusterRequest) (*Cluster, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCluster not implemented")
}
func (UnimplementedKeServiceServer) mustEmbedUnimplementedKeServiceServer() {}

// UnsafeKeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KeServiceServer will
// result in compilation errors.
type UnsafeKeServiceServer interface {
	mustEmbedUnimplementedKeServiceServer()
}

func RegisterKeServiceServer(s grpc.ServiceRegistrar, srv KeServiceServer) {
	s.RegisterService(&KeService_ServiceDesc, srv)
}

func _KeService_CreateCluster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateClusterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeServiceServer).CreateCluster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KeService_CreateCluster_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeServiceServer).CreateCluster(ctx, req.(*CreateClusterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// KeService_ServiceDesc is the grpc.ServiceDesc for KeService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var KeService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.ke.v1alpha.KeService",
	HandlerType: (*KeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateCluster",
			Handler:    _KeService_CreateCluster_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/ke/v1alpha/ke.proto",
}
