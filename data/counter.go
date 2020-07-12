package data

import (
	"context"
)

// CounterRef is a collection of documents (shards) to realize counter with high frequency.
//
type CounterRef interface {

	// CreateShards create counter and all shards, it is safe to create shards as many time as you want, normally we recreate shards when we need more shards
	//
	//	err = counter.CreateShards(ctx)
	//
	CreateShards(ctx context.Context) error

	// Increment increments a randomly picked shard.
	//
	//	err = counter.Increment(ctx, 2)
	//
	Increment(ctx context.Context, value int) error

	// Count returns a total count across all shards.
	//
	//	count, err = counter.Count(ctx)
	//
	Count(ctx context.Context) (int64, error)
}
