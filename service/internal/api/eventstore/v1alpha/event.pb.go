// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: eventstore/v1alpha/event.proto

package eventstorev1alpha

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	structpb "google.golang.org/protobuf/types/known/structpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Subject_Wildcard int32

const (
	Subject_WILDCARD_UNSPECIFIED  Subject_Wildcard = 0
	Subject_WILDCARD_SINGLE_TOKEN Subject_Wildcard = 1
	Subject_WILDCARD_MULTI_TOKEN  Subject_Wildcard = 2
)

// Enum value maps for Subject_Wildcard.
var (
	Subject_Wildcard_name = map[int32]string{
		0: "WILDCARD_UNSPECIFIED",
		1: "WILDCARD_SINGLE_TOKEN",
		2: "WILDCARD_MULTI_TOKEN",
	}
	Subject_Wildcard_value = map[string]int32{
		"WILDCARD_UNSPECIFIED":  0,
		"WILDCARD_SINGLE_TOKEN": 1,
		"WILDCARD_MULTI_TOKEN":  2,
	}
)

func (x Subject_Wildcard) Enum() *Subject_Wildcard {
	p := new(Subject_Wildcard)
	*p = x
	return p
}

func (x Subject_Wildcard) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Subject_Wildcard) Descriptor() protoreflect.EnumDescriptor {
	return file_eventstore_v1alpha_event_proto_enumTypes[0].Descriptor()
}

func (Subject_Wildcard) Type() protoreflect.EnumType {
	return &file_eventstore_v1alpha_event_proto_enumTypes[0]
}

func (x Subject_Wildcard) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Subject_Wildcard.Descriptor instead.
func (Subject_Wildcard) EnumDescriptor() ([]byte, []int) {
	return file_eventstore_v1alpha_event_proto_rawDescGZIP(), []int{4, 0}
}

type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Action    *Action                `protobuf:"bytes,2,opt,name=action,proto3" json:"action,omitempty"`
	Sequence  uint32                 `protobuf:"varint,3,opt,name=sequence,proto3" json:"sequence,omitempty"`
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_eventstore_v1alpha_event_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_eventstore_v1alpha_event_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Event.ProtoReflect.Descriptor instead.
func (*Event) Descriptor() ([]byte, []int) {
	return file_eventstore_v1alpha_event_proto_rawDescGZIP(), []int{0}
}

func (x *Event) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Event) GetAction() *Action {
	if x != nil {
		return x.Action
	}
	return nil
}

func (x *Event) GetSequence() uint32 {
	if x != nil {
		return x.Sequence
	}
	return 0
}

func (x *Event) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

type Aggregate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id              []string   `protobuf:"bytes,1,rep,name=id,proto3" json:"id,omitempty"`
	Commands        []*Command `protobuf:"bytes,2,rep,name=commands,proto3" json:"commands,omitempty"`
	CurrentSequence *uint32    `protobuf:"varint,3,opt,name=current_sequence,json=currentSequence,proto3,oneof" json:"current_sequence,omitempty"`
}

func (x *Aggregate) Reset() {
	*x = Aggregate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_eventstore_v1alpha_event_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Aggregate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Aggregate) ProtoMessage() {}

func (x *Aggregate) ProtoReflect() protoreflect.Message {
	mi := &file_eventstore_v1alpha_event_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Aggregate.ProtoReflect.Descriptor instead.
func (*Aggregate) Descriptor() ([]byte, []int) {
	return file_eventstore_v1alpha_event_proto_rawDescGZIP(), []int{1}
}

func (x *Aggregate) GetId() []string {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *Aggregate) GetCommands() []*Command {
	if x != nil {
		return x.Commands
	}
	return nil
}

func (x *Aggregate) GetCurrentSequence() uint32 {
	if x != nil && x.CurrentSequence != nil {
		return *x.CurrentSequence
	}
	return 0
}

type Command struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Action *Action `protobuf:"bytes,1,opt,name=action,proto3" json:"action,omitempty"`
}

