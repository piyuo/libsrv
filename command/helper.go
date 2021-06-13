package command

import (
	"github.com/piyuo/libsrv/command/simple"
)

// ok instace, it use frequently
//
var OK = &simple.OK{}

// Error return error response with code
//
//	return command.Error("INVALID_EMAIL")
//
func Error(errCode string) interface{} {
	return &simple.Error{
		Code: errCode,
	}
}

// IsOK return true if object is PbOK
//
//	is := command.IsOK(response)
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
//	is := command.IsError(response,"INVALID_EMAIL")
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

// GetErrorCode return error code if object is PbError otherwise return empty
//
//	code := command.GetErrorCode(response) // "INVALID_EMAIL"
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
//	return command.Text("hi")
//
func String(text string) interface{} {
	return &simple.String{
		Value: text,
	}
}

// IsString return true if object is PbString and code is entered
//
//	is := command.IsString(response,"hi")
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

// Number return int response
//
//	return command.PbInt(101)
//
func Number(number int32) interface{} {
	return &simple.Number{
		Value: number,
	}
}

// IsInt return true if object is PbInt and value is entered
//
//	is := command.IsInt(response,42)
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
//	return command.NewPbBool(true)
//
func Bool(value bool) interface{} {
	return &simple.Bool{
		Value: value,
	}
}

// IsBool return true if object is PbBool and value is entered
//
//	is := command.IsBool(response,42)
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
