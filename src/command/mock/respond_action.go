package mock

import (
	"context"

	"github.com/piyuo/libsrv/src/command/pb"
)

// Do comments
// you can return a response to user and error will be log to server
//
// do not return nil on response
func (a *RespondAction) Do(ctx context.Context) (interface{}, error) {
	return &pb.OK{}, nil
}
