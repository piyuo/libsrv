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

// IsOK return true if object is shared.Err and code is Empty
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

// IsError return true if object is shared.Err and code is the same
//
//	is := command.IsError(response,"INVALID_EMAIL")
//
func IsError(obj interface{}, errCode string) bool {
	if obj == nil {
		return false
	}
	switch obj.(type) {
	case *shared.PbError:
		e := obj.(*shared.PbError)
		if e.Code == errCode {
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

// Int return int response
//
//	return command.PbInt(101)
//
func Int(num int32) interface{} {
	return &shared.PbInt{
		Value: num,
	}
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
