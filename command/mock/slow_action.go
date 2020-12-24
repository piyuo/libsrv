package mock

import (
	"context"
	"time"

	shared "github.com/piyuo/libsrv/command/shared"
	//data "github.com/piyuo/libsrv/data"
)

// Do return PbError
//
func (a *SlowAction) Do(ctx context.Context) (interface{}, error) {
	time.Sleep(time.Duration(2) * time.Millisecond)
	return &shared.PbError{
		Code: "",
	}, nil

}
