// Code generated by protoc-gen-go. DO NOT EDIT.
// source: text.proto

package shared

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

//*
//  String Response
//
type Text struct {
	Value                string   `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-" firestore:"-"`
	XXX_unrecognized     []byte   `json:"-" firestore:"-"`
	XXX_sizecache        int32    `json:"-" firestore:"-"`
}

func (m *Text) Reset()         { *m = Text{} }
func (m *Text) String() string { return proto.CompactTextString(m) }
func (*Text) ProtoMessage()    {}
func (*Text) Descriptor() ([]byte, []int) {
	return fileDescriptor_a8e73d1ce47f9297, []int{0}
}

func (m *Text) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Text.Unmarshal(m, b)
}
func (m *Text) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Text.Marshal(b, m, deterministic)
}
func (m *Text) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Text.Merge(m, src)
}
func (m *Text) XXX_Size() int {
	return xxx_messageInfo_Text.Size(m)
}
func (m *Text) XXX_DiscardUnknown() {
	xxx_messageInfo_Text.DiscardUnknown(m)
}

var xxx_messageInfo_Text proto.InternalMessageInfo

func (m *Text) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func init() {
	proto.RegisterType((*Text)(nil), "Text")
}

func init() { proto.RegisterFile("text.proto", fileDescriptor_a8e73d1ce47f9297) }

var fileDescriptor_a8e73d1ce47f9297 = []byte{
	// 67 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2a, 0x49, 0xad, 0x28,
	0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x57, 0x92, 0xe1, 0x62, 0x09, 0x49, 0xad, 0x28, 0x11, 0x12,
	0xe1, 0x62, 0x2d, 0x4b, 0xcc, 0x29, 0x4d, 0x95, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x82, 0x70,
	0x92, 0xd8, 0xc0, 0x8a, 0x8c, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0xec, 0x35, 0x1a, 0x8d, 0x32,
	0x00, 0x00, 0x00,
}


func (m *Text) XXX_MapID() uint16 {
	return 3
}

func (m *Text) XXX_MapName() string {
	return "Text"
}
