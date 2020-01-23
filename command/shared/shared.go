package shared

import (
	"context"

	app "github.com/piyuo/go-libsrv/app"
	sharedcommands "github.com/piyuo/go-libsrv/commend/shared/commands"
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

// NeedToken return token or ErrorResponse
//
// 	token, errResp := shared.NeedToken(ctx)
// 	if errResp != nil {
// 		return errResp, nil
// 	}
func NeedToken(ctx context.Context) (app.Token, interface{}) {
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
func OK() interface{} {
	return &sharedcommands.Err{
		Code: 0,
	}
}

//Error return  error response with code
func Error(code ErrorCode, tag string) interface{} {
	return errorInt32(int32(code), tag)
}

func errorInt32(code int32, tag string) interface{} {
	return &sharedcommands.Err{
		Code: code,
		Tag:  tag,
	}
}
