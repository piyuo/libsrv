package command

import (
	shared "github.com/piyuo/libsrv/command/shared"
)

// OK return empty string
//
//	return shared.OK(),nil
func OK() interface{} {
	return &shared.Err{
		Code: "",
	}
}

// Error return  error response with code
//
//	return shared.Error(shared.ErrorUnknown),nil
func Error(errCode string) interface{} {
	return &shared.Err{
		Code: errCode,
	}
}

// String return string response
//
//	return shared.Text("hi"),nil
func String(text string) interface{} {
	return &shared.Text{
		Value: text,
	}
}

// Number return number response
//
//	return shared.Number(101),nil
func Number(num int64) interface{} {
	return &shared.Num{
		Value: num,
	}
}
