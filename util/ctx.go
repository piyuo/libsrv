package util

import (
	"context"
	"time"
)

// CanceledCtx return canceled ctx for testing canceled ctx
//
func CanceledCtx() context.Context {
	dateline := time.Now().Add(time.Duration(1) * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), dateline)
	defer cancel()
	time.Sleep(time.Duration(2) * time.Millisecond)
	return ctx
}
