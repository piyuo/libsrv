// Code generated by protoc-gen-go. DO NOT EDIT.
// source: num.proto

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
//  Number Response
//
type Num struct {
	Value                int64    `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-" firestore:"-"`
	XXX_unrecognized     []byte   `json:"-" firestore:"-"`
	XXX_sizecache        int32    `json:"-" firestore:"-"`
}

func (m *Num) Reset()         { *m = Num{} }
func (m *Num) String() string { return proto.CompactTextString(m) }
func (*Num) ProtoMessage()    {}
func (*Num) Descriptor() ([]byte, []int) {
	return fileDescriptor_b52d19e3737dc567, []int{0}
}

func (m *Num) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Num.Unmarshal(m, b)
}
func (m *Num) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Num.Marshal(b, m, deterministic)
}
func (m *Num) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Num.Merge(m, src)
}
func (m *Num) XXX_Size() int {
	return xxx_messageInfo_Num.Size(m)
}
func (m *Num) XXX_DiscardUnknown() {
	xxx_messageInfo_Num.DiscardUnknown(m)
}

var xxx_messageInfo_Num proto.InternalMessageInfo

func (m *Num) GetValue() int64 {
	if m != nil {
		return m.Value
	}
	return 0
}

func init() {
	proto.RegisterType((*Num)(nil), "Num")
}

func init() { proto.RegisterFile("num.proto", fileDescriptor_b52d19e3737dc567) }

var fileDescriptor_b52d19e3737dc567 = []byte{
	// 65 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcc, 0x2b, 0xcd, 0xd5,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x57, 0x92, 0xe6, 0x62, 0xf6, 0x2b, 0xcd, 0x15, 0x12, 0xe1, 0x62,
	0x2d, 0x4b, 0xcc, 0x29, 0x4d, 0x95, 0x60, 0x54, 0x60, 0xd4, 0x60, 0x0e, 0x82, 0x70, 0x92, 0xd8,
	0xc0, 0x6a, 0x8c, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x78, 0xdb, 0xc5, 0x5f, 0x30, 0x00, 0x00,
	0x00,
}


func (m *Num) XXX_MapID() uint16 {
	return 2
}

func (m *Num) XXX_MapName() string {
	return "Num"
}
