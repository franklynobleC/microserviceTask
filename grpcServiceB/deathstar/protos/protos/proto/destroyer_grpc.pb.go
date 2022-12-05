// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: destroyer.proto

package destroyerService

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

// DestroyerServiceClient is the client API for DestroyerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DestroyerServiceClient interface {
	AcquireTarget(ctx context.Context, in *DestroyerRequest, opts ...grpc.CallOption) (*DestroyerResponse, error)
}

type destroyerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDestroyerServiceClient(cc grpc.ClientConnInterface) DestroyerServiceClient {
	return &destroyerServiceClient{cc}
}

func (c *destroyerServiceClient) AcquireTarget(ctx context.Context, in *DestroyerRequest, opts ...grpc.CallOption) (*DestroyerResponse, error) {
	out := new(DestroyerResponse)
	err := c.cc.Invoke(ctx, "/destroyer.DestroyerService/AcquireTarget", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DestroyerServiceServer is the server API for DestroyerService service.
// All implementations should embed UnimplementedDestroyerServiceServer
// for forward compatibility
type DestroyerServiceServer interface {
	AcquireTarget(context.Context, *DestroyerRequest) (*DestroyerResponse, error)
}

// UnimplementedDestroyerServiceServer should be embedded to have forward compatible implementations.
type UnimplementedDestroyerServiceServer struct {
}

func (UnimplementedDestroyerServiceServer) AcquireTarget(context.Context, *DestroyerRequest) (*DestroyerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AcquireTarget not implemented")
}

// UnsafeDestroyerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DestroyerServiceServer will
// result in compilation errors.
type UnsafeDestroyerServiceServer interface {
	mustEmbedUnimplementedDestroyerServiceServer()
}

func RegisterDestroyerServiceServer(s grpc.ServiceRegistrar, srv DestroyerServiceServer) {
	s.RegisterService(&DestroyerService_ServiceDesc, srv)
}

func _DestroyerService_AcquireTarget_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DestroyerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DestroyerServiceServer).AcquireTarget(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/destroyer.DestroyerService/AcquireTarget",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DestroyerServiceServer).AcquireTarget(ctx, req.(*DestroyerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DestroyerService_ServiceDesc is the grpc.ServiceDesc for DestroyerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DestroyerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "destroyer.DestroyerService",
	HandlerType: (*DestroyerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AcquireTarget",
			Handler:    _DestroyerService_AcquireTarget_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "destroyer.proto",
}