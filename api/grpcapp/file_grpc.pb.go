// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: api/grpcapp/file.proto

package grpcapp

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

// FileServiceClient is the client API for FileService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FileServiceClient interface {
	DeleteFileById(ctx context.Context, in *DeleteFileByIdParam, opts ...grpc.CallOption) (*DeleteFileByIdResult, error)
	RetrieveFileById(ctx context.Context, in *RetrieveFileByIdParam, opts ...grpc.CallOption) (FileService_RetrieveFileByIdClient, error)
	UploadFile(ctx context.Context, opts ...grpc.CallOption) (FileService_UploadFileClient, error)
}

type fileServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFileServiceClient(cc grpc.ClientConnInterface) FileServiceClient {
	return &fileServiceClient{cc}
}

func (c *fileServiceClient) DeleteFileById(ctx context.Context, in *DeleteFileByIdParam, opts ...grpc.CallOption) (*DeleteFileByIdResult, error) {
	out := new(DeleteFileByIdResult)
	err := c.cc.Invoke(ctx, "/file.v1.FileService/DeleteFileById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileServiceClient) RetrieveFileById(ctx context.Context, in *RetrieveFileByIdParam, opts ...grpc.CallOption) (FileService_RetrieveFileByIdClient, error) {
	stream, err := c.cc.NewStream(ctx, &FileService_ServiceDesc.Streams[0], "/file.v1.FileService/RetrieveFileById", opts...)
	if err != nil {
		return nil, err
	}
	x := &fileServiceRetrieveFileByIdClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type FileService_RetrieveFileByIdClient interface {
	Recv() (*RetrieveFileByIdResult, error)
	grpc.ClientStream
}

type fileServiceRetrieveFileByIdClient struct {
	grpc.ClientStream
}

func (x *fileServiceRetrieveFileByIdClient) Recv() (*RetrieveFileByIdResult, error) {
	m := new(RetrieveFileByIdResult)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *fileServiceClient) UploadFile(ctx context.Context, opts ...grpc.CallOption) (FileService_UploadFileClient, error) {
	stream, err := c.cc.NewStream(ctx, &FileService_ServiceDesc.Streams[1], "/file.v1.FileService/UploadFile", opts...)
	if err != nil {
		return nil, err
	}
	x := &fileServiceUploadFileClient{stream}
	return x, nil
}

type FileService_UploadFileClient interface {
	Send(*UploadFileParam) error
	CloseAndRecv() (*UploadFileResult, error)
	grpc.ClientStream
}

type fileServiceUploadFileClient struct {
	grpc.ClientStream
}

func (x *fileServiceUploadFileClient) Send(m *UploadFileParam) error {
	return x.ClientStream.SendMsg(m)
}

func (x *fileServiceUploadFileClient) CloseAndRecv() (*UploadFileResult, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(UploadFileResult)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// FileServiceServer is the server API for FileService service.
// All implementations should embed UnimplementedFileServiceServer
// for forward compatibility
type FileServiceServer interface {
	DeleteFileById(context.Context, *DeleteFileByIdParam) (*DeleteFileByIdResult, error)
	RetrieveFileById(*RetrieveFileByIdParam, FileService_RetrieveFileByIdServer) error
	UploadFile(FileService_UploadFileServer) error
}

// UnimplementedFileServiceServer should be embedded to have forward compatible implementations.
type UnimplementedFileServiceServer struct {
}

func (UnimplementedFileServiceServer) DeleteFileById(context.Context, *DeleteFileByIdParam) (*DeleteFileByIdResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteFileById not implemented")
}
func (UnimplementedFileServiceServer) RetrieveFileById(*RetrieveFileByIdParam, FileService_RetrieveFileByIdServer) error {
	return status.Errorf(codes.Unimplemented, "method RetrieveFileById not implemented")
}
func (UnimplementedFileServiceServer) UploadFile(FileService_UploadFileServer) error {
	return status.Errorf(codes.Unimplemented, "method UploadFile not implemented")
}

// UnsafeFileServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FileServiceServer will
// result in compilation errors.
type UnsafeFileServiceServer interface {
	mustEmbedUnimplementedFileServiceServer()
}

func RegisterFileServiceServer(s grpc.ServiceRegistrar, srv FileServiceServer) {
	s.RegisterService(&FileService_ServiceDesc, srv)
}

func _FileService_DeleteFileById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteFileByIdParam)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServiceServer).DeleteFileById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/file.v1.FileService/DeleteFileById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServiceServer).DeleteFileById(ctx, req.(*DeleteFileByIdParam))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileService_RetrieveFileById_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(RetrieveFileByIdParam)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(FileServiceServer).RetrieveFileById(m, &fileServiceRetrieveFileByIdServer{stream})
}

type FileService_RetrieveFileByIdServer interface {
	Send(*RetrieveFileByIdResult) error
	grpc.ServerStream
}

type fileServiceRetrieveFileByIdServer struct {
	grpc.ServerStream
}

func (x *fileServiceRetrieveFileByIdServer) Send(m *RetrieveFileByIdResult) error {
	return x.ServerStream.SendMsg(m)
}

func _FileService_UploadFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FileServiceServer).UploadFile(&fileServiceUploadFileServer{stream})
}

type FileService_UploadFileServer interface {
	SendAndClose(*UploadFileResult) error
	Recv() (*UploadFileParam, error)
	grpc.ServerStream
}

type fileServiceUploadFileServer struct {
	grpc.ServerStream
}

func (x *fileServiceUploadFileServer) SendAndClose(m *UploadFileResult) error {
	return x.ServerStream.SendMsg(m)
}

func (x *fileServiceUploadFileServer) Recv() (*UploadFileParam, error) {
	m := new(UploadFileParam)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// FileService_ServiceDesc is the grpc.ServiceDesc for FileService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FileService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "file.v1.FileService",
	HandlerType: (*FileServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DeleteFileById",
			Handler:    _FileService_DeleteFileById_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "RetrieveFileById",
			Handler:       _FileService_RetrieveFileById_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "UploadFile",
			Handler:       _FileService_UploadFile_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "api/grpcapp/file.proto",
}
