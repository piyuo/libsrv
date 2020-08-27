package data

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
)

// CounterFirestore implement Counter
//
type CounterFirestore struct {
	Counter `firestore:"-"`

	MetaFirestore `firestore:"-"`

	loc *time.Location

	native time.Time

	callRX bool

	value interface{}

	shardPick string

	shardExist bool
}

// isShardExist return true if shard already exist
//
func (c *CounterFirestore) isShardExist(ctx context.Context, ref *firestore.DocumentRef) (bool, error) {
	snapshot, err := c.conn.tx.Get(ref)
	if snapshot != nil && !snapshot.Exists() {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// IncrementRX increments a randomly picked shard. must used it in transaction with IncrementWX()
//
//	err = counter.IncrementRX(ctx,1)
//
func (c *CounterFirestore) IncrementRX(ctx context.Context, value interface{}) error {
	if c.conn.tx == nil {
		return errors.New("IncrementRX() must run in transaction")
	}
	c.callRX = true
	c.value = value
	c.shardPick = strconv.Itoa(rand.Intn(c.numShards)) //random pick a shard
	fmt.Printf("counter pick:" + c.shardPick + "\n")
	exist, err := c.isShardExist(ctx, c.getPickedAllRef())
	if err != nil {
		return err
	}
	c.shardExist = exist
	return nil
}

// getPickedAllRef return picked all period ref
//
func (c *CounterFirestore) getPickedAllRef() *firestore.DocumentRef {
	return c.conn.getDocRef(c.tableName, c.id+CounterPeriodAll+"."+c.shardPick)
}

// IncrementWX commit IncrementRX()
//
//	err = counter.IncrementWX(ctx)
//
func (c *CounterFirestore) IncrementWX(ctx context.Context) error {
	if c.conn.tx == nil {
		return errors.New("IncrementWX() must run in transaction")
	}
	if c.callRX == false {
		return errors.New("IncrementWX() need call IncrementRX() first")
	}

	year := strconv.Itoa(c.native.Year())
	month := strconv.Itoa(int(c.native.Month()))
	day := strconv.Itoa(int(c.native.Day()))
	hour := strconv.Itoa(int(c.native.Hour()))
	yearRef := c.conn.getDocRef(c.tableName, c.id+year+"."+c.shardPick)
	monthRef := c.conn.getDocRef(c.tableName, c.id+year+"-"+month+"."+c.shardPick)
	dayRef := c.conn.getDocRef(c.tableName, c.id+year+"-"+month+"-"+day+"."+c.shardPick)
	hourRef := c.conn.getDocRef(c.tableName, c.id+year+"-"+month+"-"+day+"-"+hour+"."+c.shardPick)

	if c.shardExist {
		if err := c.incrementShard(c.getPickedAllRef(), c.value); err != nil {
			return errors.Wrap(err, "Failed to increment shard all")
		}
		if err := c.incrementShard(yearRef, c.value); err != nil {
			return errors.Wrap(err, "Failed to increment shard year")
		}
		if err := c.incrementShard(monthRef, c.value); err != nil {
			return errors.Wrap(err, "Failed to increment shard month")
		}
		if err := c.incrementShard(dayRef, c.value); err != nil {
			return errors.Wrap(err, "Failed to increment shard day")
		}
		if err := c.incrementShard(hourRef, c.value); err != nil {
			return errors.Wrap(err, "Failed to increment shard hour")
		}
	} else {
		shard := map[string]interface{}{
			MetaID:    c.id,
			MetaValue: c.value,
		}

		shard[CounterPeriod] = HierarchyAll
		if err := c.createShard(c.getPickedAllRef(), shard); err != nil {
			return errors.Wrap(err, "Failed to create shard all")
		}

		shard[CounterPeriod] = HierarchyYear
		shard[CounterDate] = time.Date(c.native.Year(), time.Month(1), 01, 0, 0, 0, 0, c.loc)
		if err := c.createShard(yearRef, shard); err != nil {
			return errors.Wrap(err, "Failed to create shard year")
		}

		shard[CounterPeriod] = HierarchyMonth
		shard[CounterDate] = time.Date(c.native.Year(), c.native.Month(), 01, 0, 0, 0, 0, c.loc)
		if err := c.createShard(monthRef, shard); err != nil {
			return errors.Wrap(err, "Failed to create shard month")
		}

		shard[CounterPeriod] = HierarchyDay
		shard[CounterDate] = time.Date(c.native.Year(), c.native.Month(), c.native.Day(), 0, 0, 0, 0, c.loc)
		if err := c.createShard(dayRef, shard); err != nil {
			return errors.Wrap(err, "Failed to create shard day")
		}

		shard[CounterPeriod] = HierarchyHour
		shard[CounterDate] = time.Date(c.native.Year(), c.native.Month(), c.native.Day(), c.native.Hour(), 0, 0, 0, c.loc)
		if err := c.createShard(hourRef, shard); err != nil {
			return errors.Wrap(err, "Failed to create shard hour")
		}
	}
	c.callRX = false
	c.value = 0
	c.shardExist = false
	c.shardPick = ""
	return nil
}

// CountAll return a total count across all period. this function not support transation cause it easily cause "Too much contention on these documents"
//
//	count, err = counter.CountAll(ctx)
//
func (c *CounterFirestore) CountAll(ctx context.Context) (float64, error) {
	tableRef := c.conn.getCollectionRef(c.tableName)
	shards := tableRef.Where(MetaID, "==", c.id).Where(CounterPeriod, "==", HierarchyAll).Documents(ctx)
	return c.countValue(shards)
}

// CountPeriod return count between from and to. this function not support transation cause it easily cause "Too much contention on these documents"
//
//	count, err = counter.CountAll(ctx)
//
func (c *CounterFirestore) CountPeriod(ctx context.Context, period string, from, to time.Time) (float64, error) {
	tableRef := c.conn.getCollectionRef(c.tableName)
	shards := tableRef.Where(MetaID, "==", c.id).Where(CounterPeriod, "==", HierarchyAll).Documents(ctx)
	return c.countValue(shards)
}

// Clear all shards
//
//	err = c.Clear(ctx)
//
func (c *CounterFirestore) Clear(ctx context.Context) error {
	return c.clear(ctx)
}

// ShardsCount returns shards count
//
//	count, err = coder.ShardsCount()
//
func (c *CounterFirestore) ShardsCount(ctx context.Context) (int, error) {
	return c.shardsCount(ctx)
}
