package command

import (
	shared "github.com/piyuo/libsrv/command/shared"
)

// pbOK single instace cause it use frequently
//
var pbOK = &shared.PbOK{}

// OK return PbOK
//
//	return command.OK(),nil
//
func OK() interface{} {
	return pbOK
}

// Error return error response with code
//
//	return command.Error("INVALID_EMAIL")
//
func Error(errCode string) interface{} {
	return &shared.PbError{
		Code: errCode,
	}
}

// IsOK return true if object is PbOK
//
//	is := command.IsOK(response)
//
func IsOK(obj interface{}) bool {
	switch obj.(type) {
	case *shared.PbOK:
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
	switch obj.(type) {
	case *shared.PbError:
		e := obj.(*shared.PbError)
		if e.Code == entered {
			return true
		}
	}
	return false
}

// String return string response
//
//	return command.Text("hi")
//
func String(text string) interface{} {
	return &shared.PbString{
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
	switch obj.(type) {
	case *shared.PbString:
		s := obj.(*shared.PbString)
		if s.Value == entered {
			return true
		}
	}
	return false
}

// Int return int response
//
//	return command.PbInt(101)
//
func Int(num int32) interface{} {
	return &shared.PbInt{
		Value: num,
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
	switch obj.(type) {
	case *shared.PbInt:
		s := obj.(*shared.PbInt)
		if s.Value == entered {
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
	return &shared.PbBool{
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
	switch obj.(type) {
	case *shared.PbBool:
		s := obj.(*shared.PbBool)
		if s.Value == entered {
			return true
		}
	}
	return false
}
