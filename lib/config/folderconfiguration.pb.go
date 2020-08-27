// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: lib/config/folderconfiguration.proto

package config

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	fs "github.com/syncthing/syncthing/lib/fs"
	github_com_syncthing_syncthing_lib_protocol "github.com/syncthing/syncthing/lib/protocol"
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

type FolderDeviceConfiguration struct {
	DeviceID     github_com_syncthing_syncthing_lib_protocol.DeviceID `protobuf:"bytes,1,opt,name=device_id,json=deviceId,proto3,customtype=github.com/syncthing/syncthing/lib/protocol.DeviceID" json:"deviceID" xml:"id,attr"`
	IntroducedBy github_com_syncthing_syncthing_lib_protocol.DeviceID `protobuf:"bytes,2,opt,name=introduced_by,json=introducedBy,proto3,customtype=github.com/syncthing/syncthing/lib/protocol.DeviceID" json:"introducedBy" xml:"introducedBy,attr"`
}

func (m *FolderDeviceConfiguration) Reset()         { *m = FolderDeviceConfiguration{} }
func (m *FolderDeviceConfiguration) String() string { return proto.CompactTextString(m) }
func (*FolderDeviceConfiguration) ProtoMessage()    {}
func (*FolderDeviceConfiguration) Descriptor() ([]byte, []int) {
	return fileDescriptor_44a9785876ed3afa, []int{0}
}
func (m *FolderDeviceConfiguration) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FolderDeviceConfiguration.Unmarshal(m, b)
}
func (m *FolderDeviceConfiguration) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FolderDeviceConfiguration.Marshal(b, m, deterministic)
}
func (m *FolderDeviceConfiguration) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FolderDeviceConfiguration.Merge(m, src)
}
func (m *FolderDeviceConfiguration) XXX_Size() int {
	return xxx_messageInfo_FolderDeviceConfiguration.Size(m)
}
func (m *FolderDeviceConfiguration) XXX_DiscardUnknown() {
	xxx_messageInfo_FolderDeviceConfiguration.DiscardUnknown(m)
}

var xxx_messageInfo_FolderDeviceConfiguration proto.InternalMessageInfo