func (x *Command) Reset() {
	*x = Command{}
	if protoimpl.UnsafeEnabled {
		mi := &file_eventstore_v1alpha_event_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Command) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Command) ProtoMessage() {}

func (x *Command) ProtoReflect() protoreflect.Message {
	mi := &file_eventstore_v1alpha_event_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Command.ProtoReflect.Descriptor instead.
func (*Command) Descriptor() ([]byte, []int) {
	return file_eventstore_v1alpha_event_proto_rawDescGZIP(), []int{2}
}

func (x *Command) GetAction() *Action {
	if x != nil {
		return x.Action
	}
	return nil
}

type Action struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Action   []string         `protobuf:"bytes,1,rep,name=action,proto3" json:"action,omitempty"`
	Revision uint32           `protobuf:"varint,2,opt,name=revision,proto3" json:"revision,omitempty"`
	Payload  *structpb.Struct `protobuf:"bytes,3,opt,name=payload,proto3,oneof" json:"payload,omitempty"`
}

func (x *Action) Reset() {
	*x = Action{}
	if protoimpl.UnsafeEnabled {
		mi := &file_eventstore_v1alpha_event_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Action) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Action) ProtoMessage() {}

func (x *Action) ProtoReflect() protoreflect.Message {
	mi := &file_eventstore_v1alpha_event_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Action.ProtoReflect.Descriptor instead.
func (*Action) Descriptor() ([]byte, []int) {
	return file_eventstore_v1alpha_event_proto_rawDescGZIP(), []int{3}
}

func (x *Action) GetAction() []string {
	if x != nil {
		return x.Action
	}
	return nil
}

func (x *Action) GetRevision() uint32 {
	if x != nil {
		return x.Revision
	}
	return 0
}

func (x *Action) GetPayload() *structpb.Struct {
	if x != nil {
		return x.Payload
	}
	return nil
}

type Subject struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Subject:
	//
	//	*Subject_Text
	//	*Subject_Wildcard_
	Subject isSubject_Subject `protobuf_oneof:"subject"`
}

func (x *Subject) Reset() {
	*x = Subject{}
	if protoimpl.UnsafeEnabled {
		mi := &file_eventstore_v1alpha_event_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Subject) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Subject) ProtoMessage() {}

