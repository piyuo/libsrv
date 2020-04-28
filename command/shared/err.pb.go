// Code generated by protoc-gen-go. DO NOT EDIT.
// source: err.proto

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
type Err struct {
	Code                 string   `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-" firestore:"-"`
	XXX_unrecognized     []byte   `json:"-" firestore:"-"`
	XXX_sizecache        int32    `json:"-" firestore:"-"`
}

func (m *Err) Reset()         { *m = Err{} }
func (m *Err) String() string { return proto.CompactTextString(m) }
func (*Err) ProtoMessage()    {}
func (*Err) Descriptor() ([]byte, []int) {
	return fileDescriptor_b4a1db73bc95ee8c, []int{0}
}

func (m *Err) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Err.Unmarshal(m, b)
}
func (m *Err) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Err.Marshal(b, m, deterministic)
}
func (m *Err) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Err.Merge(m, src)
}
func (m *Err) XXX_Size() int {
	return xxx_messageInfo_Err.Size(m)
}
func (m *Err) XXX_DiscardUnknown() {
	xxx_messageInfo_Err.DiscardUnknown(m)
}

var xxx_messageInfo_Err proto.InternalMessageInfo

func (m *Err) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func init() {
	proto.RegisterType((*Err)(nil), "Err")
}

func init() { proto.RegisterFile("err.proto", fileDescriptor_b4a1db73bc95ee8c) }

var fileDescriptor_b4a1db73bc95ee8c = []byte{
	// 64 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4c, 0x2d, 0x2a, 0xd2,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x57, 0x92, 0xe4, 0x62, 0x76, 0x2d, 0x2a, 0x12, 0x12, 0xe2, 0x62,
	0x49, 0xce, 0x4f, 0x49, 0x95, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x02, 0xb3, 0x93, 0xd8, 0xc0,
	0x2a, 0x8c, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0xec, 0x4c, 0xbc, 0x1b, 0x2e, 0x00, 0x00, 0x00,
}


func (m *Err) XXX_MapID() uint16 {
	return 0
}

func (m *Err) XXX_MapName() string {
	return "Err"
}
