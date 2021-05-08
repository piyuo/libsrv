// Code generated by protoc-gen-go. DO NOT EDIT.
// source: struct.proto

package structpb

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

// `NullValue` is a singleton enumeration to represent the null value for the
// `Value` type union.
//
//  The JSON representation for `NullValue` is JSON `null`.
type NullValue int32

const (
	// Null value.
	NullValue_NULL_VALUE NullValue = 0
)

var NullValue_name = map[int32]string{
	0: "NULL_VALUE",
}

var NullValue_value = map[string]int32{
	"NULL_VALUE": 0,
}

func (x NullValue) String() string {
	return proto.EnumName(NullValue_name, int32(x))
}

func (NullValue) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_0605f6bcb0ae6db1, []int{0}
}

func (NullValue) XXX_WellKnownType() string { return "NullValue" }

// `Struct` represents a structured data value, consisting of fields
// which map to dynamically typed values. In some languages, `Struct`
// might be supported by a native representation. For example, in
// scripting languages like JS a struct is represented as an
// object. The details of that representation are described together
// with the proto support for the language.
//
// The JSON representation for `Struct` is JSON object.
type Struct struct {
	// Unordered map of dynamically typed values.
	Fields               map[string]*Value `protobuf:"bytes,1,rep,name=fields,proto3" json:"fields,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Struct) Reset()         { *m = Struct{} }
func (m *Struct) String() string { return proto.CompactTextString(m) }
func (*Struct) ProtoMessage()    {}
func (*Struct) Descriptor() ([]byte, []int) {
	return fileDescriptor_0605f6bcb0ae6db1, []int{0}
}

func (*Struct) XXX_WellKnownType() string { return "Struct" }

func (m *Struct) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Struct.Unmarshal(m, b)
}
func (m *Struct) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Struct.Marshal(b, m, deterministic)
}
func (m *Struct) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Struct.Merge(m, src)
}
func (m *Struct) XXX_Size() int {
	return xxx_messageInfo_Struct.Size(m)
}
func (m *Struct) XXX_DiscardUnknown() {
	xxx_messageInfo_Struct.DiscardUnknown(m)
}

var xxx_messageInfo_Struct proto.InternalMessageInfo

func (m *Struct) GetFields() map[string]*Value {
	if m != nil {
		return m.Fields
	}
	return nil
}

// `Value` represents a dynamically typed value which can be either
// null, a number, a string, a boolean, a recursive struct value, or a
// list of values. A producer of value is expected to set one of that
// variants, absence of any variant indicates an error.
//
// The JSON representation for `Value` is JSON value.
type Value struct {
	// The kind of value.
	//
	// Types that are valid to be assigned to Kind:
	//	*Value_NullValue
	//	*Value_NumberValue
	//	*Value_StringValue
	//	*Value_BoolValue
	//	*Value_StructValue
	//	*Value_ListValue
	Kind                 isValue_Kind `protobuf_oneof:"kind"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Value) Reset()         { *m = Value{} }
func (m *Value) String() string { return proto.CompactTextString(m) }
func (*Value) ProtoMessage()    {}
func (*Value) Descriptor() ([]byte, []int) {
	return fileDescriptor_0605f6bcb0ae6db1, []int{1}
}

func (*Value) XXX_WellKnownType() string { return "Value" }

func (m *Value) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Value.Unmarshal(m, b)
}
func (m *Value) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Value.Marshal(b, m, deterministic)
}
func (m *Value) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Value.Merge(m, src)
}
func (m *Value) XXX_Size() int {
	return xxx_messageInfo_Value.Size(m)
}
func (m *Value) XXX_DiscardUnknown() {
	xxx_messageInfo_Value.DiscardUnknown(m)
}

var xxx_messageInfo_Value proto.InternalMessageInfo

type isValue_Kind interface {
	isValue_Kind()
}

type Value_NullValue struct {
	NullValue NullValue `protobuf:"varint,1,opt,name=null_value,json=nullValue,proto3,enum=google.protobuf.NullValue,oneof"`
}

type Value_NumberValue struct {
	NumberValue float64 `protobuf:"fixed64,2,opt,name=number_value,json=numberValue,proto3,oneof"`
}

type Value_StringValue struct {
	StringValue string `protobuf:"bytes,3,opt,name=string_value,json=stringValue,proto3,oneof"`
}

type Value_BoolValue struct {
	BoolValue bool `protobuf:"varint,4,opt,name=bool_value,json=boolValue,proto3,oneof"`
}

type Value_StructValue struct {
	StructValue *Struct `protobuf:"bytes,5,opt,name=struct_value,json=structValue,proto3,oneof"`
}

type Value_ListValue struct {
	ListValue *ListValue `protobuf:"bytes,6,opt,name=list_value,json=listValue,proto3,oneof"`
}

