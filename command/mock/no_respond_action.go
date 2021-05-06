package mock

import (
	"context"
	//data "github.com/piyuo/libsrv/data"
)

// Do return nil
//
func (a *NoRespondAction) Do(ctx context.Context) (interface{}, error) {
	return nil, nil
}
