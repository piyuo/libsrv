package shared

import (
	"context"

	app "github.com/piyuo/go-libsrv/app"
)

//ErrorCode use for code in ErrorResponse
type ErrorCode int32

const (
	//ErrorUnknown is unknown error, may happen  anywhere
	//
	//tag is log error id
	ErrorUnknown ErrorCode = 1

	//ErrorInternal is internal server error ,happen in action Main()
	//
	//tag is log error id
	ErrorInternal ErrorCode = 2

	//ErrorNeedToken need token
	ErrorNeedToken ErrorCode = 3

	//ErrorTokenExpired access token expired
	ErrorTokenExpired ErrorCode = 4

	//ErrorNeedJustLogin in high security like place order, we need user just enter  the password
	ErrorNeedJustLogin ErrorCode = 5
)

// Token return token or ErrorResponse
//
// 	token, errResp := shared.NeedToken(ctx)
// 	if errResp != nil {
// 		return errResp, nil
// 	}
func Token(ctx context.Context) (app.Token, interface{}) {
	token, err := app.TokenFromContext(ctx)
	if err != nil {
		return nil, Error(ErrorNeedToken, "")
	}
	if token.Expired() {
		return nil, Error(ErrorTokenExpired, "")
	}
	return token, nil
}

//OK return code=0 no error response
//
//	return shared.OK(),nil
func OK() interface{} {
	return &Err{
		Code: 0,
	}
}

//Error return  error response with code
//
//	return shared.Error(shared.ErrorUnknown),nil
func Error(code ErrorCode, tag string) interface{} {
	return errorInt32(int32(code), tag)
}

func errorInt32(code int32, tag string) interface{} {
	return &Err{
		Code: code,
		Tag:  tag,
	}
}

//String return string response
//
//	return shared.Text("hi"),nil
func String(text string) interface{} {
	return &Text{
		Value: text,
	}
}

//Number return number response
//
//	return shared.Number(101),nil
func Number(num int64) interface{} {
	return &Num{
		Value: num,
	}
}