func (*Value_NullValue) isValue_Kind() {}

func (*Value_NumberValue) isValue_Kind() {}

func (*Value_StringValue) isValue_Kind() {}

func (*Value_BoolValue) isValue_Kind() {}

func (*Value_StructValue) isValue_Kind() {}

func (*Value_ListValue) isValue_Kind() {}

func (m *Value) GetKind() isValue_Kind {
	if m != nil {
		return m.Kind
	}
	return nil
}

func (m *Value) GetNullValue() NullValue {
	if x, ok := m.GetKind().(*Value_NullValue); ok {
		return x.NullValue
	}
	return NullValue_NULL_VALUE
}

func (m *Value) GetNumberValue() float64 {
	if x, ok := m.GetKind().(*Value_NumberValue); ok {
		return x.NumberValue
	}
	return 0
}

func (m *Value) GetStringValue() string {
	if x, ok := m.GetKind().(*Value_StringValue); ok {
		return x.StringValue
	}
	return ""
}

func (m *Value) GetBoolValue() bool {
	if x, ok := m.GetKind().(*Value_BoolValue); ok {
		return x.BoolValue
	}
	return false
}

func (m *Value) GetStructValue() *Struct {
	if x, ok := m.GetKind().(*Value_StructValue); ok {
		return x.StructValue
	}
	return nil
}

func (m *Value) GetListValue() *ListValue {
	if x, ok := m.GetKind().(*Value_ListValue); ok {
		return x.ListValue
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Value) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Value_NullValue)(nil),
		(*Value_NumberValue)(nil),
		(*Value_StringValue)(nil),
		(*Value_BoolValue)(nil),
		(*Value_StructValue)(nil),
		(*Value_ListValue)(nil),
	}
}

