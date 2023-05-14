// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.12
// source: v1/get-shopping-lists.proto

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

type GetShoppingListsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=userId,proto3" json:"userId,omitempty"`
}

func (x *GetShoppingListsRequest) Reset() {
	*x = GetShoppingListsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_get_shopping_lists_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetShoppingListsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetShoppingListsRequest) ProtoMessage() {}

func (x *GetShoppingListsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_get_shopping_lists_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetShoppingListsRequest.ProtoReflect.Descriptor instead.
func (*GetShoppingListsRequest) Descriptor() ([]byte, []int) {
	return file_v1_get_shopping_lists_proto_rawDescGZIP(), []int{0}
}

func (x *GetShoppingListsRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type GetShoppingListsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ShoppingLists []*ShoppingListInfo `protobuf:"bytes,1,rep,name=shoppingLists,proto3" json:"shoppingLists,omitempty"`
}

func (x *GetShoppingListsResponse) Reset() {
	*x = GetShoppingListsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_get_shopping_lists_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetShoppingListsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetShoppingListsResponse) ProtoMessage() {}

func (x *GetShoppingListsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v1_get_shopping_lists_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetShoppingListsResponse.ProtoReflect.Descriptor instead.
func (*GetShoppingListsResponse) Descriptor() ([]byte, []int) {
	return file_v1_get_shopping_lists_proto_rawDescGZIP(), []int{1}
}

func (x *GetShoppingListsResponse) GetShoppingLists() []*ShoppingListInfo {
	if x != nil {
		return x.ShoppingLists
	}
	return nil
}

var File_v1_get_shopping_lists_proto protoreflect.FileDescriptor

var file_v1_get_shopping_lists_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x76, 0x31, 0x2f, 0x67, 0x65, 0x74, 0x2d, 0x73, 0x68, 0x6f, 0x70, 0x70, 0x69, 0x6e,
	0x67, 0x2d, 0x6c, 0x69, 0x73, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x76,
	0x31, 0x1a, 0x16, 0x76, 0x31, 0x2f, 0x73, 0x68, 0x6f, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x2d, 0x6c,
	0x69, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x31, 0x0a, 0x17, 0x47, 0x65, 0x74,
	0x53, 0x68, 0x6f, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x4c, 0x69, 0x73, 0x74, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x56, 0x0a, 0x18,
	0x47, 0x65, 0x74, 0x53, 0x68, 0x6f, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x4c, 0x69, 0x73, 0x74, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3a, 0x0a, 0x0d, 0x73, 0x68, 0x6f, 0x70,
	0x70, 0x69, 0x6e, 0x67, 0x4c, 0x69, 0x73, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x14, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x68, 0x6f, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x4c, 0x69, 0x73,
	0x74, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x0d, 0x73, 0x68, 0x6f, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x4c,
	0x69, 0x73, 0x74, 0x73, 0x42, 0x44, 0x5a, 0x42, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x6d, 0x65, 0x70, 0x68, 0x69, 0x73, 0x74, 0x6f, 0x6c, 0x69, 0x65, 0x2f, 0x63,
	0x68, 0x65, 0x66, 0x62, 0x6f, 0x6f, 0x6b, 0x2d, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2d,
	0x73, 0x68, 0x6f, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x2d, 0x6c, 0x69, 0x73, 0x74, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_v1_get_shopping_lists_proto_rawDescOnce sync.Once
	file_v1_get_shopping_lists_proto_rawDescData = file_v1_get_shopping_lists_proto_rawDesc
)

func file_v1_get_shopping_lists_proto_rawDescGZIP() []byte {
	file_v1_get_shopping_lists_proto_rawDescOnce.Do(func() {
		file_v1_get_shopping_lists_proto_rawDescData = protoimpl.X.CompressGZIP(file_v1_get_shopping_lists_proto_rawDescData)
	})
	return file_v1_get_shopping_lists_proto_rawDescData
}

var file_v1_get_shopping_lists_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_v1_get_shopping_lists_proto_goTypes = []interface{}{
	(*GetShoppingListsRequest)(nil),  // 0: v1.GetShoppingListsRequest
	(*GetShoppingListsResponse)(nil), // 1: v1.GetShoppingListsResponse
	(*ShoppingListInfo)(nil),         // 2: v1.ShoppingListInfo
}
var file_v1_get_shopping_lists_proto_depIdxs = []int32{
	2, // 0: v1.GetShoppingListsResponse.shoppingLists:type_name -> v1.ShoppingListInfo
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_v1_get_shopping_lists_proto_init() }
func file_v1_get_shopping_lists_proto_init() {
	if File_v1_get_shopping_lists_proto != nil {
		return
	}
	file_v1_shopping_list_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_v1_get_shopping_lists_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetShoppingListsRequest); i {
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
		file_v1_get_shopping_lists_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetShoppingListsResponse); i {
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
			RawDescriptor: file_v1_get_shopping_lists_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_v1_get_shopping_lists_proto_goTypes,
		DependencyIndexes: file_v1_get_shopping_lists_proto_depIdxs,
		MessageInfos:      file_v1_get_shopping_lists_proto_msgTypes,
	}.Build()
	File_v1_get_shopping_lists_proto = out.File
	file_v1_get_shopping_lists_proto_rawDesc = nil
	file_v1_get_shopping_lists_proto_goTypes = nil
	file_v1_get_shopping_lists_proto_depIdxs = nil
}
