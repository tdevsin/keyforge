// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.0
// 	protoc        v5.29.2
// source: keyforge.proto

package proto

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

// Request format for getting a key
type GetKeyRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Key           string                 `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"` // The key for the operation
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetKeyRequest) Reset() {
	*x = GetKeyRequest{}
	mi := &file_keyforge_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetKeyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetKeyRequest) ProtoMessage() {}

func (x *GetKeyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_keyforge_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetKeyRequest.ProtoReflect.Descriptor instead.
func (*GetKeyRequest) Descriptor() ([]byte, []int) {
	return file_keyforge_proto_rawDescGZIP(), []int{0}
}

func (x *GetKeyRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

// Response format for getting a key
type GetKeyResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Key           string                 `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`     // The key for the operation
	Value         []byte                 `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"` // The value for the operation
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetKeyResponse) Reset() {
	*x = GetKeyResponse{}
	mi := &file_keyforge_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetKeyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetKeyResponse) ProtoMessage() {}

func (x *GetKeyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_keyforge_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetKeyResponse.ProtoReflect.Descriptor instead.
func (*GetKeyResponse) Descriptor() ([]byte, []int) {
	return file_keyforge_proto_rawDescGZIP(), []int{1}
}

func (x *GetKeyResponse) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *GetKeyResponse) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

// Request format for setting a key
type SetKeyRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Key           string                 `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`     // The key for the operation
	Value         []byte                 `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"` // The value for the operation
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SetKeyRequest) Reset() {
	*x = SetKeyRequest{}
	mi := &file_keyforge_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SetKeyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetKeyRequest) ProtoMessage() {}

func (x *SetKeyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_keyforge_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetKeyRequest.ProtoReflect.Descriptor instead.
func (*SetKeyRequest) Descriptor() ([]byte, []int) {
	return file_keyforge_proto_rawDescGZIP(), []int{2}
}

func (x *SetKeyRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *SetKeyRequest) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

// Response format for setting a key
type SetKeyResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Key           string                 `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`     // The key for the operation
	Value         []byte                 `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"` // The value for the operation
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SetKeyResponse) Reset() {
	*x = SetKeyResponse{}
	mi := &file_keyforge_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SetKeyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetKeyResponse) ProtoMessage() {}

func (x *SetKeyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_keyforge_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetKeyResponse.ProtoReflect.Descriptor instead.
func (*SetKeyResponse) Descriptor() ([]byte, []int) {
	return file_keyforge_proto_rawDescGZIP(), []int{3}
}

func (x *SetKeyResponse) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *SetKeyResponse) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

// Request format for deleting a key
type DeleteKeyRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Key           string                 `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"` // The key for the operation
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteKeyRequest) Reset() {
	*x = DeleteKeyRequest{}
	mi := &file_keyforge_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteKeyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteKeyRequest) ProtoMessage() {}

func (x *DeleteKeyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_keyforge_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteKeyRequest.ProtoReflect.Descriptor instead.
func (*DeleteKeyRequest) Descriptor() ([]byte, []int) {
	return file_keyforge_proto_rawDescGZIP(), []int{4}
}

func (x *DeleteKeyRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

// Response format for deleting a key
type DeleteKeyResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Key           string                 `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"` // The key for the operation
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteKeyResponse) Reset() {
	*x = DeleteKeyResponse{}
	mi := &file_keyforge_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteKeyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteKeyResponse) ProtoMessage() {}

func (x *DeleteKeyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_keyforge_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteKeyResponse.ProtoReflect.Descriptor instead.
func (*DeleteKeyResponse) Descriptor() ([]byte, []int) {
	return file_keyforge_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteKeyResponse) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

var File_keyforge_proto protoreflect.FileDescriptor

var file_keyforge_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x6b, 0x65, 0x79, 0x66, 0x6f, 0x72, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x21, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x4b, 0x65, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x22, 0x38, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x4b, 0x65, 0x79, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x37, 0x0a,
	0x0d, 0x53, 0x65, 0x74, 0x4b, 0x65, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x38, 0x0a, 0x0e, 0x53, 0x65, 0x74, 0x4b, 0x65, 0x79,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x22, 0x24, 0x0a, 0x10, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x22, 0x25, 0x0a, 0x11, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x4b, 0x65, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x32, 0x96, 0x01,
	0x0a, 0x0a, 0x4b, 0x65, 0x79, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x29, 0x0a, 0x06,
	0x47, 0x65, 0x74, 0x4b, 0x65, 0x79, 0x12, 0x0e, 0x2e, 0x47, 0x65, 0x74, 0x4b, 0x65, 0x79, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x47, 0x65, 0x74, 0x4b, 0x65, 0x79, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x29, 0x0a, 0x06, 0x53, 0x65, 0x74, 0x4b, 0x65,
	0x79, 0x12, 0x0e, 0x2e, 0x53, 0x65, 0x74, 0x4b, 0x65, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x0f, 0x2e, 0x53, 0x65, 0x74, 0x4b, 0x65, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x32, 0x0a, 0x09, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x12,
	0x11, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x12, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x23, 0x5a, 0x21, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x64, 0x65, 0x76, 0x73, 0x69, 0x6e, 0x2f, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_keyforge_proto_rawDescOnce sync.Once
	file_keyforge_proto_rawDescData = file_keyforge_proto_rawDesc
)

func file_keyforge_proto_rawDescGZIP() []byte {
	file_keyforge_proto_rawDescOnce.Do(func() {
		file_keyforge_proto_rawDescData = protoimpl.X.CompressGZIP(file_keyforge_proto_rawDescData)
	})
	return file_keyforge_proto_rawDescData
}

var file_keyforge_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_keyforge_proto_goTypes = []any{
	(*GetKeyRequest)(nil),     // 0: GetKeyRequest
	(*GetKeyResponse)(nil),    // 1: GetKeyResponse
	(*SetKeyRequest)(nil),     // 2: SetKeyRequest
	(*SetKeyResponse)(nil),    // 3: SetKeyResponse
	(*DeleteKeyRequest)(nil),  // 4: DeleteKeyRequest
	(*DeleteKeyResponse)(nil), // 5: DeleteKeyResponse
}
var file_keyforge_proto_depIdxs = []int32{
	0, // 0: KeyService.GetKey:input_type -> GetKeyRequest
	2, // 1: KeyService.SetKey:input_type -> SetKeyRequest
	4, // 2: KeyService.DeleteKey:input_type -> DeleteKeyRequest
	1, // 3: KeyService.GetKey:output_type -> GetKeyResponse
	3, // 4: KeyService.SetKey:output_type -> SetKeyResponse
	5, // 5: KeyService.DeleteKey:output_type -> DeleteKeyResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_keyforge_proto_init() }
func file_keyforge_proto_init() {
	if File_keyforge_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_keyforge_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_keyforge_proto_goTypes,
		DependencyIndexes: file_keyforge_proto_depIdxs,
		MessageInfos:      file_keyforge_proto_msgTypes,
	}.Build()
	File_keyforge_proto = out.File
	file_keyforge_proto_rawDesc = nil
	file_keyforge_proto_goTypes = nil
	file_keyforge_proto_depIdxs = nil
}
