package test

import (
	"context"
	"time"
)

// CanceledContext return canceled ctx for testing canceled ctx
//
func CanceledContext() context.Context {
	dateline := time.Now().Add(time.Duration(1) * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), dateline)
	defer cancel()
	time.Sleep(time.Duration(2) * time.Millisecond)
	return ctx
}
