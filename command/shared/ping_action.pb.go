// Code generated by protoc-gen-go. DO NOT EDIT.
// source: ping_action.proto

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
// say hi
//
// @returns StringResponse
type PingAction struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-" firestore:"-"`
	XXX_unrecognized     []byte   `json:"-" firestore:"-"`
	XXX_sizecache        int32    `json:"-" firestore:"-"`
}

func (m *PingAction) Reset()         { *m = PingAction{} }
func (m *PingAction) String() string { return proto.CompactTextString(m) }
func (*PingAction) ProtoMessage()    {}
func (*PingAction) Descriptor() ([]byte, []int) {
	return fileDescriptor_167f9ea2e1eb67c1, []int{0}
}

func (m *PingAction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PingAction.Unmarshal(m, b)
}
func (m *PingAction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PingAction.Marshal(b, m, deterministic)
}
func (m *PingAction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PingAction.Merge(m, src)
}
func (m *PingAction) XXX_Size() int {
	return xxx_messageInfo_PingAction.Size(m)
}
func (m *PingAction) XXX_DiscardUnknown() {
	xxx_messageInfo_PingAction.DiscardUnknown(m)
}

var xxx_messageInfo_PingAction proto.InternalMessageInfo

func init() {
	proto.RegisterType((*PingAction)(nil), "PingAction")
}

func init() { proto.RegisterFile("ping_action.proto", fileDescriptor_167f9ea2e1eb67c1) }

var fileDescriptor_167f9ea2e1eb67c1 = []byte{
	// 59 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2c, 0xc8, 0xcc, 0x4b,
	0x8f, 0x4f, 0x4c, 0x2e, 0xc9, 0xcc, 0xcf, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x57, 0xe2, 0xe1,
	0xe2, 0x0a, 0xc8, 0xcc, 0x4b, 0x77, 0x04, 0x8b, 0x25, 0xb1, 0x81, 0x05, 0x8d, 0x01, 0x01, 0x00,
	0x00, 0xff, 0xff, 0xcc, 0xce, 0xfa, 0x13, 0x29, 0x00, 0x00, 0x00,
}


func (m *PingAction) XXX_MapID() uint16 {
	return 4
}

func (m *PingAction) XXX_MapName() string {
	return "PingAction"
}