// `ListValue` is a wrapper around a repeated field of values.
//
// The JSON representation for `ListValue` is JSON array.
type ListValue struct {
	// Repeated field of dynamically typed values.
	Values               []*Value `protobuf:"bytes,1,rep,name=values,proto3" json:"values,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListValue) Reset()         { *m = ListValue{} }
func (m *ListValue) String() string { return proto.CompactTextString(m) }
func (*ListValue) ProtoMessage()    {}
func (*ListValue) Descriptor() ([]byte, []int) {
	return fileDescriptor_0605f6bcb0ae6db1, []int{2}
}

func (*ListValue) XXX_WellKnownType() string { return "ListValue" }

func (m *ListValue) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListValue.Unmarshal(m, b)
}
func (m *ListValue) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListValue.Marshal(b, m, deterministic)
}
func (m *ListValue) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListValue.Merge(m, src)
}
func (m *ListValue) XXX_Size() int {
	return xxx_messageInfo_ListValue.Size(m)
}
func (m *ListValue) XXX_DiscardUnknown() {
	xxx_messageInfo_ListValue.DiscardUnknown(m)
}

var xxx_messageInfo_ListValue proto.InternalMessageInfo

func (m *ListValue) GetValues() []*Value {
	if m != nil {
		return m.Values
	}
	return nil
}

func init() {
	proto.RegisterEnum("google.protobuf.NullValue", NullValue_name, NullValue_value)
	proto.RegisterType((*Struct)(nil), "google.protobuf.Struct")
	proto.RegisterMapType((map[string]*Value)(nil), "google.protobuf.Struct.FieldsEntry")
	proto.RegisterType((*Value)(nil), "google.protobuf.Value")
	proto.RegisterType((*ListValue)(nil), "google.protobuf.ListValue")
}

func init() { proto.RegisterFile("struct.proto", fileDescriptor_0605f6bcb0ae6db1) }

var fileDescriptor_0605f6bcb0ae6db1 = []byte{
	// 410 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x91, 0xc1, 0x6b, 0xd4, 0x40,
	0x14, 0xc6, 0x77, 0x92, 0x6e, 0x30, 0x2f, 0xa5, 0x96, 0x11, 0x74, 0xa9, 0xa0, 0x61, 0x7b, 0x09,
	0x22, 0x09, 0xae, 0x17, 0x31, 0x5e, 0x0c, 0xd4, 0x16, 0x0c, 0x25, 0x46, 0x5b, 0xc1, 0xcb, 0xd2,
	0xa4, 0x69, 0x0c, 0x9d, 0x9d, 0x09, 0xc9, 0x8c, 0xb2, 0x47, 0xff, 0x0b, 0xcf, 0x1e, 0x3d, 0xfa,
	0xd7, 0x79, 0x94, 0x99, 0xc9, 0xc4, 0xb2, 0xcb, 0x9e, 0x92, 0xf7, 0xe5, 0xf7, 0xbe, 0xf7, 0xbe,
	0x17, 0xd8, 0xef, 0x79, 0x27, 0x4a, 0x1e, 0xb6, 0x1d, 0xe3, 0x0c, 0xdf, 0xaf, 0x19, 0xab, 0x49,
	0xa5, 0xab, 0x42, 0xdc, 0xcc, 0x7f, 0x22, 0x70, 0x3e, 0x2a, 0x02, 0xc7, 0xe0, 0xdc, 0x34, 0x15,
	0xb9, 0xee, 0x67, 0xc8, 0xb7, 0x03, 0x6f, 0x71, 0x1c, 0x6e, 0xc0, 0xa1, 0x06, 0xc3, 0x77, 0x8a,
	0x3a, 0xa1, 0xbc, 0x5b, 0xe7, 0x43, 0xcb, 0xd1, 0x07, 0xf0, 0xee, 0xc8, 0xf8, 0x10, 0xec, 0xdb,
	0x6a, 0x3d, 0x43, 0x3e, 0x0a, 0xdc, 0x5c, 0xbe, 0xe2, 0xe7, 0x30, 0xfd, 0x76, 0x45, 0x44, 0x35,
	0xb3, 0x7c, 0x14, 0x78, 0x8b, 0x87, 0x5b, 0xe6, 0x97, 0xf2, 0x6b, 0xae, 0xa1, 0xd7, 0xd6, 0x2b,
	0x34, 0xff, 0x63, 0xc1, 0x54, 0x89, 0x38, 0x06, 0xa0, 0x82, 0x90, 0xa5, 0x36, 0x90, 0xa6, 0x07,
	0x8b, 0xa3, 0x2d, 0x83, 0x73, 0x41, 0x88, 0xe2, 0xcf, 0x26, 0xb9, 0x4b, 0x4d, 0x81, 0x8f, 0x61,
	0x9f, 0x8a, 0x55, 0x51, 0x75, 0xcb, 0xff, 0xf3, 0xd1, 0xd9, 0x24, 0xf7, 0xb4, 0x3a, 0x42, 0x3d,
	0xef, 0x1a, 0x5a, 0x0f, 0x90, 0x2d, 0x17, 0x97, 0x90, 0x56, 0x35, 0xf4, 0x14, 0xa0, 0x60, 0xcc,
	0xac, 0xb1, 0xe7, 0xa3, 0xe0, 0x9e, 0x1c, 0x25, 0x35, 0x0d, 0xbc, 0x31, 0xd7, 0x1e, 0x90, 0xa9,
	0x8a, 0xfa, 0x68, 0xc7, 0x1d, 0x07, 0x7b, 0x51, 0xf2, 0x31, 0x25, 0x69, 0x7a, 0xd3, 0xeb, 0xa8,
	0xde, 0xed, 0x94, 0x69, 0xd3, 0xf3, 0x31, 0x25, 0x31, 0x45, 0xe2, 0xc0, 0xde, 0x6d, 0x43, 0xaf,
	0xe7, 0x31, 0xb8, 0x23, 0x81, 0x43, 0x70, 0x94, 0x99, 0xf9, 0xa3, 0xbb, 0x8e, 0x3e, 0x50, 0xcf,
	0x1e, 0x83, 0x3b, 0x1e, 0x11, 0x1f, 0x00, 0x9c, 0x5f, 0xa4, 0xe9, 0xf2, 0xf2, 0x6d, 0x7a, 0x71,
	0x72, 0x38, 0x49, 0x7e, 0x20, 0x78, 0x50, 0xb2, 0xd5, 0xa6, 0x45, 0xe2, 0xe9, 0x34, 0x99, 0xac,
	0x33, 0xf4, 0xe5, 0x45, 0xdd, 0xf0, 0xaf, 0xa2, 0x08, 0x4b, 0xb6, 0x8a, 0x6a, 0x46, 0xae, 0x68,
	0x1d, 0x19, 0x34, 0x6a, 0xf9, 0xba, 0xad, 0xfa, 0x48, 0x87, 0x8e, 0xf5, 0xa3, 0x2d, 0xfe, 0x22,
	0xf4, 0xcb, 0xb2, 0x4f, 0xb3, 0xe4, 0xb7, 0xf5, 0xe4, 0x54, 0x9b, 0x67, 0x66, 0xbf, 0xcf, 0x15,
	0x21, 0xef, 0x29, 0xfb, 0x4e, 0x3f, 0xc9, 0xce, 0xc2, 0x51, 0x56, 0x2f, 0xff, 0x05, 0x00, 0x00,
	0xff, 0xff, 0x0a, 0x91, 0x40, 0xaf, 0xd5, 0x02, 0x00, 0x00,
}