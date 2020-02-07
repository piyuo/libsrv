package shared

import (
	"context"
	"errors"

	app "github.com/piyuo/go-libsrv/app"
)

//ErrNeedToken no token mean we need user login
var ErrNeedToken = errors.New("need login")

//ErrorTokenExpired mean token expired and we need user login or refresh token
var ErrorTokenExpired = errors.New("need token")

//ErrorNeedLoginNow mean we are doing payment operation and  need confirm login again
var ErrorNeedLoginNow = errors.New("need login now")

// Token return token or ErrorResponse
//
// 	token, errResp := shared.NeedToken(ctx)
// 	if errResp != nil {
// 		return errResp, nil
// 	}
func Token(ctx context.Context) (app.Token, error) {
	token, err := app.TokenFromContext(ctx)
	if err != nil {
		return nil, ErrNeedToken
	}
	if token.Expired() {
		return nil, ErrorTokenExpired
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
func Error(code int32, msg string) interface{} {
	return errorInt32(code, msg)
}

func errorInt32(code int32, tag string) interface{} {
	return &Err{
		Code: code,
		Msg:  tag,
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
