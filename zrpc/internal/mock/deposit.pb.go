// Code generated by protoc-gen-go.
// source: deposit.proto
// DO NOT EDIT!

/*
Package mock is a generated protocol buffer package.

It is generated from these files:
	deposit.proto

It has these top-level messages:
	DepositRequest
	DepositResponse
*/
package mock

import (
	proto "github.com/golang/protobuf/proto"
	fmt "fmt"
	math "math"
)

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = proto.Marshal
	_ = fmt.Errorf
	_ = math.Inf
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type DepositRequest struct {
	Amount float32 `protobuf:"fixed32,1,opt,name=amount" json:"amount,omitempty"`
}

func (m *DepositRequest) Reset()                    { *m = DepositRequest{} }
func (m *DepositRequest) String() string            { return proto.CompactTextString(m) }
func (*DepositRequest) ProtoMessage()               {}
func (*DepositRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *DepositRequest) GetAmount() float32 {
	if m != nil {
		return m.Amount
	}
	return 0
}

type DepositResponse struct {
	Ok bool `protobuf:"varint,1,opt,name=ok" json:"ok,omitempty"`
}

func (m *DepositResponse) Reset()                    { *m = DepositResponse{} }
func (m *DepositResponse) String() string            { return proto.CompactTextString(m) }
func (*DepositResponse) ProtoMessage()               {}
func (*DepositResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *DepositResponse) GetOk() bool {
	if m != nil {
		return m.Ok
	}
	return false
}

func init() {
	proto.RegisterType((*DepositRequest)(nil), "mock.DepositRequest")
	proto.RegisterType((*DepositResponse)(nil), "mock.DepositResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ context.Context
	_ grpc.ClientConn
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for DepositService service

type DepositServiceClient interface {
	Deposit(ctx context.Context, in *DepositRequest, opts ...grpc.CallOption) (*DepositResponse, error)
}

type depositServiceClient struct {
	cc *grpc.ClientConn
}

func NewDepositServiceClient(cc *grpc.ClientConn) DepositServiceClient {
	return &depositServiceClient{cc}
}

func (c *depositServiceClient) Deposit(ctx context.Context, in *DepositRequest, opts ...grpc.CallOption) (*DepositResponse, error) {
	out := new(DepositResponse)
	err := grpc.Invoke(ctx, "/mock.DepositService/Deposit", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for DepositService service

type DepositServiceServer interface {
	Deposit(context.Context, *DepositRequest) (*DepositResponse, error)
}

func RegisterDepositServiceServer(s *grpc.Server, srv DepositServiceServer) {
	s.RegisterService(&_DepositService_serviceDesc, srv)
}

func _DepositService_Deposit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DepositRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DepositServiceServer).Deposit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mock.DepositService/Deposit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DepositServiceServer).Deposit(ctx, req.(*DepositRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _DepositService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "mock.DepositService",
	HandlerType: (*DepositServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Deposit",
			Handler:    _DepositService_Deposit_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "deposit.proto",
}

func init() { proto.RegisterFile("deposit.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 139 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x4d, 0x49, 0x2d, 0xc8,
	0x2f, 0xce, 0x2c, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xc9, 0xcd, 0x4f, 0xce, 0x56,
	0xd2, 0xe0, 0xe2, 0x73, 0x81, 0x08, 0x07, 0xa5, 0x16, 0x96, 0xa6, 0x16, 0x97, 0x08, 0x89, 0x71,
	0xb1, 0x25, 0xe6, 0xe6, 0x97, 0xe6, 0x95, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0x30, 0x05, 0x41, 0x79,
	0x4a, 0x8a, 0x5c, 0xfc, 0x70, 0x95, 0xc5, 0x05, 0xf9, 0x79, 0xc5, 0xa9, 0x42, 0x7c, 0x5c, 0x4c,
	0xf9, 0xd9, 0x60, 0x65, 0x1c, 0x41, 0x4c, 0xf9, 0xd9, 0x46, 0x1e, 0x70, 0xc3, 0x82, 0x53, 0x8b,
	0xca, 0x32, 0x93, 0x53, 0x85, 0xcc, 0xb8, 0xd8, 0xa1, 0x22, 0x42, 0x22, 0x7a, 0x20, 0x0b, 0xf5,
	0x50, 0x6d, 0x93, 0x12, 0x45, 0x13, 0x85, 0x98, 0x9c, 0xc4, 0x06, 0x76, 0xa3, 0x31, 0x20, 0x00,
	0x00, 0xff, 0xff, 0x62, 0x37, 0xf2, 0x36, 0xb4, 0x00, 0x00, 0x00,
}
