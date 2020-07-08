package data

import (
	"context"
	"time"
)

// CounterRef is a collection of documents (shards) to realize counter with high frequency.
//
type CounterRef interface {

	// initCounter creates a given number of shards as subcollection of specified document.
	//
	//	err = counter.init(ctx)
	//
	Init(ctx context.Context, numShards int) error

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

	// Delete counter and all shards.
	//
	//	err = counter.Delete(ctx)
	//
	Delete(ctx context.Context) error

	// GetCreateTime return object create time
	//
	GetCreateTime() time.Time

	// SetCreateTime set object create time
	//
	SetCreateTime(time.Time)

	// GetReadTime return object read time
	//
	GetReadTime() time.Time

	// SetReadTime set object read time
	//
	SetReadTime(time.Time)

	// GetUpdateTime return object update time
	//
	GetUpdateTime() time.Time

	// SetUpdateTime set object update time
	//
	SetUpdateTime(time.Time)
}
