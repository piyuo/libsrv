package data

import (
	"context"
)

// CodeRef is a collection of documents (shards) to realize code with high frequency.
//
type CodeRef interface {
	/*
		// CreateShards create shards document and collection, it is safe to create shards as many time as you want, normally we recreate shards when we need more shards
		//
		//	err = code.CreateShards(ctx)
		//
		CreateShards(ctx context.Context) error
	*/
	// Code encode uint32 number into string, please be aware serial can only generate one number per second
	//
	//	code, err := code.Code(ctx)
	//	So(c, ShouldBeEmpty)
	//
	Code(ctx context.Context) (string, error)

	// Code16 encode uint16 number into string, please be aware serial can only generate one number per second
	//
	//	c, err := code.Code16(ctx)
	//	So(c, ShouldBeEmpty)
	//
	Code16(ctx context.Context) (string, error)

	// Code64 encode uint64 serial number to string, please be aware serial can only generate one number per second
	//
	//	c, err := code.Code64(ctx)
	//	So(c, ShouldBeEmpty)
	//
	Code64(ctx context.Context) (string, error)

	// Number return code number, number is unique but not serial
	//
	//	n, err := code.Number(ctx)
	//
	Number(ctx context.Context) (int64, error)
}
