package mock

import (
	"context"

	shared "github.com/piyuo/libsrv/src/command/shared"
)

// Do comments
// you can return a response to user and error will be log to server
//
// do not return nil on response
func (a *RespondAction) Do(ctx context.Context) (interface{}, error) {
	return &shared.PbOK{}, nil
}
