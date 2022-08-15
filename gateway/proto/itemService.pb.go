// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.4
// source: proto/itemService.proto

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

type DeleteFavReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserID int64 `protobuf:"varint,1,opt,name=userID,proto3" json:"userID,omitempty"`
	ItemID int64 `protobuf:"varint,2,opt,name=itemID,proto3" json:"itemID,omitempty"`
	ShopID int64 `protobuf:"varint,3,opt,name=shopID,proto3" json:"shopID,omitempty"`
}

func (x *DeleteFavReq) Reset() {
	*x = DeleteFavReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_itemService_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteFavReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteFavReq) ProtoMessage() {}

func (x *DeleteFavReq) ProtoReflect() protoreflect.Message {
	mi := &file_proto_itemService_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteFavReq.ProtoReflect.Descriptor instead.
func (*DeleteFavReq) Descriptor() ([]byte, []int) {
	return file_proto_itemService_proto_rawDescGZIP(), []int{0}
}

func (x *DeleteFavReq) GetUserID() int64 {
	if x != nil {
		return x.UserID
	}
	return 0
}

func (x *DeleteFavReq) GetItemID() int64 {
	if x != nil {
		return x.ItemID
	}
	return 0
}

func (x *DeleteFavReq) GetShopID() int64 {
	if x != nil {
		return x.ShopID
	}
	return 0
}

type DeleteFavRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ErrorCode int32  `protobuf:"varint,1,opt,name=errorCode,proto3" json:"errorCode,omitempty"`
	ErrorMsg  string `protobuf:"bytes,2,opt,name=errorMsg,proto3" json:"errorMsg,omitempty"`
}

func (x *DeleteFavRes) Reset() {
	*x = DeleteFavRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_itemService_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteFavRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteFavRes) ProtoMessage() {}

func (x *DeleteFavRes) ProtoReflect() protoreflect.Message {
	mi := &file_proto_itemService_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteFavRes.ProtoReflect.Descriptor instead.
func (*DeleteFavRes) Descriptor() ([]byte, []int) {
	return file_proto_itemService_proto_rawDescGZIP(), []int{1}
}

func (x *DeleteFavRes) GetErrorCode() int32 {
	if x != nil {
		return x.ErrorCode
	}
	return 0
}

func (x *DeleteFavRes) GetErrorMsg() string {
	if x != nil {
		return x.ErrorMsg
	}
	return ""
}

type AddFavReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserID int64 `protobuf:"varint,1,opt,name=userID,proto3" json:"userID,omitempty"`
	ItemID int64 `protobuf:"varint,2,opt,name=itemID,proto3" json:"itemID,omitempty"`
	ShopID int64 `protobuf:"varint,3,opt,name=shopID,proto3" json:"shopID,omitempty"`
}

func (x *AddFavReq) Reset() {
	*x = AddFavReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_itemService_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddFavReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddFavReq) ProtoMessage() {}

func (x *AddFavReq) ProtoReflect() protoreflect.Message {
	mi := &file_proto_itemService_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddFavReq.ProtoReflect.Descriptor instead.
func (*AddFavReq) Descriptor() ([]byte, []int) {
	return file_proto_itemService_proto_rawDescGZIP(), []int{2}
}

func (x *AddFavReq) GetUserID() int64 {
	if x != nil {
		return x.UserID
	}
	return 0
}

func (x *AddFavReq) GetItemID() int64 {
	if x != nil {
		return x.ItemID
	}
	return 0
}

func (x *AddFavReq) GetShopID() int64 {
	if x != nil {
		return x.ShopID
	}
	return 0
}

type AddFavRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ErrorCode int32  `protobuf:"varint,1,opt,name=errorCode,proto3" json:"errorCode,omitempty"`
	ErrorMsg  string `protobuf:"bytes,2,opt,name=errorMsg,proto3" json:"errorMsg,omitempty"`
	Item      *Item  `protobuf:"bytes,3,opt,name=item,proto3" json:"item,omitempty"`
}

func (x *AddFavRes) Reset() {
	*x = AddFavRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_itemService_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddFavRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddFavRes) ProtoMessage() {}

func (x *AddFavRes) ProtoReflect() protoreflect.Message {
	mi := &file_proto_itemService_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddFavRes.ProtoReflect.Descriptor instead.
func (*AddFavRes) Descriptor() ([]byte, []int) {
	return file_proto_itemService_proto_rawDescGZIP(), []int{3}
}

func (x *AddFavRes) GetErrorCode() int32 {
	if x != nil {
		return x.ErrorCode
	}
	return 0
}

func (x *AddFavRes) GetErrorMsg() string {
	if x != nil {
		return x.ErrorMsg
	}
	return ""
}

func (x *AddFavRes) GetItem() *Item {
	if x != nil {
		return x.Item
	}
	return nil
}

type Item struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name   string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Price  int64  `protobuf:"varint,2,opt,name=price,proto3" json:"price,omitempty"`
	ShopID int64  `protobuf:"varint,3,opt,name=shopID,proto3" json:"shopID,omitempty"`
	ItemID int64  `protobuf:"varint,4,opt,name=itemID,proto3" json:"itemID,omitempty"`
}

