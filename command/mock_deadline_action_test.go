package command

import (
	"context"
)

func (a *DeadlineAction) Main(ctx context.Context) (interface{}, error) {
	return nil, context.DeadlineExceeded
}
