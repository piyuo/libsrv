package command

import (
	"context"

	app "github.com/piyuo/go-libsrv/app"
	"github.com/piyuo/go-libsrv/shared"
)

// Execute is main entry from client command
// you can return a response to user and error will be log to server
//
// do not return nil on response
func (a *TestAction) Main(ctx context.Context) (interface{}, error) {
	// get token if you need userid
	//token, errResp := shared.NeedToken(ctx)
	// if errResp != nil {
	// 	 return errResp, nil
	// }
	// if token.IsUserJustLogin() {
	// user just enter password to login iin 3 min. consider high security
	// a.placeOrder(token.UserID(),a.orderID)
	//} else {
	// user login by enter password or refresh token, consider low security
	// a.putToShopCart(token.UserID(),a.itemID)
	// }

	// use sys.LogInfo to print message to the console
	app.LogInfo(ctx, "hi")

	// log significant events to google cloud
	// sys.LogNotice(ctx, "hi")
	// sys.LogWarning(ctx, "hi")
	// sys.LogCritical(ctx, "hi")
	// sys.LogAlert(ctx, "hi")

	// no need to log error, just return error and client will get internal server error, error will log to google cloud

	// return custom response to client
	//return &StringResponse{Text: "hello"}

	// return error code to client
	//return Error(ErrorNeedJustLogin)

	// return shared.OK() if nothing else to return
	return shared.OK(), nil

	//do not return nil, it will result internal server error
}