type FolderConfiguration struct {
	ID                      string                      `protobuf:"bytes,1,opt,name=id,proto3" json:"id" xml:"id,attr"`
	Label                   string                      `protobuf:"bytes,2,opt,name=label,proto3" json:"label" xml:"label,omitempty" restart:"false"`
	FilesystemType          fs.FilesystemType           `protobuf:"varint,3,opt,name=filesystem_type,json=filesystemType,proto3,enum=fs.FilesystemType" json:"filesystemType" xml:"filesystemType"`
	Path                    string                      `protobuf:"bytes,4,opt,name=path,proto3" json:"path" xml:"path"`
	Type                    FolderType                  `protobuf:"varint,5,opt,name=type,proto3,enum=config.FolderType" json:"type" xml:"type"`
	Devices                 []FolderDeviceConfiguration `protobuf:"bytes,6,rep,name=devices,proto3" json:"devices" xml:"device"`
	RescanIntervalS         int                         `protobuf:"varint,7,opt,name=rescan_interval_s,json=rescanIntervalS,proto3,casttype=int" json:"rescanIntervalS" xml:"rescanIntervalS" default:"3600"`
	FSWatcherEnabled        bool                        `protobuf:"varint,8,opt,name=fs_watcher_enabled,json=fsWatcherEnabled,proto3" json:"fsWatcherEnabled" xml:"fsWatcherEnabled" default:"true"`
	FSWatcherDelayS         int                         `protobuf:"varint,9,opt,name=fs_watcher_delay_s,json=fsWatcherDelayS,proto3,casttype=int" json:"fsWatcherDelayS" xml:"fsWatcherDelayS" default:"10"`
	IgnorePerms             bool                        `protobuf:"varint,10,opt,name=ignore_perms,json=ignorePerms,proto3" json:"ignorePerms" xml:"ignorePerms"`
	AutoNormalize           bool                        `protobuf:"varint,11,opt,name=auto_normalize,json=autoNormalize,proto3" json:"autoNormalize" xml:"autoNormalize" default:"true"`
	MinDiskFree             Size                        `protobuf:"bytes,12,opt,name=min_disk_free,json=minDiskFree,proto3" json:"minDiskFree" xml:"minDiskFree"`
	Versioning              VersioningConfiguration     `protobuf:"bytes,13,opt,name=versioning,proto3" json:"versioning" xml:"versioning"`
	Copiers                 int                         `protobuf:"varint,14,opt,name=copiers,proto3,casttype=int" json:"copiers" xml:"copiers"`
	PullerMaxPendingKiB     int                         `protobuf:"varint,15,opt,name=puller_max_pending_kib,json=pullerMaxPendingKib,proto3,casttype=int" json:"pullerMaxPendingKiB" xml:"pullerMaxPendingKiB"`
	Hashers                 int                         `protobuf:"varint,16,opt,name=hashers,proto3,casttype=int" json:"hashers" xml:"hashers"`
	Order                   PullOrder                   `protobuf:"varint,17,opt,name=order,proto3,enum=config.PullOrder" json:"order" xml:"order"`
	IgnoreDelete            bool                        `protobuf:"varint,18,opt,name=ignore_delete,json=ignoreDelete,proto3" json:"ignoreDelete" xml:"ignoreDelete"`
	ScanProgressIntervalS   int                         `protobuf:"varint,19,opt,name=scan_progress_interval_s,json=scanProgressIntervalS,proto3,casttype=int" json:"scanProgressIntervalS" xml:"scanProgressIntervalS"`
	PullerPauseS            int                         `protobuf:"varint,20,opt,name=puller_pause_s,json=pullerPauseS,proto3,casttype=int" json:"pullerPauseS" xml:"pullerPauseS"`
	MaxConflicts            int                         `protobuf:"varint,21,opt,name=max_conflicts,json=maxConflicts,proto3,casttype=int" json:"maxConflicts" xml:"maxConflicts" default:"-1"`
	DisableSparseFiles      bool                        `protobuf:"varint,22,opt,name=disable_sparse_files,json=disableSparseFiles,proto3" json:"disableSparseFiles" xml:"disableSparseFiles"`
	DisableTempIndexes      bool                        `protobuf:"varint,23,opt,name=disable_temp_indexes,json=disableTempIndexes,proto3" json:"disableTempIndexes" xml:"disableTempIndexes"`
	Paused                  bool                        `protobuf:"varint,24,opt,name=paused,proto3" json:"paused" xml:"paused"`
	WeakHashThresholdPct    int                         `protobuf:"varint,25,opt,name=weak_hash_threshold_pct,json=weakHashThresholdPct,proto3,casttype=int" json:"weakHashThresholdPct" xml:"weakHashThresholdPct"`
	MarkerName              string                      `protobuf:"bytes,26,opt,name=marker_name,json=markerName,proto3" json:"markerName" xml:"markerName"`
	CopyOwnershipFromParent bool                        `protobuf:"varint,27,opt,name=copy_ownership_from_parent,json=copyOwnershipFromParent,proto3" json:"copyOwnershipFromParent" xml:"copyOwnershipFromParent"`
	RawModTimeWindowS       int                         `protobuf:"varint,28,opt,name=mod_time_window_s,json=modTimeWindowS,proto3,casttype=int" json:"modTimeWindowS" xml:"modTimeWindowS"`
	MaxConcurrentWrites     int                         `protobuf:"varint,29,opt,name=max_concurrent_writes,json=maxConcurrentWrites,proto3,casttype=int" json:"maxConcurrentWrites" xml:"maxConcurrentWrites" default:"2"`
	DisableFsync            bool                        `protobuf:"varint,30,opt,name=disable_fsync,json=disableFsync,proto3" json:"disableFsync" xml:"disableFsync"`
	BlockPullOrder          BlockPullOrder              `protobuf:"varint,31,opt,name=block_pull_order,json=blockPullOrder,proto3,enum=config.BlockPullOrder" json:"blockPullOrder" xml:"blockPullOrder"`
	CopyRangeMethod         fs.CopyRangeMethod          `protobuf:"varint,32,opt,name=copy_range_method,json=copyRangeMethod,proto3,enum=fs.CopyRangeMethod" json:"copyRangeMethod" xml:"copyRangeMethod" default:"standard"`
	CaseSensitiveFS         bool                        `protobuf:"varint,33,opt,name=case_sensitive_fs,json=caseSensitiveFs,proto3" json:"caseSensitiveFS" xml:"caseSensitiveFS"`
	JunctionsAsDirs         bool                        `protobuf:"varint,34,opt,name=follow_junctions,json=followJunctions,proto3" json:"junctionsAsDirs" xml:"junctionsAsDirs"`
	// Legacy deprecated
	DeprecatedReadOnly             bool       `protobuf:"varint,9000,opt,name=read_only,json=readOnly,proto3" json:"-" xml:"ro,attr,omitempty"`                                                    // Deprecated: Do not use.
	DeprecatedMinDiskFreePct       float64    `protobuf:"fixed64,9001,opt,name=min_disk_free_pct,json=minDiskFreePct,proto3" json:"-" xml:"minDiskFreePct,omitempty"`                              // Deprecated: Do not use.
	DeprecatedPullers              int        `protobuf:"varint,9002,opt,name=pullers,proto3,casttype=int" json:"-" xml:"pullers,omitempty"`                                                       // Deprecated: Do not use.
	DeprecatedLabelAttr            string     `protobuf:"bytes,9003,opt,name=label_attr,json=labelAttr,proto3" json:"-" xml:"label,attr,omitempty"`                                                // Deprecated: Do not use.
	DeprecatedPathAttr             string     `protobuf:"bytes,9004,opt,name=path_attr,json=pathAttr,proto3" json:"-" xml:"path,attr,omitempty"`                                                   // Deprecated: Do not use.
	DeprecatedTypeAttr             FolderType `protobuf:"varint,9005,opt,name=type_attr,json=typeAttr,proto3,enum=config.FolderType" json:"-" xml:"type,attr,omitempty"`                           // Deprecated: Do not use.
	DeprecatedRescanIntervalSAttr  int        `protobuf:"varint,9006,opt,name=rescan_interval_s_attr,json=rescanIntervalSAttr,proto3,casttype=int" json:"-" xml:"rescanIntervalS,attr,omitempty"`  // Deprecated: Do not use.
	DeprecatedFsWatcherEnabledAttr bool       `protobuf:"varint,9007,opt,name=fs_watcher_enabled_attr,json=fsWatcherEnabledAttr,proto3" json:"-" xml:"fsWatcherEnabled,attr,omitempty"`            // Deprecated: Do not use.
	DeprecatedFsWatcherDelaySAttr  int        `protobuf:"varint,9008,opt,name=fs_watcher_delay_s_attr,json=fsWatcherDelaySAttr,proto3,casttype=int" json:"-" xml:"fsWatcherDelayS,attr,omitempty"` // Deprecated: Do not use.
	DeprecatedIgnorePermsAttr      bool       `protobuf:"varint,9009,opt,name=ignore_perms_attr,json=ignorePermsAttr,proto3" json:"-" xml:"ignorePerms,attr,omitempty"`                            // Deprecated: Do not use.
	DeprecatedAutoNormalizeAttr    bool       `protobuf:"varint,9010,opt,name=auto_normalize_attr,json=autoNormalizeAttr,proto3" json:"-" xml:"autoNormalize,attr,omitempty"`                      // Deprecated: Do not use.
}

func (m *FolderConfiguration) Reset()         { *m = FolderConfiguration{} }
func (m *FolderConfiguration) String() string { return proto.CompactTextString(m) }
func (*FolderConfiguration) ProtoMessage()    {}
func (*FolderConfiguration) Descriptor() ([]byte, []int) {
	return fileDescriptor_44a9785876ed3afa, []int{1}
}
func (m *FolderConfiguration) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FolderConfiguration.Unmarshal(m, b)
}
func (m *FolderConfiguration) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FolderConfiguration.Marshal(b, m, deterministic)
}
func (m *FolderConfiguration) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FolderConfiguration.Merge(m, src)
}
func (m *FolderConfiguration) XXX_Size() int {
	return xxx_messageInfo_FolderConfiguration.Size(m)
}
func (m *FolderConfiguration) XXX_DiscardUnknown() {
	xxx_messageInfo_FolderConfiguration.DiscardUnknown(m)
}

