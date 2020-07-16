package data

import (
	"context"
)

// CounterRef is a collection of documents (shards) to realize counter with high frequency.
//
type CounterRef interface {

	// IncrementRX increments a randomly picked shard. must used it in transaction with IncrementWX()
	//
	//	err = counter.IncrementRX(1)
	//
	IncrementRX(value interface{}) error

	// IncrementWX commit IncrementRX()
	//
	//	err = counter.IncrementWX()
	//
	IncrementWX() error

	// Count returns a total count across all shards. please be aware it easily cause "Too much contention on these documents"
	//
	//	count, err = counter.Count(ctx)
	//
	Count(ctx context.Context) (float64, error)

	// Reset reset counter
	//
	//	err = db.Transaction(ctx, func(ctx context.Context) error {
	//		err:= counter.Reset(ctx)
	//	})
	//
	Reset(ctx context.Context) error
}
