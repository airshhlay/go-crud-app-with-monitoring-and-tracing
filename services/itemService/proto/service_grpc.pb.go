// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.4
// source: proto/service.proto

package proto

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

// ItemServiceClient is the client API for ItemService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ItemServiceClient interface {
	DeleteFav(ctx context.Context, in *DeleteFavReq, opts ...grpc.CallOption) (*DeleteFavRes, error)
	AddFav(ctx context.Context, in *AddFavReq, opts ...grpc.CallOption) (*AddFavRes, error)
	GetFavList(ctx context.Context, in *GetFavListReq, opts ...grpc.CallOption) (*GetFavListRes, error)
}

type itemServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewItemServiceClient(cc grpc.ClientConnInterface) ItemServiceClient {
	return &itemServiceClient{cc}
}

func (c *itemServiceClient) DeleteFav(ctx context.Context, in *DeleteFavReq, opts ...grpc.CallOption) (*DeleteFavRes, error) {
	out := new(DeleteFavRes)
	err := c.cc.Invoke(ctx, "/proto.ItemService/DeleteFav", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *itemServiceClient) AddFav(ctx context.Context, in *AddFavReq, opts ...grpc.CallOption) (*AddFavRes, error) {
	out := new(AddFavRes)
	err := c.cc.Invoke(ctx, "/proto.ItemService/AddFav", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *itemServiceClient) GetFavList(ctx context.Context, in *GetFavListReq, opts ...grpc.CallOption) (*GetFavListRes, error) {
	out := new(GetFavListRes)
	err := c.cc.Invoke(ctx, "/proto.ItemService/GetFavList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ItemServiceServer is the server API for ItemService service.
// All implementations must embed UnimplementedItemServiceServer
// for forward compatibility
type ItemServiceServer interface {
	DeleteFav(context.Context, *DeleteFavReq) (*DeleteFavRes, error)
	AddFav(context.Context, *AddFavReq) (*AddFavRes, error)
	GetFavList(context.Context, *GetFavListReq) (*GetFavListRes, error)
	mustEmbedUnimplementedItemServiceServer()
}

// UnimplementedItemServiceServer must be embedded to have forward compatible implementations.
type UnimplementedItemServiceServer struct {
}

func (UnimplementedItemServiceServer) DeleteFav(context.Context, *DeleteFavReq) (*DeleteFavRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteFav not implemented")
}
func (UnimplementedItemServiceServer) AddFav(context.Context, *AddFavReq) (*AddFavRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddFav not implemented")
}
func (UnimplementedItemServiceServer) GetFavList(context.Context, *GetFavListReq) (*GetFavListRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFavList not implemented")
}
func (UnimplementedItemServiceServer) mustEmbedUnimplementedItemServiceServer() {}

// UnsafeItemServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ItemServiceServer will
// result in compilation errors.
type UnsafeItemServiceServer interface {
	mustEmbedUnimplementedItemServiceServer()
}

func RegisterItemServiceServer(s grpc.ServiceRegistrar, srv ItemServiceServer) {
	s.RegisterService(&ItemService_ServiceDesc, srv)
}

func _ItemService_DeleteFav_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteFavReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemServiceServer).DeleteFav(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ItemService/DeleteFav",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemServiceServer).DeleteFav(ctx, req.(*DeleteFavReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _ItemService_AddFav_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddFavReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemServiceServer).AddFav(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ItemService/AddFav",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemServiceServer).AddFav(ctx, req.(*AddFavReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _ItemService_GetFavList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFavListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemServiceServer).GetFavList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ItemService/GetFavList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemServiceServer).GetFavList(ctx, req.(*GetFavListReq))
	}
	return interceptor(ctx, in, info, handler)
}

// ItemService_ServiceDesc is the grpc.ServiceDesc for ItemService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ItemService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.ItemService",
	HandlerType: (*ItemServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DeleteFav",
			Handler:    _ItemService_DeleteFav_Handler,
		},
		{
			MethodName: "AddFav",
			Handler:    _ItemService_AddFav_Handler,
		},
		{
			MethodName: "GetFavList",
			Handler:    _ItemService_GetFavList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/service.proto",
}
