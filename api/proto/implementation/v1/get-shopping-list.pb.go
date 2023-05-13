// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.12
// source: v1/get-shopping-list.proto

package v1

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

type GetShoppingListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=userId,proto3" json:"userId,omitempty"`
}

func (x *GetShoppingListRequest) Reset() {
	*x = GetShoppingListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_get_shopping_list_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetShoppingListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetShoppingListRequest) ProtoMessage() {}

func (x *GetShoppingListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_get_shopping_list_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetShoppingListRequest.ProtoReflect.Descriptor instead.
func (*GetShoppingListRequest) Descriptor() ([]byte, []int) {
	return file_v1_get_shopping_list_proto_rawDescGZIP(), []int{0}
}

func (x *GetShoppingListRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type GetShoppingListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Purchases []*Purchase `protobuf:"bytes,1,rep,name=purchases,proto3" json:"purchases,omitempty"`
	Version   int32       `protobuf:"varint,2,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *GetShoppingListResponse) Reset() {
	*x = GetShoppingListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_get_shopping_list_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetShoppingListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetShoppingListResponse) ProtoMessage() {}

func (x *GetShoppingListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v1_get_shopping_list_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetShoppingListResponse.ProtoReflect.Descriptor instead.
func (*GetShoppingListResponse) Descriptor() ([]byte, []int) {
	return file_v1_get_shopping_list_proto_rawDescGZIP(), []int{1}
}

func (x *GetShoppingListResponse) GetPurchases() []*Purchase {
	if x != nil {
		return x.Purchases
	}
	return nil
}

func (x *GetShoppingListResponse) GetVersion() int32 {
	if x != nil {
		return x.Version
	}
	return 0
}

var File_v1_get_shopping_list_proto protoreflect.FileDescriptor

var file_v1_get_shopping_list_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x76, 0x31, 0x2f, 0x67, 0x65, 0x74, 0x2d, 0x73, 0x68, 0x6f, 0x70, 0x70, 0x69, 0x6e,
	0x67, 0x2d, 0x6c, 0x69, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x76, 0x31,
	0x1a, 0x12, 0x76, 0x31, 0x2f, 0x70, 0x75, 0x72, 0x63, 0x68, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x30, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x53, 0x68, 0x6f, 0x70, 0x70,
	0x69, 0x6e, 0x67, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16,
	0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x5f, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x53, 0x68, 0x6f,
	0x70, 0x70, 0x69, 0x6e, 0x67, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x2a, 0x0a, 0x09, 0x70, 0x75, 0x72, 0x63, 0x68, 0x61, 0x73, 0x65, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x75, 0x72, 0x63, 0x68, 0x61,
	0x73, 0x65, 0x52, 0x09, 0x70, 0x75, 0x72, 0x63, 0x68, 0x61, 0x73, 0x65, 0x73, 0x12, 0x18, 0x0a,
	0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x42, 0x44, 0x5a, 0x42, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x65, 0x70, 0x68, 0x69, 0x73, 0x74, 0x6f, 0x6c, 0x69,
	0x65, 0x2f, 0x63, 0x68, 0x65, 0x66, 0x62, 0x6f, 0x6f, 0x6b, 0x2d, 0x62, 0x61, 0x63, 0x6b, 0x65,
	0x6e, 0x64, 0x2d, 0x73, 0x68, 0x6f, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x2d, 0x6c, 0x69, 0x73, 0x74,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_v1_get_shopping_list_proto_rawDescOnce sync.Once
	file_v1_get_shopping_list_proto_rawDescData = file_v1_get_shopping_list_proto_rawDesc
)

func file_v1_get_shopping_list_proto_rawDescGZIP() []byte {
	file_v1_get_shopping_list_proto_rawDescOnce.Do(func() {
		file_v1_get_shopping_list_proto_rawDescData = protoimpl.X.CompressGZIP(file_v1_get_shopping_list_proto_rawDescData)
	})
	return file_v1_get_shopping_list_proto_rawDescData
}

var file_v1_get_shopping_list_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_v1_get_shopping_list_proto_goTypes = []interface{}{
	(*GetShoppingListRequest)(nil),  // 0: v1.GetShoppingListRequest
	(*GetShoppingListResponse)(nil), // 1: v1.GetShoppingListResponse
	(*Purchase)(nil),                // 2: v1.Purchase
}
var file_v1_get_shopping_list_proto_depIdxs = []int32{
	2, // 0: v1.GetShoppingListResponse.purchases:type_name -> v1.Purchase
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_v1_get_shopping_list_proto_init() }
func file_v1_get_shopping_list_proto_init() {
	if File_v1_get_shopping_list_proto != nil {
		return
	}
	file_v1_purchases_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_v1_get_shopping_list_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetShoppingListRequest); i {
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
		file_v1_get_shopping_list_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetShoppingListResponse); i {
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
			RawDescriptor: file_v1_get_shopping_list_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_v1_get_shopping_list_proto_goTypes,
		DependencyIndexes: file_v1_get_shopping_list_proto_depIdxs,
		MessageInfos:      file_v1_get_shopping_list_proto_msgTypes,
	}.Build()
	File_v1_get_shopping_list_proto = out.File
	file_v1_get_shopping_list_proto_rawDesc = nil
	file_v1_get_shopping_list_proto_goTypes = nil
	file_v1_get_shopping_list_proto_depIdxs = nil
}
