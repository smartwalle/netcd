// Code generated by protoc-gen-go. DO NOT EDIT.
// source: hw.proto

package hw

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type FirstRequest struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FirstRequest) Reset()         { *m = FirstRequest{} }
func (m *FirstRequest) String() string { return proto.CompactTextString(m) }
func (*FirstRequest) ProtoMessage()    {}
func (*FirstRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_hw_1ab919a2e4cdac16, []int{0}
}
func (m *FirstRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FirstRequest.Unmarshal(m, b)
}
func (m *FirstRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FirstRequest.Marshal(b, m, deterministic)
}
func (dst *FirstRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FirstRequest.Merge(dst, src)
}
func (m *FirstRequest) XXX_Size() int {
	return xxx_messageInfo_FirstRequest.Size(m)
}
func (m *FirstRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_FirstRequest.DiscardUnknown(m)
}

var xxx_messageInfo_FirstRequest proto.InternalMessageInfo

func (m *FirstRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type FirstResponse struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FirstResponse) Reset()         { *m = FirstResponse{} }
func (m *FirstResponse) String() string { return proto.CompactTextString(m) }
func (*FirstResponse) ProtoMessage()    {}
func (*FirstResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_hw_1ab919a2e4cdac16, []int{1}
}
func (m *FirstResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FirstResponse.Unmarshal(m, b)
}
func (m *FirstResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FirstResponse.Marshal(b, m, deterministic)
}
func (dst *FirstResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FirstResponse.Merge(dst, src)
}
func (m *FirstResponse) XXX_Size() int {
	return xxx_messageInfo_FirstResponse.Size(m)
}
func (m *FirstResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_FirstResponse.DiscardUnknown(m)
}

var xxx_messageInfo_FirstResponse proto.InternalMessageInfo

func (m *FirstResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*FirstRequest)(nil), "hw.FirstRequest")
	proto.RegisterType((*FirstResponse)(nil), "hw.FirstResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// FirstGRPCClient is the client API for FirstGRPC service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type FirstGRPCClient interface {
	FirstCall(ctx context.Context, in *FirstRequest, opts ...grpc.CallOption) (*FirstResponse, error)
}

type firstGRPCClient struct {
	cc *grpc.ClientConn
}

func NewFirstGRPCClient(cc *grpc.ClientConn) FirstGRPCClient {
	return &firstGRPCClient{cc}
}

func (c *firstGRPCClient) FirstCall(ctx context.Context, in *FirstRequest, opts ...grpc.CallOption) (*FirstResponse, error) {
	out := new(FirstResponse)
	err := c.cc.Invoke(ctx, "/hw.FirstGRPC/FirstCall", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FirstGRPCServer is the server API for FirstGRPC service.
type FirstGRPCServer interface {
	FirstCall(context.Context, *FirstRequest) (*FirstResponse, error)
}

func RegisterFirstGRPCServer(s *grpc.Server, srv FirstGRPCServer) {
	s.RegisterService(&_FirstGRPC_serviceDesc, srv)
}

func _FirstGRPC_FirstCall_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FirstRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FirstGRPCServer).FirstCall(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hw.FirstGRPC/FirstCall",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FirstGRPCServer).FirstCall(ctx, req.(*FirstRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _FirstGRPC_serviceDesc = grpc.ServiceDesc{
	ServiceName: "hw.FirstGRPC",
	HandlerType: (*FirstGRPCServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FirstCall",
			Handler:    _FirstGRPC_FirstCall_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "hw.proto",
}

func init() { proto.RegisterFile("hw.proto", fileDescriptor_hw_1ab919a2e4cdac16) }

var fileDescriptor_hw_1ab919a2e4cdac16 = []byte{
	// 136 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xc8, 0x28, 0xd7, 0x2b,
	0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xca, 0x28, 0x57, 0x52, 0xe2, 0xe2, 0x71, 0xcb, 0x2c, 0x2a,
	0x2e, 0x09, 0x4a, 0x2d, 0x2c, 0x4d, 0x2d, 0x2e, 0x11, 0x12, 0xe2, 0x62, 0xc9, 0x4b, 0xcc, 0x4d,
	0x95, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x02, 0xb3, 0x95, 0x34, 0xb9, 0x78, 0xa1, 0x6a, 0x8a,
	0x0b, 0xf2, 0xf3, 0x8a, 0x53, 0x85, 0x24, 0xb8, 0xd8, 0x73, 0x53, 0x8b, 0x8b, 0x13, 0xd3, 0x61,
	0xea, 0x60, 0x5c, 0x23, 0x7b, 0x2e, 0x4e, 0xb0, 0x52, 0xf7, 0xa0, 0x00, 0x67, 0x21, 0x23, 0x28,
	0xc7, 0x39, 0x31, 0x27, 0x47, 0x48, 0x40, 0x2f, 0xa3, 0x5c, 0x0f, 0xd9, 0x2a, 0x29, 0x41, 0x24,
	0x11, 0x88, 0xc1, 0x4a, 0x0c, 0x49, 0x6c, 0x60, 0xa7, 0x19, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff,
	0x64, 0xfb, 0x61, 0x9d, 0xa6, 0x00, 0x00, 0x00,
}
