package command

import (
	shared "github.com/piyuo/libsrv/command/shared"
)

var ok = &shared.PbError{
	Code: "",
}

// OK return empty string
//
//	return command.OK(),nil
//
func OK() interface{} {
	return ok
}

// NewPbError return  error response with code
//
//	return command.Error("INVALID_EMAIL")
//
func NewPbError(errCode string) interface{} {
	return &shared.PbError{
		Code: errCode,
	}
}

// IsOK return true if object is shared.Err and code is Empty
//
//	is := command.IsOK(response)
//
func IsOK(x interface{}) bool {
	return IsError(x, "")
}

// IsError return true if object is shared.Err and code is the same
//
//	is := command.IsError(response,"INVALID_EMAIL")
//
func IsError(x interface{}, errCode string) bool {
	if x == nil {
		return false
	}
	switch x.(type) {
	case *shared.PbError:
		e := x.(*shared.PbError)
		if e.Code == errCode {
			return true
		}
	}
	return false
}

// NewPbString return string response
//
//	return command.Text("hi")
//
func NewPbString(text string) interface{} {
	return &shared.PbString{
		Value: text,
	}
}

// NewPbInt return int response
//
//	return command.PbInt(101)
//
func NewPbInt(num int32) interface{} {
	return &shared.PbInt{
		Value: num,
	}
}

// NewPbBool return bool response
//
//	return command.NewPbBool(true)
//
func NewPbBool(value bool) interface{} {
	return &shared.PbBool{
		Value: value,
	}
}
