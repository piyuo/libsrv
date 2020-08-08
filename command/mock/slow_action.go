package mock

import (
	"context"
	"time"

	shared "github.com/piyuo/libsrv/command/shared"
	//data "github.com/piyuo/libsrv/data"
)

// Main entry for client command execution
// you can return a response to user and error will be log to server
//
// do not return nil on response
func (a *SlowAction) Do(ctx context.Context) (interface{}, error) {
	time.Sleep(time.Duration(2) * time.Millisecond)
	return &shared.Err{
		Code: "",
	}, nil

}
