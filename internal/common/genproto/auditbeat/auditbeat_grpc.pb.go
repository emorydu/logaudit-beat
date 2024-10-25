// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.20.3
// source: auditbeat.proto

package auditbeat

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	AuditBeatService_FetchBeatRule_FullMethodName = "/auditbeat.AuditBeatService/FetchBeatRule"
	AuditBeatService_Download_FullMethodName      = "/auditbeat.AuditBeatService/Download"
	AuditBeatService_UsageStatus_FullMethodName   = "/auditbeat.AuditBeatService/UsageStatus"
	AuditBeatService_Updated_FullMethodName       = "/auditbeat.AuditBeatService/Updated"
)

// AuditBeatServiceClient is the client API for AuditBeatService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuditBeatServiceClient interface {
	FetchBeatRule(ctx context.Context, in *FetchBeatRuleRequest, opts ...grpc.CallOption) (*FetchBeatRuleResponse, error)
	Download(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
	UsageStatus(ctx context.Context, in *UsageStatusRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Updated(ctx context.Context, in *UpdatedRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type auditBeatServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAuditBeatServiceClient(cc grpc.ClientConnInterface) AuditBeatServiceClient {
	return &auditBeatServiceClient{cc}
}

func (c *auditBeatServiceClient) FetchBeatRule(ctx context.Context, in *FetchBeatRuleRequest, opts ...grpc.CallOption) (*FetchBeatRuleResponse, error) {
	out := new(FetchBeatRuleResponse)
	err := c.cc.Invoke(ctx, AuditBeatService_FetchBeatRule_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *auditBeatServiceClient) Download(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, AuditBeatService_Download_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *auditBeatServiceClient) UsageStatus(ctx context.Context, in *UsageStatusRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, AuditBeatService_UsageStatus_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *auditBeatServiceClient) Updated(ctx context.Context, in *UpdatedRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, AuditBeatService_Updated_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuditBeatServiceServer is the server API for AuditBeatService service.
// All implementations should embed UnimplementedAuditBeatServiceServer
// for forward compatibility
type AuditBeatServiceServer interface {
	FetchBeatRule(context.Context, *FetchBeatRuleRequest) (*FetchBeatRuleResponse, error)
	Download(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	UsageStatus(context.Context, *UsageStatusRequest) (*emptypb.Empty, error)
	Updated(context.Context, *UpdatedRequest) (*emptypb.Empty, error)
}

// UnimplementedAuditBeatServiceServer should be embedded to have forward compatible implementations.
type UnimplementedAuditBeatServiceServer struct {
}

func (UnimplementedAuditBeatServiceServer) FetchBeatRule(context.Context, *FetchBeatRuleRequest) (*FetchBeatRuleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FetchBeatRule not implemented")
}
func (UnimplementedAuditBeatServiceServer) Download(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Download not implemented")
}
func (UnimplementedAuditBeatServiceServer) UsageStatus(context.Context, *UsageStatusRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UsageStatus not implemented")
}
func (UnimplementedAuditBeatServiceServer) Updated(context.Context, *UpdatedRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Updated not implemented")
}

// UnsafeAuditBeatServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuditBeatServiceServer will
// result in compilation errors.
type UnsafeAuditBeatServiceServer interface {
	mustEmbedUnimplementedAuditBeatServiceServer()
}

func RegisterAuditBeatServiceServer(s grpc.ServiceRegistrar, srv AuditBeatServiceServer) {
	s.RegisterService(&AuditBeatService_ServiceDesc, srv)
}

func _AuditBeatService_FetchBeatRule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FetchBeatRuleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuditBeatServiceServer).FetchBeatRule(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuditBeatService_FetchBeatRule_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuditBeatServiceServer).FetchBeatRule(ctx, req.(*FetchBeatRuleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuditBeatService_Download_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuditBeatServiceServer).Download(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuditBeatService_Download_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuditBeatServiceServer).Download(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuditBeatService_UsageStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UsageStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuditBeatServiceServer).UsageStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuditBeatService_UsageStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuditBeatServiceServer).UsageStatus(ctx, req.(*UsageStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuditBeatService_Updated_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuditBeatServiceServer).Updated(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuditBeatService_Updated_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuditBeatServiceServer).Updated(ctx, req.(*UpdatedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AuditBeatService_ServiceDesc is the grpc.ServiceDesc for AuditBeatService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuditBeatService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "auditbeat.AuditBeatService",
	HandlerType: (*AuditBeatServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FetchBeatRule",
			Handler:    _AuditBeatService_FetchBeatRule_Handler,
		},
		{
			MethodName: "Download",
			Handler:    _AuditBeatService_Download_Handler,
		},
		{
			MethodName: "UsageStatus",
			Handler:    _AuditBeatService_UsageStatus_Handler,
		},
		{
			MethodName: "Updated",
			Handler:    _AuditBeatService_Updated_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "auditbeat.proto",
}
