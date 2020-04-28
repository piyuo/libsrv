package command

import (
	"context"
	"time"
)

func (a *SlowAction) Main(ctx context.Context) (interface{}, error) {
	time.Sleep(time.Duration(2) * time.Millisecond)
	return OK(), nil
}