func (x *Subject) ProtoReflect() protoreflect.Message {
	mi := &file_eventstore_v1alpha_event_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Subject.ProtoReflect.Descriptor instead.
func (*Subject) Descriptor() ([]byte, []int) {
	return file_eventstore_v1alpha_event_proto_rawDescGZIP(), []int{4}
}

func (m *Subject) GetSubject() isSubject_Subject {
	if m != nil {
		return m.Subject
	}
	return nil
}

func (x *Subject) GetText() string {
	if x, ok := x.GetSubject().(*Subject_Text); ok {
		return x.Text
	}
	return ""
}

func (x *Subject) GetWildcard() Subject_Wildcard {
	if x, ok := x.GetSubject().(*Subject_Wildcard_); ok {
		return x.Wildcard
	}
	return Subject_WILDCARD_UNSPECIFIED
}

type isSubject_Subject interface {
	isSubject_Subject()
}

type Subject_Text struct {
	Text string `protobuf:"bytes,1,opt,name=text,proto3,oneof"`
}

type Subject_Wildcard_ struct {
	Wildcard Subject_Wildcard `protobuf:"varint,2,opt,name=wildcard,proto3,enum=eventstore.v1alpha.Subject_Wildcard,oneof"`
}

func (*Subject_Text) isSubject_Subject() {}

func (*Subject_Wildcard_) isSubject_Subject() {}

var File_eventstore_v1alpha_event_proto protoreflect.FileDescriptor

var file_eventstore_v1alpha_event_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x61,
	0x6c, 0x70, 0x68, 0x61, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x12, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x61,
	0x6c, 0x70, 0x68, 0x61, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0xa2, 0x01, 0x0a, 0x05, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x32, 0x0a,
	0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70,
	0x68, 0x61, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x08, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x39, 0x0a,
	0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x99, 0x01, 0x0a, 0x09, 0x41, 0x67, 0x67,
	0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x37, 0x0a, 0x08, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x2e, 0x43, 0x6f,
	0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52, 0x08, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x73, 0x12,
	0x2e, 0x0a, 0x10, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x71, 0x75, 0x65,
	0x6e, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x48, 0x00, 0x52, 0x0f, 0x63, 0x75, 0x72,
	0x72, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x88, 0x01, 0x01, 0x42,
	0x13, 0x0a, 0x11, 0x5f, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x71, 0x75,
	0x65, 0x6e, 0x63, 0x65, 0x22, 0x3d, 0x0a, 0x07, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12,
	0x32, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x61,
	0x6c, 0x70, 0x68, 0x61, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x06, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x22, 0x80, 0x01, 0x0a, 0x06, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16,
	0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x76, 0x69, 0x73, 0x69,
	0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x72, 0x65, 0x76, 0x69, 0x73, 0x69,
	0x6f, 0x6e, 0x12, 0x36, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x48, 0x00, 0x52, 0x07,
	0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x88, 0x01, 0x01, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x70,
	0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0xc9, 0x01, 0x0a, 0x07, 0x53, 0x75, 0x62, 0x6a, 0x65,
	0x63, 0x74, 0x12, 0x14, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x48, 0x00, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x12, 0x42, 0x0a, 0x08, 0x77, 0x69, 0x6c, 0x64,
	0x63, 0x61, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x24, 0x2e, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x2e,
	0x53, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x2e, 0x57, 0x69, 0x6c, 0x64, 0x63, 0x61, 0x72, 0x64,
	0x48, 0x00, 0x52, 0x08, 0x77, 0x69, 0x6c, 0x64, 0x63, 0x61, 0x72, 0x64, 0x22, 0x59, 0x0a, 0x08,
	0x57, 0x69, 0x6c, 0x64, 0x63, 0x61, 0x72, 0x64, 0x12, 0x18, 0x0a, 0x14, 0x57, 0x49, 0x4c, 0x44,
	0x43, 0x41, 0x52, 0x44, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44,
	0x10, 0x00, 0x12, 0x19, 0x0a, 0x15, 0x57, 0x49, 0x4c, 0x44, 0x43, 0x41, 0x52, 0x44, 0x5f, 0x53,
	0x49, 0x4e, 0x47, 0x4c, 0x45, 0x5f, 0x54, 0x4f, 0x4b, 0x45, 0x4e, 0x10, 0x01, 0x12, 0x18, 0x0a,
	0x14, 0x57, 0x49, 0x4c, 0x44, 0x43, 0x41, 0x52, 0x44, 0x5f, 0x4d, 0x55, 0x4c, 0x54, 0x49, 0x5f,
	0x54, 0x4f, 0x4b, 0x45, 0x4e, 0x10, 0x02, 0x42, 0x09, 0x0a, 0x07, 0x73, 0x75, 0x62, 0x6a, 0x65,
	0x63, 0x74, 0x42, 0xe9, 0x01, 0x0a, 0x16, 0x63, 0x6f, 0x6d, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x42, 0x0a, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x5a, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x64, 0x6c, 0x65, 0x72, 0x68, 0x75, 0x72,
	0x73, 0x74, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31,
	0x61, 0x6c, 0x70, 0x68, 0x61, 0x3b, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x74, 0x6f, 0x72, 0x65,
	0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0xa2, 0x02, 0x03, 0x45, 0x58, 0x58, 0xaa, 0x02, 0x12,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x56, 0x31, 0x61, 0x6c, 0x70,
	0x68, 0x61, 0xca, 0x02, 0x12, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x5c,
	0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0xe2, 0x02, 0x1e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73,
	0x74, 0x6f, 0x72, 0x65, 0x5c, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x5c, 0x47, 0x50, 0x42,
	0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x13, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x73, 0x74, 0x6f, 0x72, 0x65, 0x3a, 0x3a, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_eventstore_v1alpha_event_proto_rawDescOnce sync.Once
	file_eventstore_v1alpha_event_proto_rawDescData = file_eventstore_v1alpha_event_proto_rawDesc
)