func (x *Item) Reset() {
	*x = Item{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_itemService_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Item) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Item) ProtoMessage() {}

func (x *Item) ProtoReflect() protoreflect.Message {
	mi := &file_proto_itemService_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Item.ProtoReflect.Descriptor instead.
func (*Item) Descriptor() ([]byte, []int) {
	return file_proto_itemService_proto_rawDescGZIP(), []int{4}
}

func (x *Item) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Item) GetPrice() int64 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *Item) GetShopID() int64 {
	if x != nil {
		return x.ShopID
	}
	return 0
}

func (x *Item) GetItemID() int64 {
	if x != nil {
		return x.ItemID
	}
	return 0
}

type GetFavListReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserID int64 `protobuf:"varint,1,opt,name=userID,proto3" json:"userID,omitempty"`
	Page   int32 `protobuf:"varint,2,opt,name=page,proto3" json:"page,omitempty"`
}

func (x *GetFavListReq) Reset() {
	*x = GetFavListReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_itemService_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFavListReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFavListReq) ProtoMessage() {}

func (x *GetFavListReq) ProtoReflect() protoreflect.Message {
	mi := &file_proto_itemService_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFavListReq.ProtoReflect.Descriptor instead.
func (*GetFavListReq) Descriptor() ([]byte, []int) {
	return file_proto_itemService_proto_rawDescGZIP(), []int{5}
}

func (x *GetFavListReq) GetUserID() int64 {
	if x != nil {
		return x.UserID
	}
	return 0
}

func (x *GetFavListReq) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

type GetFavListRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ErrorCode  int32   `protobuf:"varint,1,opt,name=errorCode,proto3" json:"errorCode,omitempty"`
	ErrorMsg   string  `protobuf:"bytes,2,opt,name=errorMsg,proto3" json:"errorMsg,omitempty"`
	Items      []*Item `protobuf:"bytes,3,rep,name=items,proto3" json:"items,omitempty"`
	TotalPages int32   `protobuf:"varint,4,opt,name=totalPages,proto3" json:"totalPages,omitempty"`
}

func (x *GetFavListRes) Reset() {
	*x = GetFavListRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_itemService_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFavListRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFavListRes) ProtoMessage() {}

func (x *GetFavListRes) ProtoReflect() protoreflect.Message {
	mi := &file_proto_itemService_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFavListRes.ProtoReflect.Descriptor instead.
func (*GetFavListRes) Descriptor() ([]byte, []int) {
	return file_proto_itemService_proto_rawDescGZIP(), []int{6}
}

func (x *GetFavListRes) GetErrorCode() int32 {
	if x != nil {
		return x.ErrorCode
	}
	return 0
}

func (x *GetFavListRes) GetErrorMsg() string {
	if x != nil {
		return x.ErrorMsg
	}
	return ""
}

func (x *GetFavListRes) GetItems() []*Item {
	if x != nil {
		return x.Items
	}
	return nil
}

func (x *GetFavListRes) GetTotalPages() int32 {
	if x != nil {
		return x.TotalPages
	}
	return 0
}

var File_proto_itemService_proto protoreflect.FileDescriptor

var file_proto_itemService_proto_rawDesc = []byte{
	0x0a, 0x17, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x69, 0x74, 0x65, 0x6d, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x56, 0x0a, 0x0c, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x61, 0x76, 0x52, 0x65, 0x71,
	0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x44, 0x12, 0x16, 0x0a, 0x06, 0x69, 0x74, 0x65, 0x6d,
	0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x69, 0x74, 0x65, 0x6d, 0x49, 0x44,
	0x12, 0x16, 0x0a, 0x06, 0x73, 0x68, 0x6f, 0x70, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x06, 0x73, 0x68, 0x6f, 0x70, 0x49, 0x44, 0x22, 0x48, 0x0a, 0x0c, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x46, 0x61, 0x76, 0x52, 0x65, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x65, 0x72, 0x72,
	0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d,
	0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d,
	0x73, 0x67, 0x22, 0x53, 0x0a, 0x09, 0x41, 0x64, 0x64, 0x46, 0x61, 0x76, 0x52, 0x65, 0x71, 0x12,
	0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x44, 0x12, 0x16, 0x0a, 0x06, 0x69, 0x74, 0x65, 0x6d, 0x49,
	0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x69, 0x74, 0x65, 0x6d, 0x49, 0x44, 0x12,
	0x16, 0x0a, 0x06, 0x73, 0x68, 0x6f, 0x70, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x06, 0x73, 0x68, 0x6f, 0x70, 0x49, 0x44, 0x22, 0x66, 0x0a, 0x09, 0x41, 0x64, 0x64, 0x46, 0x61,
	0x76, 0x52, 0x65, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f,
	0x64, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x73, 0x67, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x73, 0x67, 0x12, 0x1f,
	0x0a, 0x04, 0x69, 0x74, 0x65, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x04, 0x69, 0x74, 0x65, 0x6d, 0x22,
	0x60, 0x0a, 0x04, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x70,
	0x72, 0x69, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x70, 0x72, 0x69, 0x63,
	0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x68, 0x6f, 0x70, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x06, 0x73, 0x68, 0x6f, 0x70, 0x49, 0x44, 0x12, 0x16, 0x0a, 0x06, 0x69, 0x74, 0x65,
	0x6d, 0x49, 0x44, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x69, 0x74, 0x65, 0x6d, 0x49,
	0x44, 0x22, 0x3b, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x46, 0x61, 0x76, 0x4c, 0x69, 0x73, 0x74, 0x52,
	0x65, 0x71, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x44, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61,
	0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x22, 0x8c,
	0x01, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x46, 0x61, 0x76, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73,
	0x12, 0x1c, 0x0a, 0x09, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x09, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x73, 0x67, 0x12, 0x21, 0x0a, 0x05, 0x69, 0x74,
	0x65, 0x6d, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x12, 0x1e, 0x0a,
	0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x50, 0x61, 0x67, 0x65, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x50, 0x61, 0x67, 0x65, 0x73, 0x32, 0xb2, 0x01,
	0x0a, 0x0b, 0x49, 0x74, 0x65, 0x6d, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x37, 0x0a,
	0x09, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x61, 0x76, 0x12, 0x13, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x61, 0x76, 0x52, 0x65, 0x71, 0x1a,
	0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x61,
	0x76, 0x52, 0x65, 0x73, 0x22, 0x00, 0x12, 0x2e, 0x0a, 0x06, 0x41, 0x64, 0x64, 0x46, 0x61, 0x76,
	0x12, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x64, 0x64, 0x46, 0x61, 0x76, 0x52,
	0x65, 0x71, 0x1a, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x64, 0x64, 0x46, 0x61,
	0x76, 0x52, 0x65, 0x73, 0x22, 0x00, 0x12, 0x3a, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x46, 0x61, 0x76,
	0x4c, 0x69, 0x73, 0x74, 0x12, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74,
	0x46, 0x61, 0x76, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x1a, 0x14, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x46, 0x61, 0x76, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73,
	0x22, 0x00, 0x42, 0x13, 0x5a, 0x11, 0x69, 0x74, 0x65, 0x6d, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_itemService_proto_rawDescOnce sync.Once
	file_proto_itemService_proto_rawDescData = file_proto_itemService_proto_rawDesc
)

func file_proto_itemService_proto_rawDescGZIP() []byte {
	file_proto_itemService_proto_rawDescOnce.Do(func() {
		file_proto_itemService_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_itemService_proto_rawDescData)
	})
	return file_proto_itemService_proto_rawDescData
}