var xxx_messageInfo_FolderConfiguration proto.InternalMessageInfo

func init() {
	proto.RegisterType((*FolderDeviceConfiguration)(nil), "config.FolderDeviceConfiguration")
	proto.RegisterType((*FolderConfiguration)(nil), "config.FolderConfiguration")
}

func init() {
	proto.RegisterFile("lib/config/folderconfiguration.proto", fileDescriptor_44a9785876ed3afa)
}

var fileDescriptor_44a9785876ed3afa = []byte{
	// 2263 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x58, 0xcd, 0x6f, 0x1c, 0x49,
	0x15, 0x77, 0x7b, 0xf3, 0x61, 0x97, 0x3f, 0xa7, 0xec, 0xc4, 0x15, 0x67, 0x77, 0xca, 0xdb, 0x0c,
	0x2b, 0x83, 0x36, 0x4e, 0xe2, 0x20, 0x0e, 0x11, 0xbb, 0x90, 0x89, 0x77, 0xc0, 0x64, 0xb3, 0x19,
	0xda, 0xd1, 0x46, 0x2c, 0xa0, 0x56, 0x7b, 0xba, 0xc6, 0xd3, 0xeb, 0xfe, 0x18, 0xba, 0x6a, 0x62,
	0x4f, 0x2e, 0x04, 0x21, 0x21, 0x90, 0xf6, 0x80, 0xcc, 0x81, 0x13, 0x12, 0x12, 0x08, 0xc1, 0xb2,
	0xcb, 0x2e, 0xfc, 0x15, 0x7b, 0x41, 0xf1, 0x11, 0x21, 0x51, 0xd2, 0x4e, 0x2e, 0x68, 0x6e, 0xcc,
	0x31, 0x27, 0x54, 0x55, 0xdd, 0x3d, 0xd5, 0x3d, 0x9d, 0x80, 0xc4, 0xc5, 0x9a, 0xfa, 0xfd, 0x5e,
	0xbd, 0xf7, 0xab, 0xd7, 0x55, 0xaf, 0xea, 0x19, 0xd4, 0x7c, 0x6f, 0xff, 0x6a, 0x2b, 0x0a, 0xdb,
	0xde, 0xc1, 0xd5, 0x76, 0xe4, 0xbb, 0x24, 0x56, 0x83, 0x5e, 0xec, 0x30, 0x2f, 0x0a, 0xb7, 0xba,
	0x71, 0xc4, 0x22, 0x78, 0x4e, 0x81, 0xeb, 0x97, 0x27, 0xac, 0x59, 0xbf, 0x4b, 0x94, 0xd1, 0xfa,
	0x05, 0x8d, 0xa4, 0xde, 0xa3, 0x14, 0x5e, 0xd7, 0xe0, 0x6e, 0xcf, 0xf7, 0xa3, 0xd8, 0x25, 0x71,
	0xc2, 0x6d, 0x6a, 0xdc, 0x43, 0x12, 0x53, 0x2f, 0x0a, 0xbd, 0xf0, 0xa0, 0x44, 0xc1, 0x3a, 0xd6,
	0x2c, 0xf7, 0xfd, 0xa8, 0x75, 0x58, 0x74, 0x05, 0x85, 0x41, 0x9b, 0x5e, 0x15, 0x82, 0x68, 0x82,
	0xbd, 0x9c, 0x60, 0xad, 0xa8, 0xdb, 0x8f, 0x9d, 0xf0, 0x80, 0x04, 0x84, 0x75, 0x22, 0x37, 0x61,
	0x67, 0xc9, 0x31, 0x53, 0x3f, 0xcd, 0x7f, 0x4f, 0x83, 0x4b, 0x0d, 0xb9, 0x9e, 0x1d, 0xf2, 0xd0,
	0x6b, 0x91, 0xdb, 0xba, 0x02, 0xf8, 0xa1, 0x01, 0x66, 0x5d, 0x89, 0xdb, 0x9e, 0x8b, 0x8c, 0x0d,
	0x63, 0x73, 0xbe, 0xfe, 0x81, 0xf1, 0x19, 0xc7, 0x53, 0xff, 0xe0, 0xf8, 0x2b, 0x07, 0x1e, 0xeb,
	0xf4, 0xf6, 0xb7, 0x5a, 0x51, 0x70, 0x95, 0xf6, 0xc3, 0x16, 0xeb, 0x78, 0xe1, 0x81, 0xf6, 0x4b,
	0x48, 0x90, 0x41, 0x5a, 0x91, 0xbf, 0xa5, 0xbc, 0xef, 0xee, 0x0c, 0x38, 0x9e, 0x49, 0x7f, 0x0f,
	0x39, 0x9e, 0x71, 0x93, 0xdf, 0x23, 0x8e, 0x17, 0x8e, 0x03, 0xff, 0xa6, 0xe9, 0xb9, 0xaf, 0x3b,
	0x8c, 0xc5, 0xe6, 0xf0, 0x49, 0xed, 0x7c, 0xf2, 0x7b, 0xf4, 0xa4, 0x96, 0xd9, 0xfd, 0xec, 0xb4,
	0x66, 0x9c, 0x9c, 0xd6, 0x32, 0x1f, 0x56, 0xca, 0xb8, 0xf0, 0xf7, 0x06, 0x58, 0xf0, 0x42, 0x16,
	0x47, 0x6e, 0xaf, 0x45, 0x5c, 0x7b, 0xbf, 0x8f, 0xa6, 0xa5, 0xe0, 0xc7, 0xff, 0x97, 0xe0, 0x21,
	0xc7, 0xf3, 0x63, 0xaf, 0xf5, 0xfe, 0x88, 0xe3, 0x35, 0x25, 0x54, 0x03, 0x33, 0xc9, 0x95, 0x09,
	0x54, 0x08, 0xb6, 0x72, 0x1e, 0xcc, 0x7f, 0x6e, 0x82, 0x15, 0x95, 0xf3, 0x7c, 0xb6, 0xdf, 0x04,
	0xd3, 0x49, 0x96, 0x67, 0xeb, 0x5b, 0x03, 0x8e, 0xa7, 0x65, 0xf4, 0x69, 0xcf, 0x7d, 0x51, 0x72,
	0x4e, 0x4e, 0x6b, 0xd3, 0xbb, 0x3b, 0xd6, 0xb4, 0xe7, 0x42, 0x1b, 0x9c, 0xf5, 0x9d, 0x7d, 0xe2,
	0xcb, 0x75, 0xcf, 0xd6, 0x77, 0x87, 0x1c, 0x2b, 0x60, 0xc4, 0xf1, 0x17, 0xe5, 0x7c, 0x39, 0x7a,
	0x3d, 0x0a, 0x3c, 0x46, 0x82, 0x2e, 0xeb, 0x9b, 0x1b, 0x31, 0xa1, 0xcc, 0x89, 0xd9, 0x4d, 0xb3,
	0xed, 0xf8, 0x94, 0x08, 0xbf, 0x4b, 0x05, 0x9b, 0xc7, 0xa7, 0xb5, 0x29, 0x4b, 0xb9, 0x81, 0x07,
	0x60, 0xa9, 0xed, 0xf9, 0x84, 0xf6, 0x29, 0x23, 0x81, 0x2d, 0xf6, 0x1b, 0x7a, 0x69, 0xc3, 0xd8,
	0x5c, 0xdc, 0x86, 0x5b, 0x6d, 0xba, 0xd5, 0xc8, 0xa8, 0xfb, 0xfd, 0x2e, 0xa9, 0x7f, 0x79, 0xc8,
	0xf1, 0x62, 0x3b, 0x87, 0x8d, 0x38, 0x5e, 0x95, 0x3a, 0xf2, 0xb0, 0x69, 0x15, 0xec, 0xe0, 0x36,
	0x38, 0xd3, 0x75, 0x58, 0x07, 0x9d, 0x91, 0x0b, 0xa9, 0x0e, 0x39, 0x96, 0xe3, 0x11, 0xc7, 0x40,
	0xce, 0x17, 0x03, 0x21, 0x56, 0xa2, 0x96, 0xfc, 0x0b, 0x1b, 0xe0, 0x8c, 0x54, 0x74, 0x36, 0x51,
	0xa4, 0x8e, 0xcc, 0x96, 0x4a, 0xb4, 0x54, 0x24, 0xfd, 0x30, 0xa5, 0x43, 0xf9, 0x11, 0x03, 0xe9,
	0x47, 0xfc, 0xb0, 0xe4, 0x5f, 0xf8, 0x7d, 0x70, 0x5e, 0x6d, 0x29, 0x8a, 0xce, 0x6d, 0xbc, 0xb4,
	0x39, 0xb7, 0xfd, 0x6a, 0xde, 0x55, 0xc9, 0x39, 0xa9, 0x63, 0xb1, 0xc3, 0x86, 0x1c, 0xa7, 0x33,
	0x47, 0x1c, 0xcf, 0xcb, 0x00, 0x6a, 0x6c, 0x5a, 0x29, 0x01, 0x7f, 0x6e, 0x80, 0x4a, 0x4c, 0x68,
	0xcb, 0x09, 0x6d, 0x2f, 0x64, 0x24, 0x7e, 0xe8, 0xf8, 0x36, 0x45, 0xe7, 0x37, 0x8c, 0xcd, 0xb3,
	0xf5, 0x1f, 0x0c, 0x39, 0x5e, 0x52, 0xe4, 0x6e, 0xc2, 0xed, 0x8d, 0x38, 0xae, 0x49, 0x4f, 0x05,
	0xdc, 0xdc, 0x70, 0x49, 0xdb, 0xe9, 0xf9, 0xec, 0xa6, 0x79, 0xe3, 0xab, 0xd7, 0xae, 0x99, 0xcf,
	0x38, 0x7e, 0xc9, 0x0b, 0x99, 0xf8, 0x80, 0x05, 0xcb, 0x67, 0x4f, 0x6a, 0x67, 0x84, 0x89, 0x55,
	0x24, 0xe0, 0xa7, 0x06, 0x80, 0x6d, 0x6a, 0x1f, 0x39, 0xac, 0xd5, 0x21, 0xb1, 0x4d, 0x42, 0x67,
	0xdf, 0x27, 0x2e, 0x9a, 0xd9, 0x30, 0x36, 0x67, 0xea, 0x3f, 0x36, 0x06, 0x1c, 0x2f, 0x37, 0xf6,
	0x1e, 0x28, 0xf6, 0x2d, 0x45, 0x0e, 0x39, 0x5e, 0x6e, 0xd3, 0x3c, 0x96, 0xed, 0xae, 0x22, 0xa1,
	0x69, 0x64, 0x71, 0x4f, 0x26, 0x7a, 0x62, 0xb2, 0x50, 0x27, 0xc8, 0x93, 0xd3, 0xda, 0x44, 0x30,
	0x6b, 0xc2, 0x1a, 0x7e, 0x94, 0x97, 0xec, 0x12, 0xdf, 0xe9, 0xdb, 0x14, 0xcd, 0xca, 0xfc, 0xfd,
	0x68, 0xc0, 0xf1, 0x52, 0xe6, 0x64, 0x47, 0x70, 0x7b, 0x22, 0xa5, 0x99, 0x17, 0x05, 0x8d, 0x38,
	0x36, 0xf3, 0x7a, 0x15, 0xae, 0xc9, 0xbd, 0x9e, 0x4b, 0x68, 0xc1, 0xee, 0xd9, 0x93, 0xda, 0xf4,
	0xf5, 0x6b, 0x27, 0xa7, 0xb5, 0x62, 0x2c, 0xab, 0x68, 0x09, 0xdf, 0x05, 0xf3, 0xde, 0x41, 0x18,
	0xc5, 0xc4, 0xee, 0x92, 0x38, 0xa0, 0x08, 0xc8, 0xd4, 0xde, 0x18, 0x72, 0x3c, 0xa7, 0xf0, 0xa6,
	0x80, 0x47, 0x1c, 0x57, 0xd4, 0xf1, 0x1e, 0x63, 0x22, 0x59, 0xba, 0x8d, 0xa5, 0x0f, 0xe0, 0x11,
	0x58, 0x74, 0x7a, 0x2c, 0xb2, 0xc3, 0x28, 0x0e, 0x1c, 0xdf, 0x7b, 0x44, 0xd0, 0x9c, 0xf4, 0xdc,
	0x1c, 0x72, 0xbc, 0x20, 0x98, 0x77, 0x52, 0x22, 0x5b, 0x6c, 0x0e, 0x2d, 0xf9, 0x32, 0xf9, 0x69,
	0xe9, 0x67, 0xb1, 0xf2, 0x30, 0x7c, 0x0f, 0x2c, 0x04, 0x5e, 0x68, 0xbb, 0x1e, 0x3d, 0xb4, 0xdb,
	0x31, 0x21, 0x68, 0x7e, 0xc3, 0xd8, 0x9c, 0xdb, 0x9e, 0x4f, 0x8f, 0xc8, 0x9e, 0xf7, 0x88, 0xd4,
	0x37, 0x93, 0xd3, 0x30, 0x17, 0x78, 0xe1, 0x8e, 0x47, 0x0f, 0x1b, 0x31, 0x21, 0xd9, 0x1a, 0x35,
	0xcc, 0xb4, 0x74, 0x0b, 0x78, 0x00, 0xc0, 0xf8, 0x26, 0x44, 0x0b, 0xd2, 0x31, 0x4e, 0x1d, 0xbf,
	0x9b, 0x31, 0xf9, 0x93, 0xf7, 0x5a, 0x12, 0x4b, 0x9b, 0x3a, 0xe2, 0x78, 0x59, 0x86, 0x1a, 0x43,
	0xa6, 0xa5, 0xf1, 0xf0, 0x0d, 0x70, 0xbe, 0x15, 0x75, 0x3d, 0x12, 0x53, 0xb4, 0x28, 0x37, 0xce,
	0x17, 0xc4, 0xd1, 0x4d, 0xa0, 0xac, 0xd6, 0x26, 0xe3, 0x74, 0x23, 0x58, 0xa9, 0x01, 0xfc, 0x9b,
	0x01, 0x2e, 0x8a, 0x3b, 0x98, 0xc4, 0x76, 0xe0, 0x1c, 0xdb, 0x5d, 0x12, 0xba, 0x5e, 0x78, 0x60,
	0x1f, 0x7a, 0xfb, 0x68, 0x49, 0xba, 0xfb, 0x95, 0x38, 0x3a, 0x2b, 0x4d, 0x69, 0x72, 0xd7, 0x39,
	0x6e, 0x2a, 0x83, 0x3b, 0x5e, 0x7d, 0xc8, 0xf1, 0x4a, 0x77, 0x12, 0x1e, 0x71, 0x7c, 0x49, 0x95,
	0xb5, 0x49, 0x4e, 0xdb, 0x87, 0xa5, 0x53, 0xcb, 0xe1, 0x93, 0xd3, 0x5a, 0x59, 0x7c, 0xab, 0xc4,
	0x76, 0x5f, 0xa4, 0xa3, 0xe3, 0xd0, 0x8e, 0x48, 0xc7, 0xf2, 0x38, 0x1d, 0x09, 0x94, 0xa5, 0x23,
	0x19, 0x8f, 0xd3, 0x91, 0x00, 0xf0, 0x16, 0x38, 0x2b, 0x5f, 0x23, 0xa8, 0x22, 0x0b, 0x6f, 0x25,
	0xfd, 0x62, 0x22, 0xfe, 0x3d, 0x41, 0xd4, 0x91, 0xb8, 0x88, 0xa4, 0xcd, 0x88, 0xe3, 0x39, 0xe9,
	0x4d, 0x8e, 0x4c, 0x4b, 0xa1, 0xf0, 0x0e, 0x58, 0x48, 0x8e, 0x89, 0x4b, 0x7c, 0xc2, 0x08, 0x82,
	0x72, 0x37, 0xbf, 0x26, 0xef, 0x5e, 0x49, 0xec, 0x48, 0x7c, 0xc4, 0x31, 0xd4, 0x0e, 0x8a, 0x02,
	0x4d, 0x2b, 0x67, 0x03, 0x8f, 0x01, 0x92, 0xe5, 0xb5, 0x1b, 0x47, 0x07, 0x31, 0xa1, 0x54, 0xaf,
	0xb3, 0x2b, 0x72, 0x7d, 0x5f, 0x1f, 0x72, 0x7c, 0x41, 0xd8, 0x34, 0x13, 0x13, 0xbd, 0xda, 0x5e,
	0x96, 0x01, 0x4a, 0xd9, 0x6c, 0xed, 0xe5, 0x93, 0xe1, 0x1e, 0x58, 0x4c, 0xf6, 0x45, 0xd7, 0xe9,
	0x51, 0x62, 0x53, 0xb4, 0x2a, 0xe3, 0x5d, 0x11, 0xeb, 0x50, 0x4c, 0x53, 0x10, 0x7b, 0xd9, 0x3a,
	0x74, 0x30, 0xf3, 0x9e, 0x33, 0x85, 0x04, 0x2c, 0x88, 0x5d, 0x26, 0x92, 0xea, 0x7b, 0x2d, 0x46,
	0xd1, 0x05, 0xe9, 0xf3, 0x1b, 0xc2, 0x67, 0xe0, 0x1c, 0xdf, 0x4e, 0xf1, 0x11, 0xc7, 0x58, 0x1d,
	0x30, 0x0d, 0xd4, 0xce, 0xf9, 0x95, 0xeb, 0x69, 0x00, 0x51, 0xbf, 0xae, 0x5c, 0xb7, 0x72, 0xb3,
	0xa1, 0x0b, 0x56, 0x5d, 0x8f, 0x8a, 0x22, 0x6b, 0xd3, 0xae, 0x13, 0x53, 0x62, 0xcb, 0x2b, 0x19,
	0x5d, 0x94, 0x5f, 0x62, 0x7b, 0xc8, 0x31, 0x4c, 0xf8, 0x3d, 0x49, 0xcb, 0xcb, 0x7e, 0xc4, 0x31,
	0x52, 0xd7, 0xdc, 0x04, 0x65, 0x5a, 0x25, 0xf6, 0x7a, 0x14, 0xf1, 0xba, 0xb0, 0xbd, 0xd0, 0x25,
	0xc7, 0x84, 0xa2, 0xb5, 0x89, 0x28, 0xf7, 0x49, 0xd0, 0xdd, 0x55, 0x6c, 0x31, 0x8a, 0x46, 0x8d,
	0xa3, 0x68, 0x20, 0xdc, 0x06, 0xe7, 0xe4, 0x07, 0x70, 0x11, 0x92, 0x7e, 0xd7, 0x87, 0x1c, 0x27,
	0x48, 0x76, 0x31, 0xab, 0xa1, 0x69, 0x25, 0x38, 0x64, 0x60, 0xed, 0x88, 0x38, 0x87, 0xb6, 0xd8,
	0xd5, 0x36, 0xeb, 0xc4, 0x84, 0x76, 0x22, 0xdf, 0xb5, 0xbb, 0x2d, 0x86, 0x2e, 0xc9, 0x84, 0xbf,
	0x31, 0xe4, 0x78, 0x55, 0x98, 0x7c, 0xcb, 0xa1, 0x9d, 0xfb, 0xa9, 0x41, 0xb3, 0xc5, 0x46, 0x1c,
	0xaf, 0x4b, 0x97, 0x65, 0x64, 0xf6, 0x51, 0x4b, 0xa7, 0xc2, 0xdb, 0x60, 0x2e, 0x70, 0xe2, 0x43,
	0x12, 0xdb, 0xa1, 0x13, 0x10, 0xb4, 0x2e, 0x9f, 0x3b, 0xa6, 0x28, 0x67, 0x0a, 0x7e, 0xc7, 0x09,
	0x48, 0x56, 0xce, 0xc6, 0x90, 0x69, 0x69, 0x3c, 0xec, 0x83, 0x75, 0xf1, 0xcc, 0xb7, 0xa3, 0xa3,
	0x90, 0xc4, 0xb4, 0xe3, 0x75, 0xed, 0x76, 0x1c, 0x05, 0x76, 0xd7, 0x89, 0x49, 0xc8, 0xd0, 0x65,
	0x99, 0x82, 0xaf, 0x0d, 0x39, 0x5e, 0x13, 0x56, 0xf7, 0x52, 0xa3, 0x46, 0x1c, 0x05, 0x4d, 0x69,
	0x32, 0xe2, 0xf8, 0x95, 0xb4, 0xe2, 0x95, 0xf1, 0xa6, 0xf5, 0xbc, 0x99, 0xf0, 0xa7, 0x06, 0xa8,
	0x04, 0x91, 0x6b, 0x33, 0x2f, 0x20, 0xf6, 0x91, 0x17, 0xba, 0xd1, 0x91, 0x4d, 0xd1, 0xcb, 0x32,
	0x61, 0xdf, 0x1b, 0x70, 0x5c, 0xb1, 0x9c, 0xa3, 0xbb, 0x91, 0x7b, 0xdf, 0x0b, 0xc8, 0x03, 0xc9,
	0x8a, 0xfb, 0x78, 0x31, 0xc8, 0x21, 0xd9, 0xa3, 0x30, 0x0f, 0xa7, 0x99, 0x3b, 0x39, 0xad, 0x4d,
	0x7a, 0xb1, 0x0a, 0x3e, 0xe0, 0x63, 0x03, 0x5c, 0x48, 0x8e, 0x49, 0xab, 0x17, 0x0b, 0x6d, 0xf6,
	0x51, 0xec, 0x31, 0x42, 0xd1, 0x2b, 0x52, 0xcc, 0xdb, 0xa2, 0xf4, 0xaa, 0x0d, 0x9f, 0xf0, 0x0f,
	0x24, 0x9d, 0xbd, 0x5d, 0x4a, 0x38, 0xed, 0xf0, 0x6c, 0x6b, 0x67, 0xc7, 0xd8, 0xb6, 0xca, 0x3c,
	0x89, 0x22, 0x96, 0xee, 0xed, 0xb6, 0xe8, 0x29, 0x50, 0x75, 0x5c, 0xc4, 0x12, 0xa2, 0x21, 0xf0,
	0xec, 0xf0, 0xeb, 0xa0, 0x69, 0xe5, 0x6c, 0xa0, 0x0f, 0x96, 0x65, 0xaf, 0x67, 0x8b, 0x5a, 0x60,
	0xab, 0xfa, 0x8a, 0x65, 0x7d, 0xbd, 0x98, 0xd6, 0xd7, 0xba, 0xe0, 0xc7, 0x45, 0x56, 0x3e, 0xb7,
	0xf7, 0x73, 0x58, 0x96, 0xd9, 0x3c, 0x6c, 0x5a, 0x05, 0x3b, 0xf8, 0x81, 0x01, 0x2a, 0x72, 0x0b,
	0xc9, 0x56, 0xd1, 0x56, 0xbd, 0x22, 0xda, 0x90, 0xf1, 0x56, 0xc4, 0xd3, 0xfe, 0x76, 0xd4, 0xed,
	0x5b, 0x82, 0xbb, 0x2b, 0xa9, 0xfa, 0x1d, 0xf1, 0xac, 0x6a, 0xe5, 0xc1, 0x11, 0xc7, 0x9b, 0xd9,
	0x36, 0xd2, 0x70, 0x2d, 0x8d, 0x94, 0x39, 0xa1, 0xeb, 0xc4, 0xae, 0xf9, 0xec, 0x49, 0x6d, 0x26,
	0x1d, 0x58, 0x45, 0x47, 0xf0, 0x77, 0x42, 0x8e, 0x23, 0x0a, 0x28, 0x09, 0xa9, 0xc7, 0xbc, 0x87,
	0x22, 0xa3, 0xe8, 0x55, 0x99, 0xce, 0x63, 0xf1, 0xc6, 0xbb, 0xed, 0x50, 0xb2, 0x97, 0x72, 0x0d,
	0xf9, 0xc6, 0x6b, 0xe5, 0xa1, 0x11, 0xc7, 0x17, 0x94, 0x98, 0x3c, 0x2e, 0x3b, 0x9c, 0xa2, 0xed,
	0x24, 0x24, 0x1e, 0x77, 0x85, 0x20, 0x56, 0xc1, 0x86, 0xc2, 0xdf, 0x1a, 0x60, 0xb9, 0x1d, 0xf9,
	0x7e, 0x74, 0x64, 0xbf, 0xdf, 0x0b, 0x5b, 0xe2, 0x39, 0x42, 0x91, 0x39, 0x56, 0xf9, 0xed, 0x14,
	0xbc, 0x45, 0x77, 0xbc, 0x98, 0x0a, 0x95, 0xef, 0xe7, 0xa1, 0x4c, 0x65, 0x01, 0x97, 0x2a, 0x8b,
	0xb6, 0x93, 0x90, 0x50, 0x59, 0x08, 0x62, 0x2d, 0x29, 0x45, 0x19, 0x0c, 0x0f, 0xc1, 0x6c, 0x4c,
	0x1c, 0xd7, 0x8e, 0x42, 0xbf, 0x8f, 0xfe, 0xd0, 0x90, 0xf2, 0xee, 0x0e, 0x38, 0x86, 0x3b, 0xa4,
	0x1b, 0x93, 0x96, 0xc3, 0x88, 0x6b, 0x11, 0xc7, 0xbd, 0x17, 0xfa, 0xfd, 0x21, 0xc7, 0xc6, 0x95,
	0xac, 0xbf, 0x8d, 0x23, 0xd9, 0x5f, 0x6a, 0xdd, 0xa2, 0xe8, 0x6f, 0x27, 0x50, 0x64, 0x58, 0x33,
	0x71, 0xe2, 0x00, 0xfe, 0x10, 0x54, 0x72, 0xcf, 0x43, 0x59, 0x3f, 0xff, 0x28, 0x82, 0x1a, 0xf5,
	0xb7, 0x06, 0x1c, 0xa3, 0x71, 0xd0, 0xbb, 0xe3, 0x97, 0x5f, 0xb3, 0xc5, 0xd2, 0xd0, 0xd5, 0xe2,
	0x1b, 0xb1, 0xd9, 0x62, 0x9a, 0x02, 0x64, 0x58, 0x8b, 0x79, 0x12, 0x7e, 0x17, 0x9c, 0x57, 0xf7,
	0x25, 0x45, 0x1f, 0x36, 0xe4, 0x59, 0x7f, 0x53, 0x14, 0x9e, 0x71, 0x20, 0xf5, 0x0e, 0xa2, 0xf9,
	0xc5, 0x25, 0x53, 0x34, 0xd7, 0xc9, 0x01, 0x47, 0x86, 0x95, 0xfa, 0x83, 0x0c, 0x00, 0xd9, 0xf7,
	0xda, 0x62, 0xc5, 0xe8, 0x4f, 0x0d, 0x59, 0x9d, 0xef, 0x8b, 0xb7, 0xdd, 0xd8, 0xfb, 0xdb, 0xc2,
	0xe0, 0x16, 0x63, 0x71, 0xea, 0x7f, 0x5d, 0x6b, 0xb4, 0x27, 0xf3, 0xb7, 0x5a, 0x46, 0x20, 0xc3,
	0x9a, 0xf5, 0x53, 0x3f, 0x30, 0x02, 0xb3, 0xa2, 0x9f, 0x55, 0x41, 0x3f, 0x52, 0x41, 0xbf, 0x93,
	0xff, 0x60, 0x4d, 0x87, 0x75, 0xf4, 0x98, 0x97, 0xb2, 0xa6, 0xb8, 0x24, 0xe4, 0x4a, 0x09, 0x2e,
	0x3e, 0x5a, 0x37, 0x71, 0x02, 0x7f, 0x62, 0x80, 0x59, 0xd1, 0xf9, 0xaa, 0x88, 0x1f, 0x37, 0x9e,
	0xdb, 0x3f, 0x17, 0x54, 0x08, 0x6c, 0x52, 0x85, 0xf0, 0x54, 0xa6, 0xa2, 0x04, 0x17, 0x2a, 0x58,
	0xe2, 0x04, 0x7e, 0x62, 0x80, 0x8b, 0x13, 0x8d, 0xb1, 0x92, 0xf4, 0x67, 0xf5, 0x5d, 0xfb, 0x03,
	0x8e, 0x5f, 0xd1, 0x77, 0x6d, 0xae, 0x9b, 0xd5, 0x95, 0x94, 0x76, 0xcc, 0x45, 0x51, 0xe3, 0x87,
	0x75, 0xf5, 0xc5, 0x96, 0xc8, 0xb0, 0x56, 0xe2, 0xc9, 0x60, 0xf0, 0x63, 0x03, 0xac, 0x4d, 0xb6,
	0xcf, 0x4a, 0xf2, 0x27, 0xea, 0xa0, 0xf5, 0x06, 0x1c, 0x57, 0xc7, 0x92, 0x1b, 0x85, 0x76, 0x56,
	0xd7, 0x5c, 0xde, 0x42, 0x97, 0x64, 0x12, 0xff, 0x17, 0x1b, 0x64, 0x58, 0xab, 0xed, 0x92, 0x40,
	0xa2, 0xdd, 0x5f, 0x9b, 0xec, 0x9d, 0x95, 0xde, 0x4f, 0x4b, 0x53, 0xdc, 0xc8, 0xb7, 0xb3, 0x93,
	0x29, 0x2e, 0xf4, 0xbb, 0x2f, 0x48, 0xf1, 0x8b, 0x2d, 0x45, 0x8a, 0xdb, 0x93, 0xc1, 0xe0, 0x2f,
	0x0d, 0x50, 0xd1, 0x1b, 0x68, 0x25, 0xf6, 0x2f, 0x2a, 0xb9, 0xed, 0x01, 0xc7, 0x97, 0xc6, 0x62,
	0x77, 0xc7, 0xfd, 0xb1, 0x2e, 0x74, 0xa3, 0xd8, 0x59, 0x97, 0xa4, 0x74, 0xfd, 0xf9, 0x34, 0x32,
	0xac, 0x25, 0x2f, 0xef, 0x19, 0xfe, 0xda, 0x00, 0x2b, 0xf9, 0xf6, 0x5b, 0xe9, 0xfa, 0xab, 0xd2,
	0xe5, 0x0f, 0x38, 0xbe, 0x3c, 0xd6, 0x75, 0x4b, 0x6f, 0xa0, 0x75, 0x65, 0x25, 0x7d, 0x79, 0x89,
	0xb6, 0x97, 0x5f, 0x64, 0x80, 0x0c, 0xab, 0xe2, 0x14, 0xfd, 0xd7, 0xbf, 0xf9, 0xd9, 0xe7, 0xd5,
	0xa9, 0xd3, 0xcf, 0xab, 0x53, 0xff, 0x1a, 0x54, 0xa7, 0x7e, 0xf1, 0xb4, 0x3a, 0xf5, 0x9b, 0xa7,
	0x55, 0xe3, 0xf4, 0x69, 0x75, 0xea, 0xef, 0x4f, 0xab, 0x53, 0xef, 0x7d, 0xe9, 0x7f, 0xf8, 0x2f,
	0xa8, 0x3a, 0xfb, 0xfb, 0xe7, 0xe4, 0x7f, 0x43, 0x6f, 0xfc, 0x27, 0x00, 0x00, 0xff, 0xff, 0xd2,
	0xf0, 0x88, 0x4a, 0x2b, 0x17, 0x00, 0x00,
}

