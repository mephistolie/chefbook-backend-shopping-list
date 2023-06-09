// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.12
// source: v1/delete-user-from-shopping-list.proto

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

type DeleteUserFromShoppingListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RequesterId    string `protobuf:"bytes,1,opt,name=requesterId,proto3" json:"requesterId,omitempty"`
	ShoppingListId string `protobuf:"bytes,2,opt,name=shoppingListId,proto3" json:"shoppingListId,omitempty"`
	UserId         string `protobuf:"bytes,3,opt,name=userId,proto3" json:"userId,omitempty"`
}

func (x *DeleteUserFromShoppingListRequest) Reset() {
	*x = DeleteUserFromShoppingListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_delete_user_from_shopping_list_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteUserFromShoppingListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteUserFromShoppingListRequest) ProtoMessage() {}

func (x *DeleteUserFromShoppingListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_delete_user_from_shopping_list_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteUserFromShoppingListRequest.ProtoReflect.Descriptor instead.
func (*DeleteUserFromShoppingListRequest) Descriptor() ([]byte, []int) {
	return file_v1_delete_user_from_shopping_list_proto_rawDescGZIP(), []int{0}
}

func (x *DeleteUserFromShoppingListRequest) GetRequesterId() string {
	if x != nil {
		return x.RequesterId
	}
	return ""
}

func (x *DeleteUserFromShoppingListRequest) GetShoppingListId() string {
	if x != nil {
		return x.ShoppingListId
	}
	return ""
}

func (x *DeleteUserFromShoppingListRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type DeleteUserFromShoppingListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *DeleteUserFromShoppingListResponse) Reset() {
	*x = DeleteUserFromShoppingListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_delete_user_from_shopping_list_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteUserFromShoppingListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteUserFromShoppingListResponse) ProtoMessage() {}

func (x *DeleteUserFromShoppingListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v1_delete_user_from_shopping_list_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteUserFromShoppingListResponse.ProtoReflect.Descriptor instead.
func (*DeleteUserFromShoppingListResponse) Descriptor() ([]byte, []int) {
	return file_v1_delete_user_from_shopping_list_proto_rawDescGZIP(), []int{1}
}

func (x *DeleteUserFromShoppingListResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_v1_delete_user_from_shopping_list_proto protoreflect.FileDescriptor

var file_v1_delete_user_from_shopping_list_proto_rawDesc = []byte{
	0x0a, 0x27, 0x76, 0x31, 0x2f, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x2d, 0x75, 0x73, 0x65, 0x72,
	0x2d, 0x66, 0x72, 0x6f, 0x6d, 0x2d, 0x73, 0x68, 0x6f, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x2d, 0x6c,
	0x69, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x76, 0x31, 0x22, 0x85, 0x01,
	0x0a, 0x21, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x46, 0x72, 0x6f, 0x6d,
	0x53, 0x68, 0x6f, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x65, 0x72,
	0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x65, 0x72, 0x49, 0x64, 0x12, 0x26, 0x0a, 0x0e, 0x73, 0x68, 0x6f, 0x70, 0x70, 0x69, 0x6e,
	0x67, 0x4c, 0x69, 0x73, 0x74, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x73,
	0x68, 0x6f, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x4c, 0x69, 0x73, 0x74, 0x49, 0x64, 0x12, 0x16, 0x0a,
	0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x3e, 0x0a, 0x22, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x55,
	0x73, 0x65, 0x72, 0x46, 0x72, 0x6f, 0x6d, 0x53, 0x68, 0x6f, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x4c,
	0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x42, 0x44, 0x5a, 0x42, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x65, 0x70, 0x68, 0x69, 0x73, 0x74, 0x6f, 0x6c, 0x69, 0x65, 0x2f,
	0x63, 0x68, 0x65, 0x66, 0x62, 0x6f, 0x6f, 0x6b, 0x2d, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64,
	0x2d, 0x73, 0x68, 0x6f, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x2d, 0x6c, 0x69, 0x73, 0x74, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_v1_delete_user_from_shopping_list_proto_rawDescOnce sync.Once
	file_v1_delete_user_from_shopping_list_proto_rawDescData = file_v1_delete_user_from_shopping_list_proto_rawDesc
)

func file_v1_delete_user_from_shopping_list_proto_rawDescGZIP() []byte {
	file_v1_delete_user_from_shopping_list_proto_rawDescOnce.Do(func() {
		file_v1_delete_user_from_shopping_list_proto_rawDescData = protoimpl.X.CompressGZIP(file_v1_delete_user_from_shopping_list_proto_rawDescData)
	})
	return file_v1_delete_user_from_shopping_list_proto_rawDescData
}

var file_v1_delete_user_from_shopping_list_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_v1_delete_user_from_shopping_list_proto_goTypes = []interface{}{
	(*DeleteUserFromShoppingListRequest)(nil),  // 0: v1.DeleteUserFromShoppingListRequest
	(*DeleteUserFromShoppingListResponse)(nil), // 1: v1.DeleteUserFromShoppingListResponse
}
var file_v1_delete_user_from_shopping_list_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_v1_delete_user_from_shopping_list_proto_init() }
func file_v1_delete_user_from_shopping_list_proto_init() {
	if File_v1_delete_user_from_shopping_list_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_v1_delete_user_from_shopping_list_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteUserFromShoppingListRequest); i {
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
		file_v1_delete_user_from_shopping_list_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteUserFromShoppingListResponse); i {
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
			RawDescriptor: file_v1_delete_user_from_shopping_list_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_v1_delete_user_from_shopping_list_proto_goTypes,
		DependencyIndexes: file_v1_delete_user_from_shopping_list_proto_depIdxs,
		MessageInfos:      file_v1_delete_user_from_shopping_list_proto_msgTypes,
	}.Build()
	File_v1_delete_user_from_shopping_list_proto = out.File
	file_v1_delete_user_from_shopping_list_proto_rawDesc = nil
	file_v1_delete_user_from_shopping_list_proto_goTypes = nil
	file_v1_delete_user_from_shopping_list_proto_depIdxs = nil
}
