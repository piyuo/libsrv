package command

import (
	"github.com/piyuo/libsrv/command/simple"
)

const (
	BLOCK_SHORT = "BLOCK_SHORT"
	BLOCK_LONG  = "BLOCK_LONG"
)

// ok instace, it use frequently
//
var OK = &simple.OK{}

// BlockShort return error response with BLOCK_SHORT
//
//	return BlockShort
//
var BlockShort = &simple.Error{
	Code: BLOCK_SHORT,
}

// BlockLong return error response with BLOCK_Long
//
//	return BlockLong
//
var BlockLong = &simple.Error{
	Code: BLOCK_LONG,
}

// Error return error response with code
//
//	return Error("INVALID_EMAIL")
//
func Error(errCode string) interface{} {
	return &simple.Error{
		Code: errCode,
	}
}

// IsOK return true if object is PbOK
//
//	is := IsOK(response)
//
func IsOK(obj interface{}) bool {
	switch obj.(type) {
	case *simple.OK:
		return true
	}
	return false
}

// IsError return true if object is PbError and code is entered
//
//	is := IsError(response,"INVALID_EMAIL")
//
func IsError(obj interface{}, entered string) bool {
	if obj == nil {
		return false
	}
	switch obj := obj.(type) {
	case *simple.Error:
		if obj.Code == entered {
			return true
		}
	}
	return false
}

// IsBlockShort return true if object is BlockShort
//
//	is := IsBlockShort(response)
//
func IsBlockShort(obj interface{}) bool {
	return IsError(obj, BLOCK_SHORT)
}

// IsBlockLong return true if object is BlockLong
//
//	is := IsBlockLong(response)
//
func IsBlockLong(obj interface{}) bool {
	return IsError(obj, BLOCK_LONG)
}

// GetErrorCode return error code if object is PbError otherwise return empty
//
//	code := GetErrorCode(response) // "INVALID_EMAIL"
//
func GetErrorCode(obj interface{}) string {
	if obj == nil {
		return ""
	}
	switch obj := obj.(type) {
	case *simple.Error:
		return obj.Code
	}
	return ""
}

// String return string response
//
//	return Text("hi")
//
func String(text string) interface{} {
	return &simple.String{
		Value: text,
	}
}

// IsString return true if object is PbString and code is entered
//
//	is := IsString(response,"hi")
//
func IsString(obj interface{}, entered string) bool {
	if obj == nil {
		return false
	}
	switch obj := obj.(type) {
	case *simple.String:
		if obj.Value == entered {
			return true
		}
	}
	return false
}

// Int return int response
//
//	return Int(101)
//
func Int(number int32) interface{} {
	return &simple.Number{
		Value: number,
	}
}

// IsInt return true if object is PbInt and value is entered
//
//	is := IsInt(response,42)
//
func IsInt(obj interface{}, entered int32) bool {
	if obj == nil {
		return false
	}
	switch obj := obj.(type) {
	case *simple.Number:
		if obj.Value == entered {
			return true
		}
	}
	return false
}

// Bool return bool response
//
//	return Bool(true)
//
func Bool(value bool) interface{} {
	return &simple.Bool{
		Value: value,
	}
}

// IsBool return true if object is PbBool and value is entered
//
//	is := IsBool(response,42)
//
func IsBool(obj interface{}, entered bool) bool {
	if obj == nil {
		return false
	}
	switch obj := obj.(type) {
	case *simple.Bool:
		if obj.Value == entered {
			return true
		}
	}
	return false
}
