package data

import (
	"context"
)

// CounterRef is a collection of documents (shards) to realize counter with high frequency.
//
type CounterRef interface {

	// CreateShards create shards document and collection, it is safe to create shards as many time as you want, normally we recreate shards when we need more shards
	//
	//	err = counter.CreateShards(ctx)
	//
	CreateShards(ctx context.Context) error

	// Increment increments a randomly picked shard. this function is slow than FastIncrement() but you don't need to create all shards first.
	//
	//	err = counter.Increment(ctx, 1)
	//
	Increment(ctx context.Context, value interface{}) error

	// FastIncrement increments a randomly picked shard. before use this function you must use createShard to create all necessary shard
	//
	//	err = counter.Increment(ctx, 1)
	//
	FastIncrement(ctx context.Context, value interface{}) error

	// Count returns a total count across all shards. avoid use this function in transation it easily cause "Too much contention on these documents"
	//
	//	count, err = counter.Count(ctx)
	//
	Count(ctx context.Context) (float64, error)
}
