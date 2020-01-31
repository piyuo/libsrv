package command

import (
	"context"

	shared "github.com/piyuo/go-libsrv/command/shared"
	log "github.com/piyuo/go-libsrv/log"
	//data "github.com/piyuo/go-libsrv/data"
)

// Main entry for client command execution
// you can return a response to user and error will be log to server
//
// do not return nil on response
func (a *TestAction) Main(ctx context.Context) (interface{}, error) {
	// get token if you need userid
	//token, errResp := shared.Token(ctx)
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

	// use log.Debug print message to the debug console
	log.Debug(ctx, "TestAction", "hi")

	// log significant events
	// log.Info(ctx, "%name", "hi")
	// log.Warning(ctx, "%name", "hi")
	// log.Critical(ctx, "%name", "hi")
	// no need to log error, just return error and client will get internal server error, error will log to google cloud

	// data operation
	//db, err := data.NewDB(ctx)
	//err = db.Put(ctx, &greet)
	//result := Greet{}
	//result.SetID(greet.ID())
	//err = db.Get(ctx, &result)

	// return custom response to client
	//return &StringResponse{Text: "hello"}

	// return error code to client
	//return Error(ErrorNeedJustLogin)

	//do not return nil, it will result internal server error
	// return shared.OK() if nothing else to return
	// return shared.Text("hi")
	// return shared.Number(101)
	return shared.OK(), nil
}
