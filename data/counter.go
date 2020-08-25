package data

import (
	"context"
)

// Counter is a collection of documents (shards) to realize counter with high frequency.
//
type Counter interface {

	// Increment increments a randomly picked shard and generate count for all/year/month/day/hour
	//
	//	err = counter.Increment(ctx,1)
	//
	Increment(ctx context.Context, value interface{}) error

	// IncrementRX increments a randomly picked shard. must used it in transaction with IncrementWX()
	//
	//	err = counter.IncrementRX(1)
	//
	IncrementRX(ctx context.Context, value interface{}) error

	// IncrementWX commit IncrementRX()
	//
	//	err = counter.IncrementWX()
	//
	IncrementWX(ctx context.Context) error

	// CountAll return a total count across all period. this function not support transation cause it easily cause "Too much contention on these documents"
	//
	//	count, err = counter.CountAll(ctx)
	//
	CountAll(ctx context.Context) (float64, error)

	// Clear all counter shards
	//
	//	err = counter.Clear(ctx)
	//
	Clear(ctx context.Context) error
}

// Hierarchy define date hierarchy
//
type Hierarchy string

const (
	// HierarchyYear Define year period
	//
	HierarchyYear Hierarchy = "Y"

	// HierarchyMonth Define month period
	//
	HierarchyMonth = "M"

	// HierarchyDay Define day period
	//
	HierarchyDay = "D"

	// HierarchyHour Define hour period
	//
	HierarchyHour = "H"

	// HierarchyAll Define all period
	//
	HierarchyAll = "A"
)
