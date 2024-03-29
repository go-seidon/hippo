// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: api/grpcapp/health.proto

package grpcapp

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CheckHealthParam struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CheckHealthParam) Reset() {
	*x = CheckHealthParam{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_grpcapp_health_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckHealthParam) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckHealthParam) ProtoMessage() {}

func (x *CheckHealthParam) ProtoReflect() protoreflect.Message {
	mi := &file_api_grpcapp_health_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckHealthParam.ProtoReflect.Descriptor instead.
func (*CheckHealthParam) Descriptor() ([]byte, []int) {
	return file_api_grpcapp_health_proto_rawDescGZIP(), []int{0}
}

type CheckHealthResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    int32            `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message string           `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Data    *CheckHealthData `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *CheckHealthResult) Reset() {
	*x = CheckHealthResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_grpcapp_health_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckHealthResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckHealthResult) ProtoMessage() {}

func (x *CheckHealthResult) ProtoReflect() protoreflect.Message {
	mi := &file_api_grpcapp_health_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckHealthResult.ProtoReflect.Descriptor instead.
func (*CheckHealthResult) Descriptor() ([]byte, []int) {
	return file_api_grpcapp_health_proto_rawDescGZIP(), []int{1}
}

func (x *CheckHealthResult) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *CheckHealthResult) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *CheckHealthResult) GetData() *CheckHealthData {
	if x != nil {
		return x.Data
	}
	return nil
}

type CheckHealthData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status  string                        `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	Details map[string]*CheckHealthDetail `protobuf:"bytes,2,rep,name=details,proto3" json:"details,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *CheckHealthData) Reset() {
	*x = CheckHealthData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_grpcapp_health_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckHealthData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckHealthData) ProtoMessage() {}

func (x *CheckHealthData) ProtoReflect() protoreflect.Message {
	mi := &file_api_grpcapp_health_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckHealthData.ProtoReflect.Descriptor instead.
func (*CheckHealthData) Descriptor() ([]byte, []int) {
	return file_api_grpcapp_health_proto_rawDescGZIP(), []int{2}
}

func (x *CheckHealthData) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *CheckHealthData) GetDetails() map[string]*CheckHealthDetail {
	if x != nil {
		return x.Details
	}
	return nil
}

type CheckHealthDetail struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name      string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Status    string `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	CheckedAt int64  `protobuf:"varint,3,opt,name=checked_at,json=checkedAt,proto3" json:"checked_at,omitempty"`
	Error     string `protobuf:"bytes,4,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *CheckHealthDetail) Reset() {
	*x = CheckHealthDetail{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_grpcapp_health_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckHealthDetail) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckHealthDetail) ProtoMessage() {}

func (x *CheckHealthDetail) ProtoReflect() protoreflect.Message {
	mi := &file_api_grpcapp_health_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckHealthDetail.ProtoReflect.Descriptor instead.
func (*CheckHealthDetail) Descriptor() ([]byte, []int) {
	return file_api_grpcapp_health_proto_rawDescGZIP(), []int{3}
}

func (x *CheckHealthDetail) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CheckHealthDetail) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *CheckHealthDetail) GetCheckedAt() int64 {
	if x != nil {
		return x.CheckedAt
	}
	return 0
}

func (x *CheckHealthDetail) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_api_grpcapp_health_proto protoreflect.FileDescriptor

var file_api_grpcapp_health_proto_rawDesc = []byte{
	0x0a, 0x18, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x61, 0x70, 0x70, 0x2f, 0x68, 0x65,
	0x61, 0x6c, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x68, 0x65, 0x61, 0x6c,
	0x74, 0x68, 0x2e, 0x76, 0x31, 0x22, 0x12, 0x0a, 0x10, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x48, 0x65,
	0x61, 0x6c, 0x74, 0x68, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x22, 0x71, 0x0a, 0x11, 0x43, 0x68, 0x65,
	0x63, 0x6b, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x12,
	0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f,
	0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2e, 0x0a, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x68, 0x65, 0x61,
	0x6c, 0x74, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x48, 0x65, 0x61, 0x6c,
	0x74, 0x68, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0xc6, 0x01, 0x0a,
	0x0f, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x44, 0x61, 0x74, 0x61,
	0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x41, 0x0a, 0x07, 0x64, 0x65, 0x74, 0x61,
	0x69, 0x6c, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x68, 0x65, 0x61, 0x6c,
	0x74, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x48, 0x65, 0x61, 0x6c, 0x74,
	0x68, 0x44, 0x61, 0x74, 0x61, 0x2e, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x52, 0x07, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x1a, 0x58, 0x0a, 0x0c, 0x44,
	0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x32, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x68,
	0x65, 0x61, 0x6c, 0x74, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x48, 0x65,
	0x61, 0x6c, 0x74, 0x68, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x74, 0x0a, 0x11, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x48, 0x65,
	0x61, 0x6c, 0x74, 0x68, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x65,
	0x64, 0x5f, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63, 0x68, 0x65, 0x63,
	0x6b, 0x65, 0x64, 0x41, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x32, 0x59, 0x0a, 0x0d, 0x48,
	0x65, 0x61, 0x6c, 0x74, 0x68, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x48, 0x0a, 0x0b,
	0x43, 0x68, 0x65, 0x63, 0x6b, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x12, 0x1b, 0x2e, 0x68, 0x65,
	0x61, 0x6c, 0x74, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x48, 0x65, 0x61,
	0x6c, 0x74, 0x68, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x1a, 0x1c, 0x2e, 0x68, 0x65, 0x61, 0x6c, 0x74,
	0x68, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x42, 0x0b, 0x5a, 0x09, 0x2e, 0x2f, 0x67, 0x72, 0x70, 0x63,
	0x61, 0x70, 0x70, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_grpcapp_health_proto_rawDescOnce sync.Once
	file_api_grpcapp_health_proto_rawDescData = file_api_grpcapp_health_proto_rawDesc
)

func file_api_grpcapp_health_proto_rawDescGZIP() []byte {
	file_api_grpcapp_health_proto_rawDescOnce.Do(func() {
		file_api_grpcapp_health_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_grpcapp_health_proto_rawDescData)
	})
	return file_api_grpcapp_health_proto_rawDescData
}

var file_api_grpcapp_health_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_api_grpcapp_health_proto_goTypes = []interface{}{
	(*CheckHealthParam)(nil),  // 0: health.v1.CheckHealthParam
	(*CheckHealthResult)(nil), // 1: health.v1.CheckHealthResult
	(*CheckHealthData)(nil),   // 2: health.v1.CheckHealthData
	(*CheckHealthDetail)(nil), // 3: health.v1.CheckHealthDetail
	nil,                       // 4: health.v1.CheckHealthData.DetailsEntry
}
var file_api_grpcapp_health_proto_depIdxs = []int32{
	2, // 0: health.v1.CheckHealthResult.data:type_name -> health.v1.CheckHealthData
	4, // 1: health.v1.CheckHealthData.details:type_name -> health.v1.CheckHealthData.DetailsEntry
	3, // 2: health.v1.CheckHealthData.DetailsEntry.value:type_name -> health.v1.CheckHealthDetail
	0, // 3: health.v1.HealthService.CheckHealth:input_type -> health.v1.CheckHealthParam
	1, // 4: health.v1.HealthService.CheckHealth:output_type -> health.v1.CheckHealthResult
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_api_grpcapp_health_proto_init() }
func file_api_grpcapp_health_proto_init() {
	if File_api_grpcapp_health_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_grpcapp_health_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckHealthParam); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_grpcapp_health_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckHealthResult); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_grpcapp_health_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckHealthData); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_grpcapp_health_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckHealthDetail); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_grpcapp_health_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_grpcapp_health_proto_goTypes,
		DependencyIndexes: file_api_grpcapp_health_proto_depIdxs,
		MessageInfos:      file_api_grpcapp_health_proto_msgTypes,
	}.Build()
	File_api_grpcapp_health_proto = out.File
	file_api_grpcapp_health_proto_rawDesc = nil
	file_api_grpcapp_health_proto_goTypes = nil
	file_api_grpcapp_health_proto_depIdxs = nil
}
