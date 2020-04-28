package command

import (
	"context"
	"errors"

	app "github.com/piyuo/go-libsrv/app"
	shared "github.com/piyuo/go-libsrv/command/shared"
)

//ErrAccessTokenRequired mean service need access  token
var ErrAccessTokenRequired = errors.New("access token required")

//ErrAccessTokenExpired mean access token is expired, client need use refresh token to get new access token
var ErrAccessTokenExpired = errors.New("access token expired")

//ErrPaymentTokenRequired mean service need access toke that generate from user enter password in 5 min
var ErrPaymentTokenRequired = errors.New("payment token required")

// Token return token or ErrorResponse
//
// 	token, errResp := shared.NeedToken(ctx)
// 	if errResp != nil {
// 		return errResp, nil
// 	}
func Token(ctx context.Context) (app.Token, error) {
	token, err := app.TokenFromContext(ctx)
	if err != nil {
		return nil, ErrAccessTokenRequired
	}
	if token.Expired() {
		return nil, ErrAccessTokenExpired
	}
	return token, nil
}

//OK return empty string
//
//	return shared.OK(),nil
func OK() interface{} {
	return &shared.Err{
		Code: "",
	}
}

//Error return  error response with code
//
//	return shared.Error(shared.ErrorUnknown),nil
func Error(errCode string) interface{} {
	return &shared.Err{
		Code: errCode,
	}
}

//String return string response
//
//	return shared.Text("hi"),nil
func String(text string) interface{} {
	return &shared.Text{
		Value: text,
	}
}

//Number return number response
//
//	return shared.Number(101),nil
func Number(num int64) interface{} {
	return &shared.Num{
		Value: num,
	}
}
