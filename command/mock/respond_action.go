package mock

import (
	"context"

	shared "github.com/piyuo/libsrv/command/shared"
)

// Main entry for client command execution
// you can return a response to user and error will be log to server
//
// do not return nil on response
func (a *RespondAction) Main(ctx context.Context) (interface{}, error) {
	return &shared.Err{
		Code: "",
	}, nil
}
