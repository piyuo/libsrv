package shared

import (
	"context"
	//data "github.com/piyuo/go-libsrv/data"
)

// Main entry for client command execution
// you can return a response to user and error will be log to server
//
// do not return nil on response
func (a *Num) Main(ctx context.Context) (interface{}, error) {
	return OK(), nil
}