func file_eventstore_v1alpha_event_proto_rawDescGZIP() []byte {
	file_eventstore_v1alpha_event_proto_rawDescOnce.Do(func() {
		file_eventstore_v1alpha_event_proto_rawDescData = protoimpl.X.CompressGZIP(file_eventstore_v1alpha_event_proto_rawDescData)
	})
	return file_eventstore_v1alpha_event_proto_rawDescData
}

var file_eventstore_v1alpha_event_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_eventstore_v1alpha_event_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_eventstore_v1alpha_event_proto_goTypes = []interface{}{
	(Subject_Wildcard)(0),         // 0: eventstore.v1alpha.Subject.Wildcard
	(*Event)(nil),                 // 1: eventstore.v1alpha.Event
	(*Aggregate)(nil),             // 2: eventstore.v1alpha.Aggregate
	(*Command)(nil),               // 3: eventstore.v1alpha.Command
	(*Action)(nil),                // 4: eventstore.v1alpha.Action
	(*Subject)(nil),               // 5: eventstore.v1alpha.Subject
	(*timestamppb.Timestamp)(nil), // 6: google.protobuf.Timestamp
	(*structpb.Struct)(nil),       // 7: google.protobuf.Struct
}
var file_eventstore_v1alpha_event_proto_depIdxs = []int32{
	4, // 0: eventstore.v1alpha.Event.action:type_name -> eventstore.v1alpha.Action
	6, // 1: eventstore.v1alpha.Event.created_at:type_name -> google.protobuf.Timestamp
	3, // 2: eventstore.v1alpha.Aggregate.commands:type_name -> eventstore.v1alpha.Command
	4, // 3: eventstore.v1alpha.Command.action:type_name -> eventstore.v1alpha.Action
	7, // 4: eventstore.v1alpha.Action.payload:type_name -> google.protobuf.Struct
	0, // 5: eventstore.v1alpha.Subject.wildcard:type_name -> eventstore.v1alpha.Subject.Wildcard
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_eventstore_v1alpha_event_proto_init() }
func file_eventstore_v1alpha_event_proto_init() {
	if File_eventstore_v1alpha_event_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_eventstore_v1alpha_event_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Event); i {
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
		file_eventstore_v1alpha_event_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Aggregate); i {
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
		file_eventstore_v1alpha_event_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Command); i {
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
		file_eventstore_v1alpha_event_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Action); i {
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
		file_eventstore_v1alpha_event_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Subject); i {
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
	file_eventstore_v1alpha_event_proto_msgTypes[1].OneofWrappers = []interface{}{}
	file_eventstore_v1alpha_event_proto_msgTypes[3].OneofWrappers = []interface{}{}
	file_eventstore_v1alpha_event_proto_msgTypes[4].OneofWrappers = []interface{}{
		(*Subject_Text)(nil),
		(*Subject_Wildcard_)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_eventstore_v1alpha_event_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_eventstore_v1alpha_event_proto_goTypes,
		DependencyIndexes: file_eventstore_v1alpha_event_proto_depIdxs,
		EnumInfos:         file_eventstore_v1alpha_event_proto_enumTypes,
		MessageInfos:      file_eventstore_v1alpha_event_proto_msgTypes,
	}.Build()
	File_eventstore_v1alpha_event_proto = out.File
	file_eventstore_v1alpha_event_proto_rawDesc = nil
	file_eventstore_v1alpha_event_proto_goTypes = nil
	file_eventstore_v1alpha_event_proto_depIdxs = nil
}
