package data

import (
	"context"
	"time"
)

// CounterRef is a collection of documents (shards) to realize counter with high frequency.
//
type CounterRef interface {

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

	// GetCreateTime return object create time
	//
	GetCreateTime() time.Time

	// GetReadTime return object read time
	//
	GetReadTime() time.Time

	// GetUpdateTime return object update time
	//
	GetUpdateTime() time.Time
}
