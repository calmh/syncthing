// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: lib/config/deviceconfiguration.proto

package config

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	github_com_syncthing_syncthing_lib_protocol "github.com/syncthing/syncthing/lib/protocol"
	protocol "github.com/syncthing/syncthing/lib/protocol"
	_ "github.com/syncthing/syncthing/proto/ext"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type DeviceConfiguration struct {
	DeviceID                               github_com_syncthing_syncthing_lib_protocol.DeviceID `protobuf:"bytes,1,opt,name=device_id,json=deviceId,proto3,customtype=github.com/syncthing/syncthing/lib/protocol.DeviceID" json:"deviceID" xml:"id,attr"`
	Name                                   string                                               `protobuf:"bytes,2,opt,name=name,proto3" json:"name" xml:"name"`
	Addresses                              []string                                             `protobuf:"bytes,3,rep,name=addresses,proto3" json:"addresses" xml:"address" default:"dynamic"`
	Compression                            protocol.Compression                                 `protobuf:"varint,4,opt,name=compression,proto3,enum=protocol.Compression" json:"compression" xml:"compression"`
	CertName                               string                                               `protobuf:"bytes,5,opt,name=cert_name,json=certName,proto3" json:"certName" xml:"certName"`
	Introducer                             bool                                                 `protobuf:"varint,6,opt,name=introducer,proto3" json:"introducer" xml:"introducer"`
	SkipIntroductionRemovals               bool                                                 `protobuf:"varint,7,opt,name=skip_introduction_removals,json=skipIntroductionRemovals,proto3" json:"skipIntroductionRemovals" xml:"skipIntroductionRemovals"`
	IntroducedBy                           github_com_syncthing_syncthing_lib_protocol.DeviceID `protobuf:"bytes,8,opt,name=introduced_by,json=introducedBy,proto3,customtype=github.com/syncthing/syncthing/lib/protocol.DeviceID" json:"introducedBy" xml:"introducedBy"`
	Paused                                 bool                                                 `protobuf:"varint,9,opt,name=paused,proto3" json:"paused" xml:"paused"`
	AllowedNetworks                        []string                                             `protobuf:"bytes,10,rep,name=allowed_networks,json=allowedNetworks,proto3" json:"allowedNetworks" xml:"allowedNetwork"`
	AutoAcceptFolders                      bool                                                 `protobuf:"varint,11,opt,name=auto_accept_folders,json=autoAcceptFolders,proto3" json:"autoAcceptFolders" xml:"autoAcceptFolders"`
	MaxSendKbps                            int                                                  `protobuf:"varint,12,opt,name=max_send_kbps,json=maxSendKbps,proto3,casttype=int" json:"maxSendKbps" xml:"maxSendKbps"`
	MaxRecvKbps                            int                                                  `protobuf:"varint,13,opt,name=max_recv_kbps,json=maxRecvKbps,proto3,casttype=int" json:"maxRecvKbps" xml:"maxRecvKbps"`
	IgnoredFolders                         []ObservedFolder                                     `protobuf:"bytes,14,rep,name=ignored_folders,json=ignoredFolders,proto3" json:"ignoredFolders" xml:"ignoredFolder"`
	PendingFolders                         []ObservedFolder                                     `protobuf:"bytes,15,rep,name=pending_folders,json=pendingFolders,proto3" json:"pendingFolders" xml:"pendingFolder"`
	MaxRequestKiB                          int                                                  `protobuf:"varint,16,opt,name=max_request_kib,json=maxRequestKib,proto3,casttype=int" json:"maxRequestKiB" xml:"maxRequestKiB,omitempty"`
	DeprecatedNameAttr                     string                                               `protobuf:"bytes,9001,opt,name=name_attr,json=nameAttr,proto3" json:"-" xml:"name,attr,omitempty"`                                                                                          // Deprecated: Do not use.
	DeprecatedCompressionAttr              protocol.Compression                                 `protobuf:"varint,9002,opt,name=compression_attr,json=compressionAttr,proto3,enum=protocol.Compression" json:"-" xml:"compression,attr,omitempty"`                                          // Deprecated: Do not use.
	DeprecatedCertNameAttr                 string                                               `protobuf:"bytes,9003,opt,name=cert_name_attr,json=certNameAttr,proto3" json:"-" xml:"certName,attr,omitempty"`                                                                             // Deprecated: Do not use.
	DeprecatedIntroducerAttr               bool                                                 `protobuf:"varint,9004,opt,name=introducer_attr,json=introducerAttr,proto3" json:"-" xml:"introducer,attr,omitempty"`                                                                       // Deprecated: Do not use.
	DeprecatedSkipIntroductionRemovalsAttr bool                                                 `protobuf:"varint,9005,opt,name=skip_introduction_removals_attr,json=skipIntroductionRemovalsAttr,proto3" json:"-" xml:"skipIntroductionRemovals,attr,omitempty"`                           // Deprecated: Do not use.
	DeprecatedIntroducedByAttr             github_com_syncthing_syncthing_lib_protocol.DeviceID `protobuf:"bytes,9006,opt,name=introduced_by_attr,json=introducedByAttr,proto3,customtype=github.com/syncthing/syncthing/lib/protocol.DeviceID" json:"-" xml:"introducedBy,attr,omitempty"` // Deprecated: Do not use.
}