var file_proto_itemService_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_proto_itemService_proto_goTypes = []interface{}{
	(*DeleteFavReq)(nil),  // 0: proto.DeleteFavReq
	(*DeleteFavRes)(nil),  // 1: proto.DeleteFavRes
	(*AddFavReq)(nil),     // 2: proto.AddFavReq
	(*AddFavRes)(nil),     // 3: proto.AddFavRes
	(*Item)(nil),          // 4: proto.Item
	(*GetFavListReq)(nil), // 5: proto.GetFavListReq
	(*GetFavListRes)(nil), // 6: proto.GetFavListRes
}
var file_proto_itemService_proto_depIdxs = []int32{
	4, // 0: proto.AddFavRes.item:type_name -> proto.Item
	4, // 1: proto.GetFavListRes.items:type_name -> proto.Item
	0, // 2: proto.ItemService.DeleteFav:input_type -> proto.DeleteFavReq
	2, // 3: proto.ItemService.AddFav:input_type -> proto.AddFavReq
	5, // 4: proto.ItemService.GetFavList:input_type -> proto.GetFavListReq
	1, // 5: proto.ItemService.DeleteFav:output_type -> proto.DeleteFavRes
	3, // 6: proto.ItemService.AddFav:output_type -> proto.AddFavRes
	6, // 7: proto.ItemService.GetFavList:output_type -> proto.GetFavListRes
	5, // [5:8] is the sub-list for method output_type
	2, // [2:5] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_itemService_proto_init() }
func file_proto_itemService_proto_init() {
	if File_proto_itemService_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_itemService_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteFavReq); i {
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
		file_proto_itemService_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteFavRes); i {
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
		file_proto_itemService_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddFavReq); i {
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
		file_proto_itemService_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddFavRes); i {
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
		file_proto_itemService_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Item); i {
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
		file_proto_itemService_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFavListReq); i {
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
		file_proto_itemService_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFavListRes); i {
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
			RawDescriptor: file_proto_itemService_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_itemService_proto_goTypes,
		DependencyIndexes: file_proto_itemService_proto_depIdxs,
		MessageInfos:      file_proto_itemService_proto_msgTypes,
	}.Build()
	File_proto_itemService_proto = out.File
	file_proto_itemService_proto_rawDesc = nil
	file_proto_itemService_proto_goTypes = nil
	file_proto_itemService_proto_depIdxs = nil
}
