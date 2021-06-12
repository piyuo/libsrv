package mock

import (
	"context"

	"github.com/piyuo/libsrv/command/types"
)

func (c *CmdBigData) Do(ctx context.Context) (interface{}, error) {
	return &types.String{
		Value: c.GetSample(),
	}, nil
}

// GetSample return large text sample
//
func (c *CmdBigData) GetSample() string {
	return "Go is expressive, concise, clean, and efficient. Its concurrency mechanisms make it easy to write programs that get the most out of multicore and networked machines, while its novel type system enables flexible and modular program construction. Go compiles quickly to machine code yet has the convenience of garbage collection and the power of run-time reflection. It's a fast, statically typed, compiled language that feels like a dynamically typed, interpreted language."
}
