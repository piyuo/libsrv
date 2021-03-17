package db

import (
	"context"
	"time"
)

// Counter is a collection of documents (shards) to realize counter with high frequency.
//
type Counter interface {

	// IncrementRX increments a randomly picked shard. must used it in transaction with IncrementWX()
	//
	//	err := Transaction(ctx, func(ctx context.Context,tx db.Transaction) error {
	//		err = counter.IncrementRX(ctx,transaction)
	//		err = counter.IncrementWX(ctx,transaction,1)
	//	})
	//
	IncrementRX(ctx context.Context, transaction Transaction) error

	// IncrementWX commit IncrementRX()
	//
	//	err := Transaction(ctx, func(ctx context.Context,tx db.Transaction) error {
	//		err = counter.IncrementRX(ctx,transaction)
	//		err = counter.IncrementWX(ctx,transaction,1)
	//	})
	//
	IncrementWX(ctx context.Context, transaction Transaction, value interface{}) error

	// CountAll return a total count across all period. this function not support transation cause it easily cause "Too much contention on these documents"
	//
	//	count, err = counter.CountAll(ctx)
	//
	CountAll(ctx context.Context) (float64, error)

	// CountPeriod return count between from and to. this function not support transation cause it easily cause "Too much contention on these documents"
	//
	//	from := time.Date(now.Year()-1, 01, 01, 0, 0, 0, 0, time.UTC)
	//	to := time.Date(now.Year()+1, 01, 01, 0, 0, 0, 0, time.UTC)
	//	count, err := counter.CountPeriod(ctx, HierarchyYear, from, to)
	//
	CountPeriod(ctx context.Context, hierarchy Hierarchy, from, to time.Time) (float64, error)

	// DetailPeriod return detail between from and to. this function not support transation cause it easily cause "Too much contention on these documents"
	//
	//	dict, err = counter.DetailPeriod(ctx)
	//
	DetailPeriod(ctx context.Context, hierarchy Hierarchy, from, to time.Time) (map[time.Time]float64, error)

	// Clear all shards
	//
	//	err = counter.Clear(ctx,100)
	//
	Clear(ctx context.Context, max int) (bool, error)

	// ShardsCount returns shards count
	//
	//	count, err = counter.ShardsCount(ctx)
	//
	ShardsCount(ctx context.Context) (int, error)
}

// Hierarchy define date hierarchy
//
type Hierarchy string

const (
	// HierarchyYear define year period
	//
	HierarchyYear Hierarchy = "Year"

	// HierarchyMonth define month period
	//
	HierarchyMonth = "Month"

	// HierarchyDay define day period
	//
	HierarchyDay = "Day"

	// HierarchyHour define hour period
	//
	HierarchyHour = "Hour"

	// HierarchyTotal define total period
	//
	HierarchyTotal = "Total"
)

// CounterDateLevel field in shard
//
const CounterDateLevel = "Level"

// CounterTime field in shard
//
const CounterTime = "Time"
