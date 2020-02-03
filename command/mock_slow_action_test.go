package command

import (
	"context"
	"time"

	shared "github.com/piyuo/go-libsrv/command/shared"
)

func (a *SlowAction) Main(ctx context.Context) (interface{}, error) {
	time.Sleep(time.Duration(2) * time.Millisecond)
	return shared.OK(), nil
}
