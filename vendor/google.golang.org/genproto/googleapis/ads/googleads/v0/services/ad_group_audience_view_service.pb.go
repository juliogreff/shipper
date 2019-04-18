// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/ads/googleads/v0/services/ad_group_audience_view_service.proto

package services // import "google.golang.org/genproto/googleapis/ads/googleads/v0/services"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import resources "google.golang.org/genproto/googleapis/ads/googleads/v0/resources"
import _ "google.golang.org/genproto/googleapis/api/annotations"

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

// Request message for [AdGroupAudienceViewService.GetAdGoupAudienceView][].
type GetAdGroupAudienceViewRequest struct {
	// The resource name of the ad group audience view to fetch.
	ResourceName         string   `protobuf:"bytes,1,opt,name=resource_name,json=resourceName,proto3" json:"resource_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetAdGroupAudienceViewRequest) Reset()         { *m = GetAdGroupAudienceViewRequest{} }
func (m *GetAdGroupAudienceViewRequest) String() string { return proto.CompactTextString(m) }
func (*GetAdGroupAudienceViewRequest) ProtoMessage()    {}
func (*GetAdGroupAudienceViewRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ad_group_audience_view_service_63c90b1e31981bbe, []int{0}
}
func (m *GetAdGroupAudienceViewRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAdGroupAudienceViewRequest.Unmarshal(m, b)
}
func (m *GetAdGroupAudienceViewRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAdGroupAudienceViewRequest.Marshal(b, m, deterministic)
}
func (dst *GetAdGroupAudienceViewRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAdGroupAudienceViewRequest.Merge(dst, src)
}
func (m *GetAdGroupAudienceViewRequest) XXX_Size() int {
	return xxx_messageInfo_GetAdGroupAudienceViewRequest.Size(m)
}
func (m *GetAdGroupAudienceViewRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAdGroupAudienceViewRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetAdGroupAudienceViewRequest proto.InternalMessageInfo

func (m *GetAdGroupAudienceViewRequest) GetResourceName() string {
	if m != nil {
		return m.ResourceName
	}
	return ""
}

func init() {
	proto.RegisterType((*GetAdGroupAudienceViewRequest)(nil), "google.ads.googleads.v0.services.GetAdGroupAudienceViewRequest")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// AdGroupAudienceViewServiceClient is the client API for AdGroupAudienceViewService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AdGroupAudienceViewServiceClient interface {
	// Returns the requested ad group audience view in full detail.
	GetAdGroupAudienceView(ctx context.Context, in *GetAdGroupAudienceViewRequest, opts ...grpc.CallOption) (*resources.AdGroupAudienceView, error)
}

type adGroupAudienceViewServiceClient struct {
	cc *grpc.ClientConn
}

func NewAdGroupAudienceViewServiceClient(cc *grpc.ClientConn) AdGroupAudienceViewServiceClient {
	return &adGroupAudienceViewServiceClient{cc}
}

func (c *adGroupAudienceViewServiceClient) GetAdGroupAudienceView(ctx context.Context, in *GetAdGroupAudienceViewRequest, opts ...grpc.CallOption) (*resources.AdGroupAudienceView, error) {
	out := new(resources.AdGroupAudienceView)
	err := c.cc.Invoke(ctx, "/google.ads.googleads.v0.services.AdGroupAudienceViewService/GetAdGroupAudienceView", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AdGroupAudienceViewServiceServer is the server API for AdGroupAudienceViewService service.
type AdGroupAudienceViewServiceServer interface {
	// Returns the requested ad group audience view in full detail.
	GetAdGroupAudienceView(context.Context, *GetAdGroupAudienceViewRequest) (*resources.AdGroupAudienceView, error)
}

func RegisterAdGroupAudienceViewServiceServer(s *grpc.Server, srv AdGroupAudienceViewServiceServer) {
	s.RegisterService(&_AdGroupAudienceViewService_serviceDesc, srv)
}

func _AdGroupAudienceViewService_GetAdGroupAudienceView_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAdGroupAudienceViewRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdGroupAudienceViewServiceServer).GetAdGroupAudienceView(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.ads.googleads.v0.services.AdGroupAudienceViewService/GetAdGroupAudienceView",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdGroupAudienceViewServiceServer).GetAdGroupAudienceView(ctx, req.(*GetAdGroupAudienceViewRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _AdGroupAudienceViewService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "google.ads.googleads.v0.services.AdGroupAudienceViewService",
	HandlerType: (*AdGroupAudienceViewServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAdGroupAudienceView",
			Handler:    _AdGroupAudienceViewService_GetAdGroupAudienceView_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "google/ads/googleads/v0/services/ad_group_audience_view_service.proto",
}

func init() {
	proto.RegisterFile("google/ads/googleads/v0/services/ad_group_audience_view_service.proto", fileDescriptor_ad_group_audience_view_service_63c90b1e31981bbe)
}

var fileDescriptor_ad_group_audience_view_service_63c90b1e31981bbe = []byte{
	// 377 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x52, 0xbd, 0x4a, 0xec, 0x40,
	0x14, 0x26, 0xb9, 0x70, 0xe1, 0x86, 0x7b, 0x9b, 0x14, 0x17, 0x09, 0x8a, 0xcb, 0xba, 0x85, 0x6c,
	0x31, 0x13, 0x14, 0x16, 0x19, 0x71, 0x25, 0x8b, 0x12, 0x2b, 0x59, 0x56, 0x48, 0x21, 0x81, 0x30,
	0x66, 0x86, 0x10, 0xd8, 0x64, 0x62, 0x4e, 0x92, 0x2d, 0xc4, 0xc6, 0xc2, 0x17, 0xf0, 0x0d, 0x2c,
	0x7d, 0x14, 0x3b, 0xf1, 0x15, 0xac, 0xac, 0x7c, 0x04, 0xc9, 0x4e, 0x26, 0x20, 0x6e, 0xdc, 0xee,
	0x63, 0xe6, 0xfb, 0x39, 0xf3, 0x9d, 0x31, 0x4e, 0x23, 0x21, 0xa2, 0x39, 0xc7, 0x94, 0x01, 0x96,
	0xb0, 0x46, 0x95, 0x8d, 0x81, 0xe7, 0x55, 0x1c, 0x72, 0xc0, 0x94, 0x05, 0x51, 0x2e, 0xca, 0x2c,
	0xa0, 0x25, 0x8b, 0x79, 0x1a, 0xf2, 0xa0, 0x8a, 0xf9, 0x22, 0x68, 0xee, 0x51, 0x96, 0x8b, 0x42,
	0x98, 0x3d, 0xa9, 0x45, 0x94, 0x01, 0x6a, 0x6d, 0x50, 0x65, 0x23, 0x65, 0x63, 0x8d, 0xbb, 0x82,
	0x72, 0x0e, 0xa2, 0xcc, 0xbb, 0x93, 0x64, 0x82, 0xb5, 0xa9, 0xf4, 0x59, 0x8c, 0x69, 0x9a, 0x8a,
	0x82, 0x16, 0xb1, 0x48, 0x41, 0xde, 0xf6, 0x4f, 0x8c, 0x2d, 0x97, 0x17, 0x0e, 0x73, 0x6b, 0xbd,
	0xd3, 0xc8, 0xbd, 0x98, 0x2f, 0x66, 0xfc, 0xba, 0xe4, 0x50, 0x98, 0x3b, 0xc6, 0x3f, 0x15, 0x14,
	0xa4, 0x34, 0xe1, 0x1b, 0x5a, 0x4f, 0xdb, 0xfd, 0x33, 0xfb, 0xab, 0x0e, 0xcf, 0x69, 0xc2, 0xf7,
	0x3e, 0x34, 0xc3, 0x5a, 0xe1, 0x71, 0x21, 0xdf, 0x60, 0xbe, 0x68, 0xc6, 0xff, 0xd5, 0x29, 0xe6,
	0x31, 0x5a, 0x57, 0x00, 0xfa, 0x71, 0x3e, 0x6b, 0xd4, 0x69, 0xd0, 0xf6, 0x83, 0x56, 0xc8, 0xfb,
	0xe3, 0xbb, 0xd7, 0xb7, 0x07, 0xfd, 0xc0, 0x1c, 0xd5, 0x55, 0xde, 0x7c, 0x79, 0xe2, 0x51, 0x58,
	0x42, 0x21, 0x12, 0x9e, 0x03, 0x1e, 0x62, 0xfa, 0x5d, 0x0b, 0x78, 0x78, 0x3b, 0xb9, 0xd7, 0x8d,
	0x41, 0x28, 0x92, 0xb5, 0xe3, 0x4f, 0xb6, 0xbb, 0x8b, 0x99, 0xd6, 0x2b, 0x98, 0x6a, 0x97, 0x67,
	0x8d, 0x49, 0x24, 0xe6, 0x34, 0x8d, 0x90, 0xc8, 0x23, 0x1c, 0xf1, 0x74, 0xb9, 0x20, 0xb5, 0xf2,
	0x2c, 0x86, 0xee, 0xaf, 0x76, 0xa8, 0xc0, 0xa3, 0xfe, 0xcb, 0x75, 0x9c, 0x27, 0xbd, 0xe7, 0x4a,
	0x43, 0x87, 0x01, 0x92, 0xb0, 0x46, 0x9e, 0x8d, 0x9a, 0x60, 0x78, 0x56, 0x14, 0xdf, 0x61, 0xe0,
	0xb7, 0x14, 0xdf, 0xb3, 0x7d, 0x45, 0x79, 0xd7, 0x07, 0xf2, 0x9c, 0x10, 0x87, 0x01, 0x21, 0x2d,
	0x89, 0x10, 0xcf, 0x26, 0x44, 0xd1, 0xae, 0x7e, 0x2f, 0xe7, 0xdc, 0xff, 0x0c, 0x00, 0x00, 0xff,
	0xff, 0xc4, 0x1d, 0x81, 0x8a, 0x11, 0x03, 0x00, 0x00,
}
