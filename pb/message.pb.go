// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v5.26.1
// source: pb/message.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProtoVersion string     `protobuf:"bytes,1,opt,name=ProtoVersion,proto3" json:"ProtoVersion,omitempty"`
	Content      *anypb.Any `protobuf:"bytes,2,opt,name=Content,proto3" json:"Content,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_pb_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_pb_message_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetProtoVersion() string {
	if x != nil {
		return x.ProtoVersion
	}
	return ""
}

func (x *Message) GetContent() *anypb.Any {
	if x != nil {
		return x.Content
	}
	return nil
}

type Subscribe struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Queue string `protobuf:"bytes,1,opt,name=Queue,proto3" json:"Queue,omitempty"`
}

func (x *Subscribe) Reset() {
	*x = Subscribe{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Subscribe) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Subscribe) ProtoMessage() {}

func (x *Subscribe) ProtoReflect() protoreflect.Message {
	mi := &file_pb_message_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Subscribe.ProtoReflect.Descriptor instead.
func (*Subscribe) Descriptor() ([]byte, []int) {
	return file_pb_message_proto_rawDescGZIP(), []int{1}
}

func (x *Subscribe) GetQueue() string {
	if x != nil {
		return x.Queue
	}
	return ""
}

type Unsubscribe struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Queue string `protobuf:"bytes,1,opt,name=Queue,proto3" json:"Queue,omitempty"`
}

func (x *Unsubscribe) Reset() {
	*x = Unsubscribe{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_message_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Unsubscribe) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Unsubscribe) ProtoMessage() {}

func (x *Unsubscribe) ProtoReflect() protoreflect.Message {
	mi := &file_pb_message_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Unsubscribe.ProtoReflect.Descriptor instead.
func (*Unsubscribe) Descriptor() ([]byte, []int) {
	return file_pb_message_proto_rawDescGZIP(), []int{2}
}

func (x *Unsubscribe) GetQueue() string {
	if x != nil {
		return x.Queue
	}
	return ""
}

var File_pb_message_proto protoreflect.FileDescriptor

var file_pb_message_proto_rawDesc = []byte{
	0x0a, 0x10, 0x70, 0x62, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x5d, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x22, 0x0a, 0x0c,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0c, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x12, 0x2e, 0x0a, 0x07, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x07, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x22, 0x21, 0x0a, 0x09, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x51, 0x75, 0x65, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x51, 0x75,
	0x65, 0x75, 0x65, 0x22, 0x23, 0x0a, 0x0b, 0x55, 0x6e, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69,
	0x62, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x51, 0x75, 0x65, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x51, 0x75, 0x65, 0x75, 0x65, 0x42, 0x06, 0x5a, 0x04, 0x2e, 0x2f, 0x70, 0x62,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pb_message_proto_rawDescOnce sync.Once
	file_pb_message_proto_rawDescData = file_pb_message_proto_rawDesc
)

func file_pb_message_proto_rawDescGZIP() []byte {
	file_pb_message_proto_rawDescOnce.Do(func() {
		file_pb_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_pb_message_proto_rawDescData)
	})
	return file_pb_message_proto_rawDescData
}

var file_pb_message_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_pb_message_proto_goTypes = []interface{}{
	(*Message)(nil),     // 0: pb.Message
	(*Subscribe)(nil),   // 1: pb.Subscribe
	(*Unsubscribe)(nil), // 2: pb.Unsubscribe
	(*anypb.Any)(nil),   // 3: google.protobuf.Any
}
var file_pb_message_proto_depIdxs = []int32{
	3, // 0: pb.Message.Content:type_name -> google.protobuf.Any
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_pb_message_proto_init() }
func file_pb_message_proto_init() {
	if File_pb_message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pb_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
		file_pb_message_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Subscribe); i {
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
		file_pb_message_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Unsubscribe); i {
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
			RawDescriptor: file_pb_message_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pb_message_proto_goTypes,
		DependencyIndexes: file_pb_message_proto_depIdxs,
		MessageInfos:      file_pb_message_proto_msgTypes,
	}.Build()
	File_pb_message_proto = out.File
	file_pb_message_proto_rawDesc = nil
	file_pb_message_proto_goTypes = nil
	file_pb_message_proto_depIdxs = nil
}
