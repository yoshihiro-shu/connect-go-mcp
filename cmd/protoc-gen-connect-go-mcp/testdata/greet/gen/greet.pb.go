//*
// Greet related messages.
//
// This file is really just an example. The data model is completely
// fictional.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: greet.proto

package greetv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Greeting request
type GreetRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"` // Name
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GreetRequest) Reset() {
	*x = GreetRequest{}
	mi := &file_greet_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GreetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GreetRequest) ProtoMessage() {}

func (x *GreetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_greet_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GreetRequest.ProtoReflect.Descriptor instead.
func (*GreetRequest) Descriptor() ([]byte, []int) {
	return file_greet_proto_rawDescGZIP(), []int{0}
}

func (x *GreetRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type GreetResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Greeting      string                 `protobuf:"bytes,1,opt,name=greeting,proto3" json:"greeting,omitempty"` // Greeting
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GreetResponse) Reset() {
	*x = GreetResponse{}
	mi := &file_greet_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GreetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GreetResponse) ProtoMessage() {}

func (x *GreetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_greet_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GreetResponse.ProtoReflect.Descriptor instead.
func (*GreetResponse) Descriptor() ([]byte, []int) {
	return file_greet_proto_rawDescGZIP(), []int{1}
}

func (x *GreetResponse) GetGreeting() string {
	if x != nil {
		return x.Greeting
	}
	return ""
}

// Ping request
type PingRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Message       string                 `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"` // Message
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PingRequest) Reset() {
	*x = PingRequest{}
	mi := &file_greet_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingRequest) ProtoMessage() {}

func (x *PingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_greet_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingRequest.ProtoReflect.Descriptor instead.
func (*PingRequest) Descriptor() ([]byte, []int) {
	return file_greet_proto_rawDescGZIP(), []int{2}
}

func (x *PingRequest) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type PingResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Message       string                 `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"` // Message
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PingResponse) Reset() {
	*x = PingResponse{}
	mi := &file_greet_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingResponse) ProtoMessage() {}

func (x *PingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_greet_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingResponse.ProtoReflect.Descriptor instead.
func (*PingResponse) Descriptor() ([]byte, []int) {
	return file_greet_proto_rawDescGZIP(), []int{3}
}

func (x *PingResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_greet_proto protoreflect.FileDescriptor

const file_greet_proto_rawDesc = "" +
	"\n" +
	"\vgreet.proto\x12\bgreet.v1\"\"\n" +
	"\fGreetRequest\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\"+\n" +
	"\rGreetResponse\x12\x1a\n" +
	"\bgreeting\x18\x01 \x01(\tR\bgreeting\"'\n" +
	"\vPingRequest\x12\x18\n" +
	"\amessage\x18\x01 \x01(\tR\amessage\"(\n" +
	"\fPingResponse\x12\x18\n" +
	"\amessage\x18\x01 \x01(\tR\amessage2\x83\x01\n" +
	"\fGreetService\x12:\n" +
	"\x05Greet\x12\x16.greet.v1.GreetRequest\x1a\x17.greet.v1.GreetResponse\"\x00\x127\n" +
	"\x04Ping\x12\x15.greet.v1.PingRequest\x1a\x16.greet.v1.PingResponse\"\x00BGZEgithub.com/yoshihiro-shu/connect-go-mcp/internal/gen/greet/v1;greetv1b\x06proto3"

var (
	file_greet_proto_rawDescOnce sync.Once
	file_greet_proto_rawDescData []byte
)

func file_greet_proto_rawDescGZIP() []byte {
	file_greet_proto_rawDescOnce.Do(func() {
		file_greet_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_greet_proto_rawDesc), len(file_greet_proto_rawDesc)))
	})
	return file_greet_proto_rawDescData
}

var file_greet_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_greet_proto_goTypes = []any{
	(*GreetRequest)(nil),  // 0: greet.v1.GreetRequest
	(*GreetResponse)(nil), // 1: greet.v1.GreetResponse
	(*PingRequest)(nil),   // 2: greet.v1.PingRequest
	(*PingResponse)(nil),  // 3: greet.v1.PingResponse
}
var file_greet_proto_depIdxs = []int32{
	0, // 0: greet.v1.GreetService.Greet:input_type -> greet.v1.GreetRequest
	2, // 1: greet.v1.GreetService.Ping:input_type -> greet.v1.PingRequest
	1, // 2: greet.v1.GreetService.Greet:output_type -> greet.v1.GreetResponse
	3, // 3: greet.v1.GreetService.Ping:output_type -> greet.v1.PingResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_greet_proto_init() }
func file_greet_proto_init() {
	if File_greet_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_greet_proto_rawDesc), len(file_greet_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_greet_proto_goTypes,
		DependencyIndexes: file_greet_proto_depIdxs,
		MessageInfos:      file_greet_proto_msgTypes,
	}.Build()
	File_greet_proto = out.File
	file_greet_proto_goTypes = nil
	file_greet_proto_depIdxs = nil
}