func (m *FolderDeviceConfiguration) ProtoSize() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.DeviceID.ProtoSize()
	n += 1 + l + sovFolderconfiguration(uint64(l))
	l = m.IntroducedBy.ProtoSize()
	n += 1 + l + sovFolderconfiguration(uint64(l))
	return n
}

func (m *FolderConfiguration) ProtoSize() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ID)
	if l > 0 {
		n += 1 + l + sovFolderconfiguration(uint64(l))
	}
	l = len(m.Label)
	if l > 0 {
		n += 1 + l + sovFolderconfiguration(uint64(l))
	}
	if m.FilesystemType != 0 {
		n += 1 + sovFolderconfiguration(uint64(m.FilesystemType))
	}
	l = len(m.Path)
	if l > 0 {
		n += 1 + l + sovFolderconfiguration(uint64(l))
	}
	if m.Type != 0 {
		n += 1 + sovFolderconfiguration(uint64(m.Type))
	}
	if len(m.Devices) > 0 {
		for _, e := range m.Devices {
			l = e.ProtoSize()
			n += 1 + l + sovFolderconfiguration(uint64(l))
		}
	}
	if m.RescanIntervalS != 0 {
		n += 1 + sovFolderconfiguration(uint64(m.RescanIntervalS))
	}
	if m.FSWatcherEnabled {
		n += 2
	}
	if m.FSWatcherDelayS != 0 {
		n += 1 + sovFolderconfiguration(uint64(m.FSWatcherDelayS))
	}
	if m.IgnorePerms {
		n += 2
	}
	if m.AutoNormalize {
		n += 2
	}
	l = m.MinDiskFree.ProtoSize()
	n += 1 + l + sovFolderconfiguration(uint64(l))
	l = m.Versioning.ProtoSize()
	n += 1 + l + sovFolderconfiguration(uint64(l))
	if m.Copiers != 0 {
		n += 1 + sovFolderconfiguration(uint64(m.Copiers))
	}
	if m.PullerMaxPendingKiB != 0 {
		n += 1 + sovFolderconfiguration(uint64(m.PullerMaxPendingKiB))
	}
	if m.Hashers != 0 {
		n += 2 + sovFolderconfiguration(uint64(m.Hashers))
	}
	if m.Order != 0 {
		n += 2 + sovFolderconfiguration(uint64(m.Order))
	}
	if m.IgnoreDelete {
		n += 3
	}
	if m.ScanProgressIntervalS != 0 {
		n += 2 + sovFolderconfiguration(uint64(m.ScanProgressIntervalS))
	}
	if m.PullerPauseS != 0 {
		n += 2 + sovFolderconfiguration(uint64(m.PullerPauseS))
	}
	if m.MaxConflicts != 0 {
		n += 2 + sovFolderconfiguration(uint64(m.MaxConflicts))
	}
	if m.DisableSparseFiles {
		n += 3
	}
	if m.DisableTempIndexes {
		n += 3
	}
	if m.Paused {
		n += 3
	}
	if m.WeakHashThresholdPct != 0 {
		n += 2 + sovFolderconfiguration(uint64(m.WeakHashThresholdPct))
	}
	l = len(m.MarkerName)
	if l > 0 {
		n += 2 + l + sovFolderconfiguration(uint64(l))
	}
	if m.CopyOwnershipFromParent {
		n += 3
	}
	if m.RawModTimeWindowS != 0 {
		n += 2 + sovFolderconfiguration(uint64(m.RawModTimeWindowS))
	}
	if m.MaxConcurrentWrites != 0 {
		n += 2 + sovFolderconfiguration(uint64(m.MaxConcurrentWrites))
	}
	if m.DisableFsync {
		n += 3
	}
	if m.BlockPullOrder != 0 {
		n += 2 + sovFolderconfiguration(uint64(m.BlockPullOrder))
	}
	if m.CopyRangeMethod != 0 {
		n += 2 + sovFolderconfiguration(uint64(m.CopyRangeMethod))
	}
	if m.CaseSensitiveFS {
		n += 3
	}
	if m.JunctionsAsDirs {
		n += 3
	}
	if m.DeprecatedReadOnly {
		n += 4
	}
	if m.DeprecatedMinDiskFreePct != 0 {
		n += 11
	}
	if m.DeprecatedPullers != 0 {
		n += 3 + sovFolderconfiguration(uint64(m.DeprecatedPullers))
	}
	l = len(m.DeprecatedLabelAttr)
	if l > 0 {
		n += 3 + l + sovFolderconfiguration(uint64(l))
	}
	l = len(m.DeprecatedPathAttr)
	if l > 0 {
		n += 3 + l + sovFolderconfiguration(uint64(l))
	}
	if m.DeprecatedTypeAttr != 0 {
		n += 3 + sovFolderconfiguration(uint64(m.DeprecatedTypeAttr))
	}
	if m.DeprecatedRescanIntervalSAttr != 0 {
		n += 3 + sovFolderconfiguration(uint64(m.DeprecatedRescanIntervalSAttr))
	}
	if m.DeprecatedFsWatcherEnabledAttr {
		n += 4
	}
	if m.DeprecatedFsWatcherDelaySAttr != 0 {
		n += 3 + sovFolderconfiguration(uint64(m.DeprecatedFsWatcherDelaySAttr))
	}
	if m.DeprecatedIgnorePermsAttr {
		n += 4
	}
	if m.DeprecatedAutoNormalizeAttr {
		n += 4
	}
	return n
}

func sovFolderconfiguration(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozFolderconfiguration(x uint64) (n int) {
	return sovFolderconfiguration(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
