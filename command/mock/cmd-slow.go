package mock

import (
	"context"
	"time"

	"github.com/piyuo/libsrv/command/simple"
)

// Do comments
//
//	return OK if success
//	return "INVALID_XXX" when something wrong
//
func (c *CmdSlow) Do(ctx context.Context) (interface{}, error) {
	time.Sleep(time.Duration(2) * time.Millisecond)
	return &simple.Error{
		Code: "",
	}, nil
}
