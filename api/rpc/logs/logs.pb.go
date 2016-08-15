// Code generated by protoc-gen-go.
// source: logs.proto
// DO NOT EDIT!

/*
Package logs is a generated protocol buffer package.

It is generated from these files:
	logs.proto

It has these top-level messages:
	LogEntry
	GetRequest
	GetReply
*/
package logs

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

type LogEntry struct {
	Timestamp   string `protobuf:"bytes,1,opt,name=timestamp" json:"timestamp,omitempty"`
	TimeId      string `protobuf:"bytes,2,opt,name=time_id,json=timeId" json:"time_id,omitempty"`
	ServiceId   string `protobuf:"bytes,3,opt,name=service_id,json=serviceId" json:"service_id,omitempty"`
	ServiceName string `protobuf:"bytes,4,opt,name=service_name,json=serviceName" json:"service_name,omitempty"`
	Message     string `protobuf:"bytes,5,opt,name=message" json:"message,omitempty"`
	ContainerId string `protobuf:"bytes,6,opt,name=container_id,json=containerId" json:"container_id,omitempty"`
	NodeId      string `protobuf:"bytes,7,opt,name=node_id,json=nodeId" json:"node_id,omitempty"`
}

func (m *LogEntry) Reset()                    { *m = LogEntry{} }
func (m *LogEntry) String() string            { return proto.CompactTextString(m) }
func (*LogEntry) ProtoMessage()               {}
func (*LogEntry) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type GetRequest struct {
	Timestamp   string `protobuf:"bytes,1,opt,name=timestamp" json:"timestamp,omitempty"`
	ServiceId   string `protobuf:"bytes,2,opt,name=service_id,json=serviceId" json:"service_id,omitempty"`
	ServiceName string `protobuf:"bytes,3,opt,name=service_name,json=serviceName" json:"service_name,omitempty"`
	Message     string `protobuf:"bytes,4,opt,name=message" json:"message,omitempty"`
	ContainerId string `protobuf:"bytes,5,opt,name=container_id,json=containerId" json:"container_id,omitempty"`
	NodeId      string `protobuf:"bytes,6,opt,name=node_id,json=nodeId" json:"node_id,omitempty"`
	From        int64  `protobuf:"zigzag64,7,opt,name=from" json:"from,omitempty"`
	Size        int64  `protobuf:"zigzag64,8,opt,name=size" json:"size,omitempty"`
}