func (m *DeviceConfiguration) Reset()         { *m = DeviceConfiguration{} }
func (m *DeviceConfiguration) String() string { return proto.CompactTextString(m) }
func (*DeviceConfiguration) ProtoMessage()    {}
func (*DeviceConfiguration) Descriptor() ([]byte, []int) {
	return fileDescriptor_744b782bd13071dd, []int{0}
}
func (m *DeviceConfiguration) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeviceConfiguration.Unmarshal(m, b)
}
func (m *DeviceConfiguration) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeviceConfiguration.Marshal(b, m, deterministic)
}
func (m *DeviceConfiguration) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeviceConfiguration.Merge(m, src)
}
func (m *DeviceConfiguration) XXX_Size() int {
	return xxx_messageInfo_DeviceConfiguration.Size(m)
}
func (m *DeviceConfiguration) XXX_DiscardUnknown() {
	xxx_messageInfo_DeviceConfiguration.DiscardUnknown(m)
}

var xxx_messageInfo_DeviceConfiguration proto.InternalMessageInfo

func init() {
	proto.RegisterType((*DeviceConfiguration)(nil), "config.DeviceConfiguration")
}

func init() {
	proto.RegisterFile("lib/config/deviceconfiguration.proto", fileDescriptor_744b782bd13071dd)
}

