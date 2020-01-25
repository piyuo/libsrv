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
	log.Info(ctx, "hi")
	// log.Notice(ctx, "hi")
	// log.Warning(ctx, "hi")
	// log.Critical(ctx, "hi")
	// log.Alert(ctx, "hi")
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

	// return shared.OK() if nothing else to return
	//do not return nil, it will result internal server error
	return shared.OK(), nil
}
