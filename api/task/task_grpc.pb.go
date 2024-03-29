// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.22.2
// source: api/proto/task.proto

package task

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
	TaskService_SubscribeTask_FullMethodName = "/task.TaskService/SubscribeTask"
)

// TaskServiceClient is the client API for TaskService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TaskServiceClient interface {
	// SubscribeTask register a task client to the task server and wait to receive task
	SubscribeTask(ctx context.Context, in *ClientInfo, opts ...grpc.CallOption) (TaskService_SubscribeTaskClient, error)
}

type taskServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTaskServiceClient(cc grpc.ClientConnInterface) TaskServiceClient {
	return &taskServiceClient{cc}
}

func (c *taskServiceClient) SubscribeTask(ctx context.Context, in *ClientInfo, opts ...grpc.CallOption) (TaskService_SubscribeTaskClient, error) {
	stream, err := c.cc.NewStream(ctx, &TaskService_ServiceDesc.Streams[0], TaskService_SubscribeTask_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &taskServiceSubscribeTaskClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type TaskService_SubscribeTaskClient interface {
	Recv() (*TaskInfo, error)
	grpc.ClientStream
}

type taskServiceSubscribeTaskClient struct {
	grpc.ClientStream
}

func (x *taskServiceSubscribeTaskClient) Recv() (*TaskInfo, error) {
	m := new(TaskInfo)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TaskServiceServer is the server API for TaskService service.
// All implementations must embed UnimplementedTaskServiceServer
// for forward compatibility
type TaskServiceServer interface {
	// SubscribeTask register a task client to the task server and wait to receive task
	SubscribeTask(*ClientInfo, TaskService_SubscribeTaskServer) error
	mustEmbedUnimplementedTaskServiceServer()
}

// UnimplementedTaskServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTaskServiceServer struct {
}

func (UnimplementedTaskServiceServer) SubscribeTask(*ClientInfo, TaskService_SubscribeTaskServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeTask not implemented")
}
func (UnimplementedTaskServiceServer) mustEmbedUnimplementedTaskServiceServer() {}

// UnsafeTaskServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TaskServiceServer will
// result in compilation errors.
type UnsafeTaskServiceServer interface {
	mustEmbedUnimplementedTaskServiceServer()
}

func RegisterTaskServiceServer(s grpc.ServiceRegistrar, srv TaskServiceServer) {
	s.RegisterService(&TaskService_ServiceDesc, srv)
}

func _TaskService_SubscribeTask_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ClientInfo)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TaskServiceServer).SubscribeTask(m, &taskServiceSubscribeTaskServer{stream})
}

type TaskService_SubscribeTaskServer interface {
	Send(*TaskInfo) error
	grpc.ServerStream
}

type taskServiceSubscribeTaskServer struct {
	grpc.ServerStream
}

func (x *taskServiceSubscribeTaskServer) Send(m *TaskInfo) error {
	return x.ServerStream.SendMsg(m)
}

// TaskService_ServiceDesc is the grpc.ServiceDesc for TaskService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TaskService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "task.TaskService",
	HandlerType: (*TaskServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SubscribeTask",
			Handler:       _TaskService_SubscribeTask_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/proto/task.proto",
}