func (m *GetRequest) Reset()                    { *m = GetRequest{} }
func (m *GetRequest) String() string            { return proto.CompactTextString(m) }
func (*GetRequest) ProtoMessage()               {}
func (*GetRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type GetReply struct {
	Entries []*LogEntry `protobuf:"bytes,1,rep,name=entries" json:"entries,omitempty"`
}

func (m *GetReply) Reset()                    { *m = GetReply{} }
func (m *GetReply) String() string            { return proto.CompactTextString(m) }
func (*GetReply) ProtoMessage()               {}
func (*GetReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *GetReply) GetEntries() []*LogEntry {
	if m != nil {
		return m.Entries
	}
	return nil
}

func init() {
	proto.RegisterType((*LogEntry)(nil), "logs.LogEntry")
	proto.RegisterType((*GetRequest)(nil), "logs.GetRequest")
	proto.RegisterType((*GetReply)(nil), "logs.GetReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for Logs service

type LogsClient interface {
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetReply, error)
}

type logsClient struct {
	cc *grpc.ClientConn
}

func NewLogsClient(cc *grpc.ClientConn) LogsClient {
	return &logsClient{cc}
}

func (c *logsClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetReply, error) {
	out := new(GetReply)
	err := grpc.Invoke(ctx, "/logs.Logs/Get", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Logs service

type LogsServer interface {
	Get(context.Context, *GetRequest) (*GetReply, error)
}

func RegisterLogsServer(s *grpc.Server, srv LogsServer) {
	s.RegisterService(&_Logs_serviceDesc, srv)
}

func _Logs_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogsServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/logs.Logs/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogsServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Logs_serviceDesc = grpc.ServiceDesc{
	ServiceName: "logs.Logs",
	HandlerType: (*LogsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _Logs_Get_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fileDescriptor0,
}

func init() { proto.RegisterFile("logs.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 302 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x84, 0x92, 0xc1, 0x4a, 0xc3, 0x40,
	0x10, 0x86, 0x4d, 0x93, 0x26, 0xed, 0x54, 0x8a, 0xcc, 0xc5, 0x45, 0x14, 0x6a, 0x4e, 0xf5, 0x52,
	0xb0, 0xfa, 0x0a, 0x52, 0x02, 0xc5, 0x43, 0x5e, 0x40, 0x62, 0x77, 0x0c, 0x0b, 0xdd, 0xdd, 0xb8,
	0xbb, 0x0a, 0xf5, 0x4d, 0x7d, 0x05, 0x9f, 0x42, 0x76, 0xd3, 0xd8, 0x7a, 0x49, 0x6f, 0xf3, 0xff,
	0xff, 0x4c, 0x32, 0x5f, 0x26, 0x00, 0x5b, 0x5d, 0xdb, 0x45, 0x63, 0xb4, 0xd3, 0x98, 0xf8, 0x3a,
	0xff, 0x8e, 0x60, 0xb4, 0xd6, 0xf5, 0x93, 0x72, 0x66, 0x87, 0xd7, 0x30, 0x76, 0x42, 0x92, 0x75,
	0x95, 0x6c, 0x58, 0x34, 0x8b, 0xe6, 0xe3, 0xf2, 0x60, 0xe0, 0x25, 0x64, 0x5e, 0xbc, 0x08, 0xce,
	0x06, 0x21, 0x4b, 0xbd, 0x2c, 0x38, 0xde, 0x00, 0x58, 0x32, 0x9f, 0x62, 0x13, 0xb2, 0xb8, 0x9d,
	0xdb, 0x3b, 0x05, 0xc7, 0x5b, 0x38, 0xef, 0x62, 0x55, 0x49, 0x62, 0x49, 0x68, 0x98, 0xec, 0xbd,
	0xe7, 0x4a, 0x12, 0x32, 0xc8, 0x24, 0x59, 0x5b, 0xd5, 0xc4, 0x86, 0x21, 0xed, 0xa4, 0x1f, 0xde,
	0x68, 0xe5, 0x2a, 0xa1, 0xc8, 0xf8, 0xa7, 0xa7, 0xed, 0xf0, 0x9f, 0x57, 0x70, 0xbf, 0x97, 0xd2,
	0x3c, 0xbc, 0x3b, 0x6b, 0xf7, 0xf2, 0xb2, 0xe0, 0xf9, 0x4f, 0x04, 0xb0, 0x22, 0x57, 0xd2, 0xfb,
	0x07, 0x59, 0x77, 0x82, 0xee, 0x3f, 0xc4, 0xe0, 0x14, 0x44, 0xdc, 0x0b, 0x91, 0xf4, 0x43, 0x0c,
	0x7b, 0x21, 0xd2, 0x63, 0x08, 0x44, 0x48, 0xde, 0x8c, 0x96, 0x01, 0x0d, 0xcb, 0x50, 0x7b, 0xcf,
	0x8a, 0x2f, 0x62, 0xa3, 0xd6, 0xf3, 0x75, 0xfe, 0x08, 0xa3, 0xc0, 0xda, 0x6c, 0x77, 0x38, 0x87,
	0x8c, 0x94, 0x33, 0x82, 0x2c, 0x8b, 0x66, 0xf1, 0x7c, 0xb2, 0x9c, 0x2e, 0xc2, 0xe1, 0xbb, 0x43,
	0x97, 0x5d, 0xbc, 0xbc, 0x87, 0x64, 0xad, 0x6b, 0x8b, 0x77, 0x10, 0xaf, 0xc8, 0xe1, 0x45, 0xdb,
	0x77, 0xf8, 0x68, 0x57, 0xd3, 0x23, 0xa7, 0xd9, 0xee, 0xf2, 0xb3, 0xd7, 0x34, 0xfc, 0x3e, 0x0f,
	0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x75, 0x5a, 0xe1, 0x76, 0x4c, 0x02, 0x00, 0x00,
}