var fileDescriptor_744b782bd13071dd = []byte{
	// 1135 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x56, 0xcf, 0x6f, 0x1b, 0x45,
	0x14, 0xf6, 0x90, 0x36, 0x8d, 0x27, 0xb1, 0x9d, 0x4e, 0x4a, 0xb2, 0x31, 0xd4, 0x63, 0x56, 0x16,
	0x18, 0x68, 0x1c, 0x14, 0x38, 0xb5, 0x02, 0xa9, 0xdb, 0x28, 0x28, 0xaa, 0x68, 0xc5, 0xf4, 0x04,
	0x52, 0xb5, 0xac, 0x77, 0x27, 0xc9, 0x2a, 0xf6, 0xee, 0xb2, 0xbb, 0x4e, 0x63, 0x89, 0x23, 0x12,
	0x48, 0x50, 0x09, 0xf5, 0x2f, 0xe0, 0x9a, 0xf2, 0x43, 0xfc, 0x19, 0xbd, 0x39, 0x12, 0x12, 0x02,
	0x0e, 0x23, 0xd5, 0xb9, 0xa0, 0x3d, 0xee, 0xb1, 0x27, 0xb4, 0x33, 0xeb, 0xf5, 0xac, 0x6b, 0x47,
	0x91, 0x7a, 0x9b, 0xfd, 0xbe, 0x37, 0xdf, 0xfb, 0xde, 0xf3, 0xcb, 0x9b, 0xc0, 0x46, 0xc7, 0x6e,
	0x6f, 0x9a, 0xae, 0xb3, 0x67, 0xef, 0x6f, 0x5a, 0xf4, 0xc8, 0x36, 0xa9, 0xf8, 0xe8, 0xf9, 0x46,
	0x68, 0xbb, 0x4e, 0xcb, 0xf3, 0xdd, 0xd0, 0x45, 0xf3, 0x02, 0xac, 0xae, 0x26, 0xd1, 0x1c, 0x32,
	0xdd, 0xce, 0x66, 0x9b, 0x7a, 0x82, 0xaf, 0xae, 0x4b, 0x2a, 0x6e, 0x3b, 0xa0, 0xfe, 0x11, 0xb5,
	0x52, 0xaa, 0x48, 0x8f, 0x43, 0x71, 0x54, 0xff, 0xbc, 0x06, 0x57, 0xb6, 0x79, 0x8e, 0x3b, 0x72,
	0x0e, 0xf4, 0x14, 0xc0, 0xa2, 0xc8, 0xad, 0xdb, 0x96, 0x02, 0xea, 0xa0, 0xb9, 0xa4, 0xfd, 0x08,
	0x9e, 0x31, 0x5c, 0xf8, 0x97, 0xe1, 0x8f, 0xf6, 0xed, 0xf0, 0xa0, 0xd7, 0x6e, 0x99, 0x6e, 0x77,
	0x33, 0xe8, 0x3b, 0x66, 0x78, 0x60, 0x3b, 0xfb, 0xd2, 0x49, 0x76, 0xd4, 0x12, 0xea, 0xbb, 0xdb,
	0x43, 0x86, 0x17, 0x46, 0xe7, 0x88, 0xe1, 0x05, 0x2b, 0x3d, 0xc7, 0x0c, 0x97, 0x8e, 0xbb, 0x9d,
	0x9b, 0xaa, 0x6d, 0xdd, 0x30, 0xc2, 0xd0, 0x57, 0xa3, 0x41, 0xe3, 0x4a, 0x7a, 0x8e, 0x07, 0x8d,
	0x2c, 0xee, 0xfb, 0xd3, 0x06, 0x78, 0x72, 0xda, 0xc8, 0x34, 0xc8, 0x88, 0xb1, 0xd0, 0x7b, 0xf0,
	0x92, 0x63, 0x74, 0xa9, 0xf2, 0x5a, 0x1d, 0x34, 0x8b, 0xda, 0x6a, 0xc4, 0x30, 0xff, 0x8e, 0x19,
	0x86, 0x5c, 0x39, 0xf9, 0x50, 0x09, 0xc7, 0xd0, 0x17, 0xb0, 0x68, 0x58, 0x96, 0x4f, 0x83, 0x80,
	0x06, 0xca, 0x5c, 0x7d, 0xae, 0x59, 0xd4, 0x6e, 0x45, 0x0c, 0x8f, 0xc1, 0x98, 0x61, 0xcc, 0x6f,
	0xa5, 0x88, 0x5a, 0xb7, 0xe8, 0x9e, 0xd1, 0xeb, 0x84, 0x37, 0x55, 0xab, 0xef, 0x18, 0x5d, 0xdb,
	0x54, 0x5f, 0x0c, 0x1a, 0x57, 0xd2, 0x33, 0x19, 0x5f, 0x44, 0x0f, 0xe1, 0xa2, 0xe9, 0x76, 0xbd,
	0xe4, 0xcb, 0x76, 0x1d, 0xe5, 0x52, 0x1d, 0x34, 0xcb, 0x5b, 0xaf, 0xb7, 0xb2, 0x4e, 0xdc, 0x19,
	0x93, 0x5a, 0x23, 0x62, 0x58, 0x8e, 0x8e, 0x19, 0xbe, 0xca, 0xb3, 0x4a, 0x98, 0x4a, 0xe4, 0x08,
	0x74, 0x0b, 0x16, 0x4d, 0xea, 0x87, 0x3a, 0x2f, 0xf5, 0x32, 0x2f, 0xb5, 0x96, 0x34, 0x32, 0x01,
	0xef, 0x89, 0x72, 0xcb, 0x42, 0x22, 0x05, 0x54, 0x92, 0x71, 0x48, 0x83, 0xd0, 0x76, 0x42, 0xdf,
	0xb5, 0x7a, 0x26, 0xf5, 0x95, 0xf9, 0x3a, 0x68, 0x2e, 0x68, 0x6a, 0xc4, 0xb0, 0x84, 0xc6, 0x0c,
	0x2f, 0x8b, 0x1f, 0x22, 0x83, 0x54, 0x22, 0xf1, 0xe8, 0x1b, 0x58, 0x0d, 0x0e, 0x6d, 0x4f, 0x1f,
	0x41, 0xc9, 0xa0, 0xe8, 0x3e, 0xed, 0xba, 0x47, 0x46, 0x27, 0x50, 0xae, 0x70, 0xcd, 0x4f, 0x22,
	0x86, 0x95, 0x24, 0x6a, 0x57, 0x0a, 0x22, 0x69, 0x4c, 0xcc, 0x70, 0x8d, 0x67, 0x98, 0x15, 0xa0,
	0x92, 0x99, 0x77, 0xd1, 0x0f, 0x00, 0x96, 0x32, 0x33, 0x96, 0xde, 0xee, 0x2b, 0x0b, 0x7c, 0x2a,
	0xf7, 0x5e, 0x65, 0x28, 0x23, 0x86, 0x97, 0xc6, 0xa2, 0x5a, 0x3f, 0x66, 0x18, 0xe5, 0x7b, 0x60,
	0x69, 0x7d, 0x35, 0x19, 0x3d, 0x92, 0x8b, 0x43, 0x5b, 0x70, 0xde, 0x33, 0x7a, 0x01, 0xb5, 0x94,
	0x22, 0xaf, 0xbb, 0x1a, 0x31, 0x9c, 0x22, 0x31, 0xc3, 0x4b, 0x5c, 0x43, 0x7c, 0xaa, 0x24, 0xc5,
	0xd1, 0x01, 0x5c, 0x36, 0x3a, 0x1d, 0xf7, 0x11, 0xb5, 0x74, 0x87, 0x86, 0x8f, 0x5c, 0xff, 0x30,
	0x50, 0x20, 0x9f, 0xc0, 0x8f, 0x23, 0x86, 0x2b, 0x29, 0x77, 0x2f, 0xa5, 0x62, 0x86, 0xaf, 0x89,
	0x39, 0xcc, 0xe1, 0xc9, 0x9f, 0x47, 0x39, 0x0f, 0x91, 0xc9, 0xab, 0xe8, 0x2b, 0xb8, 0x62, 0xf4,
	0x42, 0x57, 0x37, 0x4c, 0x93, 0x7a, 0xa1, 0xbe, 0xe7, 0x76, 0x2c, 0xea, 0x07, 0xca, 0x22, 0xb7,
	0xfa, 0x41, 0xc4, 0xf0, 0xd5, 0x84, 0xbe, 0xcd, 0xd9, 0x1d, 0x41, 0xc6, 0x0c, 0xaf, 0x89, 0x74,
	0x93, 0x8c, 0x4a, 0x5e, 0x8e, 0x46, 0xf7, 0x61, 0xa9, 0x6b, 0x1c, 0xeb, 0x01, 0x75, 0x2c, 0xfd,
	0xb0, 0xed, 0x05, 0xca, 0x52, 0x1d, 0x34, 0x2f, 0x6b, 0xef, 0x27, 0x63, 0xdd, 0x35, 0x8e, 0x1f,
	0x50, 0xc7, 0xba, 0xdb, 0xf6, 0x82, 0x6c, 0xac, 0x25, 0x4c, 0x7d, 0xc1, 0xf0, 0x9c, 0xed, 0x84,
	0x44, 0x0e, 0x1c, 0x09, 0xfa, 0xd4, 0x3c, 0x12, 0x82, 0xa5, 0x9c, 0x20, 0xa1, 0xe6, 0xd1, 0xa4,
	0xe0, 0x08, 0xcb, 0x09, 0x8e, 0x40, 0xe4, 0xc0, 0x8a, 0xbd, 0xef, 0xb8, 0x3e, 0xb5, 0xb2, 0xfa,
	0xcb, 0xf5, 0xb9, 0xe6, 0xe2, 0xd6, 0x6a, 0x4b, 0x6c, 0xc5, 0xd6, 0xfd, 0x74, 0x2b, 0x8a, 0x9a,
	0xb4, 0x8d, 0x64, 0x90, 0x22, 0x86, 0xcb, 0xe9, 0xb5, 0x71, 0x63, 0x56, 0xc4, 0x48, 0xc8, 0xb0,
	0x4a, 0x26, 0xc2, 0x92, 0x7c, 0x1e, 0x75, 0x2c, 0xdb, 0xd9, 0xcf, 0xf2, 0x55, 0x2e, 0x96, 0x2f,
	0xbd, 0x36, 0x99, 0x2f, 0x07, 0xab, 0x64, 0x22, 0x0c, 0xfd, 0x01, 0x60, 0x45, 0x74, 0xec, 0xeb,
	0x1e, 0x0d, 0x42, 0xfd, 0xd0, 0x6e, 0x2b, 0xcb, 0xbc, 0x67, 0xdf, 0x81, 0x21, 0xc3, 0xa5, 0xcf,
	0x92, 0x5e, 0x70, 0xea, 0xae, 0xad, 0x45, 0x0c, 0x97, 0xba, 0x32, 0x10, 0x33, 0x7c, 0x7d, 0xdc,
	0xc7, 0x11, 0x7a, 0xc3, 0xed, 0xda, 0x21, 0xed, 0x7a, 0x61, 0x7f, 0xd4, 0xd3, 0x68, 0xd0, 0x58,
	0x9b, 0x11, 0x12, 0x0f, 0x1a, 0x79, 0xcd, 0x27, 0xa7, 0x8d, 0x7c, 0x56, 0x92, 0xe3, 0xdb, 0xc8,
	0x85, 0xc5, 0x64, 0x79, 0xe9, 0xc9, 0x6a, 0x57, 0x4e, 0x76, 0xf8, 0x0a, 0xfb, 0x7c, 0xc8, 0x30,
	0xda, 0xa6, 0x9e, 0x4f, 0x4d, 0x23, 0xa4, 0x56, 0xb2, 0xac, 0x6e, 0x87, 0xa1, 0x1f, 0x31, 0x0c,
	0x36, 0x62, 0x86, 0xd7, 0xb3, 0x05, 0xce, 0x1f, 0x04, 0xc9, 0x60, 0x34, 0x68, 0xac, 0x4c, 0xc1,
	0x15, 0x40, 0x16, 0x9c, 0x54, 0x04, 0x9d, 0x00, 0xb8, 0x2c, 0xad, 0x50, 0x91, 0xf8, 0xe9, 0xce,
	0x79, 0x8b, 0x79, 0x6f, 0xc8, 0xf0, 0xfa, 0xd8, 0x8f, 0x44, 0xc9, 0xb6, 0xea, 0x93, 0xbb, 0x7a,
	0x8a, 0xbb, 0xea, 0x6c, 0x5a, 0x01, 0xa4, 0x62, 0xe6, 0x95, 0xd1, 0xb7, 0x00, 0x96, 0xb3, 0xfd,
	0x2e, 0x9c, 0xfe, 0x22, 0x5a, 0xf4, 0x70, 0xc8, 0xf0, 0xaa, 0x64, 0x29, 0xdd, 0xe9, 0xb2, 0x9f,
	0xeb, 0xb9, 0xc5, 0x3f, 0xc5, 0xcc, 0xda, 0x0c, 0x4e, 0x01, 0x64, 0xc9, 0x94, 0x04, 0xd1, 0x63,
	0x00, 0x2b, 0xe3, 0x9d, 0x2f, 0x7c, 0xfc, 0xba, 0xc3, 0x17, 0x87, 0x39, 0x64, 0x58, 0x19, 0xfb,
	0xd8, 0xcd, 0xa2, 0x64, 0x27, 0x78, 0xe2, 0x09, 0x99, 0xe2, 0x65, 0x7d, 0x26, 0xab, 0x00, 0x52,
	0xb6, 0x73, 0xb2, 0xe8, 0x2f, 0x00, 0xf1, 0xec, 0x57, 0x47, 0xf8, 0xfb, 0x4d, 0xf8, 0x7b, 0x9c,
	0xcc, 0xfd, 0xdb, 0x63, 0x83, 0x0f, 0x66, 0x3c, 0x25, 0xb2, 0xdd, 0x8d, 0x73, 0xdf, 0xa3, 0x29,
	0xe6, 0xdf, 0xb9, 0x60, 0xac, 0x02, 0xc8, 0x9b, 0xc1, 0x39, 0x06, 0xd0, 0x3f, 0x00, 0xa2, 0xdc,
	0x7b, 0x26, 0x6a, 0xf9, 0x7d, 0x87, 0xbf, 0x6a, 0x27, 0xaf, 0xfa, 0xbf, 0x56, 0x75, 0xca, 0x0f,
	0x65, 0x69, 0x7d, 0xb9, 0xf6, 0xb7, 0x5e, 0x7a, 0xe9, 0xa6, 0xd4, 0xfb, 0xc6, 0x39, 0x7c, 0xf2,
	0x2e, 0x2a, 0x80, 0x2c, 0xdb, 0x13, 0x09, 0xb4, 0x4f, 0x9f, 0x3d, 0xaf, 0x15, 0x4e, 0x9f, 0xd7,
	0x0a, 0xff, 0x0d, 0x6b, 0x85, 0x9f, 0xce, 0x6a, 0x85, 0x9f, 0xcf, 0x6a, 0xe0, 0xf4, 0xac, 0x56,
	0xf8, 0xfb, 0xac, 0x56, 0xf8, 0xf2, 0xdd, 0x0b, 0x94, 0x24, 0x56, 0x66, 0x7b, 0x9e, 0x97, 0xf6,
	0xe1, 0xff, 0x01, 0x00, 0x00, 0xff, 0xff, 0x3d, 0xd6, 0xab, 0x8e, 0x13, 0x0b, 0x00, 0x00,
}

