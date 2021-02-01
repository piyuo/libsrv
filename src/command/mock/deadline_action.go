package mock

import (
	"context"
	//data "github.com/piyuo/libsrv/src/data"
)

// Do return DeadlineExceeded
//
func (a *DeadlineAction) Do(ctx context.Context) (interface{}, error) {
	return nil, context.DeadlineExceeded
}
