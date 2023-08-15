// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: bench.proto

package bench

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

// BenchClient is the client API for Bench service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BenchClient interface {
	// StartPlan 启动压测
	StartPlan(ctx context.Context, in *StartPlanRequest, opts ...grpc.CallOption) (*StartPlanResponse, error)
	// SendChat 发送消息
	SendChat(ctx context.Context, in *SendChatRequest, opts ...grpc.CallOption) (*SendChatResponse, error)
	// CustomSend 自定义消息发送
	CustomSend(ctx context.Context, in *CustomSendRequest, opts ...grpc.CallOption) (*CustomSendResponse, error)
	// ClosePlan 清理任务
	ClosePlan(ctx context.Context, in *ClosePlanRequest, opts ...grpc.CallOption) (*ClosePlanResponse, error)
	// ListPlan 查询压测计划
	ListPlan(ctx context.Context, in *ListPlanRequest, opts ...grpc.CallOption) (*ListPlanResponse, error)
	// DetailPlan 查询压测计划详情
	DetailPlan(ctx context.Context, in *DetailPlanRequest, opts ...grpc.CallOption) (*DetailPlanResponse, error)
}

type benchClient struct {
	cc grpc.ClientConnInterface
}

func NewBenchClient(cc grpc.ClientConnInterface) BenchClient {
	return &benchClient{cc}
}

func (c *benchClient) StartPlan(ctx context.Context, in *StartPlanRequest, opts ...grpc.CallOption) (*StartPlanResponse, error) {
	out := new(StartPlanResponse)
	err := c.cc.Invoke(ctx, "/bench.Bench/StartPlan", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *benchClient) SendChat(ctx context.Context, in *SendChatRequest, opts ...grpc.CallOption) (*SendChatResponse, error) {
	out := new(SendChatResponse)
	err := c.cc.Invoke(ctx, "/bench.Bench/SendChat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *benchClient) CustomSend(ctx context.Context, in *CustomSendRequest, opts ...grpc.CallOption) (*CustomSendResponse, error) {
	out := new(CustomSendResponse)
	err := c.cc.Invoke(ctx, "/bench.Bench/CustomSend", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *benchClient) ClosePlan(ctx context.Context, in *ClosePlanRequest, opts ...grpc.CallOption) (*ClosePlanResponse, error) {
	out := new(ClosePlanResponse)
	err := c.cc.Invoke(ctx, "/bench.Bench/ClosePlan", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *benchClient) ListPlan(ctx context.Context, in *ListPlanRequest, opts ...grpc.CallOption) (*ListPlanResponse, error) {
	out := new(ListPlanResponse)
	err := c.cc.Invoke(ctx, "/bench.Bench/ListPlan", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *benchClient) DetailPlan(ctx context.Context, in *DetailPlanRequest, opts ...grpc.CallOption) (*DetailPlanResponse, error) {
	out := new(DetailPlanResponse)
	err := c.cc.Invoke(ctx, "/bench.Bench/DetailPlan", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BenchServer is the server API for Bench service.
// All implementations must embed UnimplementedBenchServer
// for forward compatibility
type BenchServer interface {
	// StartPlan 启动压测
	StartPlan(context.Context, *StartPlanRequest) (*StartPlanResponse, error)
	// SendChat 发送消息
	SendChat(context.Context, *SendChatRequest) (*SendChatResponse, error)
	// CustomSend 自定义消息发送
	CustomSend(context.Context, *CustomSendRequest) (*CustomSendResponse, error)
	// ClosePlan 清理任务
	ClosePlan(context.Context, *ClosePlanRequest) (*ClosePlanResponse, error)
	// ListPlan 查询压测计划
	ListPlan(context.Context, *ListPlanRequest) (*ListPlanResponse, error)
	// DetailPlan 查询压测计划详情
	DetailPlan(context.Context, *DetailPlanRequest) (*DetailPlanResponse, error)
	mustEmbedUnimplementedBenchServer()
}

// UnimplementedBenchServer must be embedded to have forward compatible implementations.
type UnimplementedBenchServer struct {
}

func (UnimplementedBenchServer) StartPlan(context.Context, *StartPlanRequest) (*StartPlanResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartPlan not implemented")
}
func (UnimplementedBenchServer) SendChat(context.Context, *SendChatRequest) (*SendChatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendChat not implemented")
}
func (UnimplementedBenchServer) CustomSend(context.Context, *CustomSendRequest) (*CustomSendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CustomSend not implemented")
}
func (UnimplementedBenchServer) ClosePlan(context.Context, *ClosePlanRequest) (*ClosePlanResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClosePlan not implemented")
}
func (UnimplementedBenchServer) ListPlan(context.Context, *ListPlanRequest) (*ListPlanResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPlan not implemented")
}
func (UnimplementedBenchServer) DetailPlan(context.Context, *DetailPlanRequest) (*DetailPlanResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DetailPlan not implemented")
}
func (UnimplementedBenchServer) mustEmbedUnimplementedBenchServer() {}

// UnsafeBenchServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BenchServer will
// result in compilation errors.
type UnsafeBenchServer interface {
	mustEmbedUnimplementedBenchServer()
}

func RegisterBenchServer(s grpc.ServiceRegistrar, srv BenchServer) {
	s.RegisterService(&Bench_ServiceDesc, srv)
}

func _Bench_StartPlan_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartPlanRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BenchServer).StartPlan(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bench.Bench/StartPlan",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BenchServer).StartPlan(ctx, req.(*StartPlanRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bench_SendChat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendChatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BenchServer).SendChat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bench.Bench/SendChat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BenchServer).SendChat(ctx, req.(*SendChatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bench_CustomSend_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CustomSendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BenchServer).CustomSend(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bench.Bench/CustomSend",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BenchServer).CustomSend(ctx, req.(*CustomSendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bench_ClosePlan_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClosePlanRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BenchServer).ClosePlan(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bench.Bench/ClosePlan",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BenchServer).ClosePlan(ctx, req.(*ClosePlanRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bench_ListPlan_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListPlanRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BenchServer).ListPlan(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bench.Bench/ListPlan",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BenchServer).ListPlan(ctx, req.(*ListPlanRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bench_DetailPlan_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DetailPlanRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BenchServer).DetailPlan(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bench.Bench/DetailPlan",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BenchServer).DetailPlan(ctx, req.(*DetailPlanRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Bench_ServiceDesc is the grpc.ServiceDesc for Bench service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Bench_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "bench.Bench",
	HandlerType: (*BenchServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StartPlan",
			Handler:    _Bench_StartPlan_Handler,
		},
		{
			MethodName: "SendChat",
			Handler:    _Bench_SendChat_Handler,
		},
		{
			MethodName: "CustomSend",
			Handler:    _Bench_CustomSend_Handler,
		},
		{
			MethodName: "ClosePlan",
			Handler:    _Bench_ClosePlan_Handler,
		},
		{
			MethodName: "ListPlan",
			Handler:    _Bench_ListPlan_Handler,
		},
		{
			MethodName: "DetailPlan",
			Handler:    _Bench_DetailPlan_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "bench.proto",
}