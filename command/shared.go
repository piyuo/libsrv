package command

import (
	shared "github.com/piyuo/libsrv/command/shared"
)

var ok = &shared.Err{
	Code: "",
}

// OK return empty string
//
//	return command.OK(),nil
//
func OK() interface{} {
	return ok
}

// Error return  error response with code
//
//	return command.Error("INVALID_EMAIL")
//
func Error(errCode string) interface{} {
	return &shared.Err{
		Code: errCode,
	}
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
	case *shared.Err:
		e := x.(*shared.Err)
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
	return &shared.Text{
		Value: text,
	}
}

// Number return number response
//
//	return command.Number(101)
//
func Number(num int64) interface{} {
	return &shared.Num{
		Value: num,
	}
}