func (m *DeviceConfiguration) ProtoSize() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.DeviceID.ProtoSize()
	n += 1 + l + sovDeviceconfiguration(uint64(l))
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovDeviceconfiguration(uint64(l))
	}
	if len(m.Addresses) > 0 {
		for _, s := range m.Addresses {
			l = len(s)
			n += 1 + l + sovDeviceconfiguration(uint64(l))
		}
	}
	if m.Compression != 0 {
		n += 1 + sovDeviceconfiguration(uint64(m.Compression))
	}
	l = len(m.CertName)
	if l > 0 {
		n += 1 + l + sovDeviceconfiguration(uint64(l))
	}
	if m.Introducer {
		n += 2
	}
	if m.SkipIntroductionRemovals {
		n += 2
	}
	l = m.IntroducedBy.ProtoSize()
	n += 1 + l + sovDeviceconfiguration(uint64(l))
	if m.Paused {
		n += 2
	}
	if len(m.AllowedNetworks) > 0 {
		for _, s := range m.AllowedNetworks {
			l = len(s)
			n += 1 + l + sovDeviceconfiguration(uint64(l))
		}
	}
	if m.AutoAcceptFolders {
		n += 2
	}
	if m.MaxSendKbps != 0 {
		n += 1 + sovDeviceconfiguration(uint64(m.MaxSendKbps))
	}
	if m.MaxRecvKbps != 0 {
		n += 1 + sovDeviceconfiguration(uint64(m.MaxRecvKbps))
	}
	if len(m.IgnoredFolders) > 0 {
		for _, e := range m.IgnoredFolders {
			l = e.ProtoSize()
			n += 1 + l + sovDeviceconfiguration(uint64(l))
		}
	}
	if len(m.PendingFolders) > 0 {
		for _, e := range m.PendingFolders {
			l = e.ProtoSize()
			n += 1 + l + sovDeviceconfiguration(uint64(l))
		}
	}
	if m.MaxRequestKiB != 0 {
		n += 2 + sovDeviceconfiguration(uint64(m.MaxRequestKiB))
	}
	l = len(m.DeprecatedNameAttr)
	if l > 0 {
		n += 3 + l + sovDeviceconfiguration(uint64(l))
	}
	if m.DeprecatedCompressionAttr != 0 {
		n += 3 + sovDeviceconfiguration(uint64(m.DeprecatedCompressionAttr))
	}
	l = len(m.DeprecatedCertNameAttr)
	if l > 0 {
		n += 3 + l + sovDeviceconfiguration(uint64(l))
	}
	if m.DeprecatedIntroducerAttr {
		n += 4
	}
	if m.DeprecatedSkipIntroductionRemovalsAttr {
		n += 4
	}
	l = m.DeprecatedIntroducedByAttr.ProtoSize()
	n += 3 + l + sovDeviceconfiguration(uint64(l))
	return n
}

func sovDeviceconfiguration(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozDeviceconfiguration(x uint64) (n int) {
	return sovDeviceconfiguration(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
