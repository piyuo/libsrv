package mock

import (
	"context"

	shared "github.com/piyuo/libsrv/src/command/shared"
)

// GetSample return large text sample
//
func (a *BigDataAction) GetSample() string {
	return "Go is expressive, concise, clean, and efficient. Its concurrency mechanisms make it easy to write programs that get the most out of multicore and networked machines, while its novel type system enables flexible and modular program construction. Go compiles quickly to machine code yet has the convenience of garbage collection and the power of run-time reflection. It's a fast, statically typed, compiled language that feels like a dynamically typed, interpreted language."
}

// Do something
// you can return a response to user and error will be log to server
//
// do not return nil on response
func (a *BigDataAction) Do(ctx context.Context) (interface{}, error) {
	return &shared.PbString{
		Value: a.GetSample(),
	}, nil
}
