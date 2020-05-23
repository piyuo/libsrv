package mock

import (
	"context"
	//data "github.com/piyuo/libsrv/data"
)

// Main entry for client command execution
// you can return a response to user and error will be log to server
//
// do not return nil on response
func (a *NoRespondAction) Main(ctx context.Context) (interface{}, error) {
	return nil, nil
}
