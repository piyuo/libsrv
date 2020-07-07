package data

import (
	"context"
	"time"
)

// Counter is a collection of documents (shards) to realize counter with high frequency.
//
type Counter interface {

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

	// CreateTime return object create time
	//
	CreateTime() time.Time

	// SetCreateTime set object create time
	//
	SetCreateTime(time.Time)

	// ReadTime return object read time
	//
	ReadTime() time.Time

	// SetReadTime set object read time
	//
	SetReadTime(time.Time)

	// UpdateTime return object update time
	//
	UpdateTime() time.Time

	// SetUpdateTime set object update time
	//
	SetUpdateTime(time.Time)
}
