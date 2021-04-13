package mock

import (
	"context"
	"time"

	"github.com/piyuo/libsrv/src/command/pb"
	//data "github.com/piyuo/libsrv/src/data"
)

// Do return PbError
//
func (a *SlowAction) Do(ctx context.Context) (interface{}, error) {
	time.Sleep(time.Duration(2) * time.Millisecond)
	return &pb.Error{
		Code: "",
	}, nil

}
