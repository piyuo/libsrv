package command

import (
	"context"

	shared "github.com/piyuo/go-libsrv/command/shared"
	log "github.com/piyuo/go-libsrv/log"
	//data "github.com/piyuo/go-libsrv/data"
)

// Main entry for client command execution, need return response to client, client will get response use following code
//
//	var response = await service.send(action);
//
// return error will Intercept by command service on server and client. client response will be null
func (a *TestAction) Main(ctx context.Context) (interface{}, error) {
	// get token if you need user id
	//token, err := shared.Token(ctx)
	//if err != nil {
	//	return nil, err
	//}
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
	// log.Alert(ctx, "%name", "hi")
	// no need to log error, just return error

	// data operation
	//db, err := data.NewDB(ctx)
	//err = db.Put(ctx, &greet)
	//result := Greet{}
	//result.SetID(greet.ID())
	//err = db.Get(ctx, &result)

	// return custom response to client
	//return &StringResponse{Text: "hello"}

	// if you need return error to client, use
	//
	//return shared.Error(code,"error message")
	//
	// client will use following code to find is response error
	//	var response = await service.send(action);
	//	if(response != null && response is Err){print('response is error')}

	// return shared.Text("hi")
	// return shared.Number(101)
	// return shared.OK() if nothing else to return
	return shared.OK(), nil
}
