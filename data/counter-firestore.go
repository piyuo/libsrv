package data

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"github.com/piyuo/libsrv/util"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

// CounterFirestore implement Counter
//
type CounterFirestore struct {
	Counter `firestore:"-"`

	ShardsFirestore `firestore:"-"`

	loc *time.Location

	native time.Time

	incrementCallRX bool

	incrementValue interface{}

	incrementShardPick string

	incrementShardExist bool
}

// IncrementRX increments a randomly picked shard. must used it in transaction with IncrementWX()
//
//	err = counter.IncrementRX(ctx,1)
//
func (c *CounterFirestore) IncrementRX(ctx context.Context, value interface{}) error {
	if c.conn.tx == nil {
		return errors.New("This function must run in transaction")
	}
	c.incrementCallRX = true
	c.incrementValue = value
	c.incrementShardExist = false
	c.incrementShardPick = strconv.Itoa(rand.Intn(c.numShards)) //random pick a shard
	exist, err := c.isShardExist(ctx, c.incrementShardPick)
	if err != nil {
		return err
	}
	c.incrementShardExist = exist
	return nil
}

// getPickedAllRef return picked all period ref
//
func (c *CounterFirestore) getPickedAllRef() *firestore.DocumentRef {
	return c.conn.getDocRef(MetaDoc, MetaCounter+c.id+CounterPeriodAll+"."+c.incrementShardPick)
}

// isShardExist return true if shard already exist
//
func (c *CounterFirestore) isShardExist(ctx context.Context, pick string) (bool, error) {
	snapshot, err := c.conn.tx.Get(c.getPickedAllRef())
	if snapshot != nil && !snapshot.Exists() {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// IncrementWX commit IncrementRX()
//
//	err = counter.IncrementWX(ctx)
//
func (c *CounterFirestore) IncrementWX(ctx context.Context) error {
	if c.conn.tx == nil {
		return errors.New("This function must run in transaction")
	}
	if c.incrementCallRX == false {
		return errors.New("WX() function need call NumberRX() first")
	}

	year := strconv.Itoa(c.native.Year())
	month := strconv.Itoa(int(c.native.Month()))
	day := strconv.Itoa(int(c.native.Day()))
	hour := strconv.Itoa(int(c.native.Hour()))
	yearRef := c.conn.getDocRef(MetaDoc, MetaCounter+c.id+year+"."+c.incrementShardPick)
	monthRef := c.conn.getDocRef(MetaDoc, MetaCounter+c.id+year+"-"+month+"."+c.incrementShardPick)
	dayRef := c.conn.getDocRef(MetaDoc, MetaCounter+c.id+year+"-"+month+"-"+day+"."+c.incrementShardPick)
	hourRef := c.conn.getDocRef(MetaDoc, MetaCounter+c.id+year+"-"+month+"-"+day+"-"+hour+"."+c.incrementShardPick)

	if c.incrementShardExist {
		if err := c.incrementShard(ctx, c.getPickedAllRef()); err != nil {
			return errors.Wrap(err, "Failed to increment shard all")
		}
		if err := c.incrementShard(ctx, yearRef); err != nil {
			return errors.Wrap(err, "Failed to increment shard year")
		}
		if err := c.incrementShard(ctx, monthRef); err != nil {
			return errors.Wrap(err, "Failed to increment shard month")
		}
		if err := c.incrementShard(ctx, dayRef); err != nil {
			return errors.Wrap(err, "Failed to increment shard day")
		}
		if err := c.incrementShard(ctx, hourRef); err != nil {
			return errors.Wrap(err, "Failed to increment shard hour")
		}
	} else {
		shard := map[string]interface{}{
			CounterType:  MetaCounter,
			CounterID:    c.id,
			CounterCount: c.incrementValue,
		}

		shard[CounterPeriod] = HierarchyAll
		if err := c.createShard(ctx, c.getPickedAllRef(), shard); err != nil {
			return errors.Wrap(err, "Failed to create shard all")
		}

		shard[CounterPeriod] = HierarchyYear
		shard[CounterDate] = time.Date(c.native.Year(), time.Month(1), 01, 0, 0, 0, 0, c.loc)
		if err := c.createShard(ctx, yearRef, shard); err != nil {
			return errors.Wrap(err, "Failed to create shard year")
		}

		shard[CounterPeriod] = HierarchyMonth
		shard[CounterDate] = time.Date(c.native.Year(), c.native.Month(), 01, 0, 0, 0, 0, c.loc)
		if err := c.createShard(ctx, monthRef, shard); err != nil {
			return errors.Wrap(err, "Failed to create shard month")
		}

		shard[CounterPeriod] = HierarchyDay
		shard[CounterDate] = time.Date(c.native.Year(), c.native.Month(), c.native.Day(), 0, 0, 0, 0, c.loc)
		if err := c.createShard(ctx, dayRef, shard); err != nil {
			return errors.Wrap(err, "Failed to create shard day")
		}

		shard[CounterPeriod] = HierarchyHour
		shard[CounterDate] = time.Date(c.native.Year(), c.native.Month(), c.native.Day(), c.native.Hour(), 0, 0, 0, c.loc)
		if err := c.createShard(ctx, hourRef, shard); err != nil {
			return errors.Wrap(err, "Failed to create shard hour")
		}
	}
	c.incrementCallRX = false
	c.incrementValue = 0
	c.incrementShardExist = false
	c.incrementShardPick = ""
	return nil
}

// createShard create a shard
//
func (c *CounterFirestore) createShard(ctx context.Context, ref *firestore.DocumentRef, shard map[string]interface{}) error {
	err := c.conn.tx.Set(ref, shard, firestore.MergeAll)
	if err != nil {
		return errors.Wrap(err, "failed to create shard: "+ref.ID)
	}
	return nil
}

// incrementShard increment shard count
//
func (c *CounterFirestore) incrementShard(ctx context.Context, ref *firestore.DocumentRef) error {
	err := c.conn.tx.Update(ref, []firestore.Update{
		{Path: CounterCount, Value: firestore.Increment(c.incrementValue)},
	})
	if err != nil {
		return errors.Wrap(err, "failed to increment shard: "+ref.ID)
	}
	return nil
}

// countShards returns a total count on given shards
//
//	count, err = counter.countShards(ctx)
//
func (c *CounterFirestore) countShards(ctx context.Context, shards *firestore.DocumentIterator) (float64, error) {
	if err := c.assert(ctx); err != nil {
		return 0, err
	}

	defer shards.Stop()
	var total float64
	for {
		snotshot, err := shards.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, errors.Wrap(err, "failed to iterator shards: "+c.errorID())
		}

		iTotal := snotshot.Data()[CounterCount]
		shardCount, err := util.ToFloat64(iTotal)
		if err != nil {
			return 0, errors.Wrapf(err, "failed to get count on shards, invalid dataType %T, want float64: "+c.errorID(), iTotal)
		}
		total += shardCount
	}
	return total, nil
}

// CountAll return a total count across all period. this function not support transation cause it easily cause "Too much contention on these documents"
//
//	count, err = counter.CountAll(ctx)
//
func (c *CounterFirestore) CountAll(ctx context.Context) (float64, error) {
	metaRef := c.conn.getCollectionRef(MetaDoc)
	shards := metaRef.Where(CounterType, "==", MetaCounter).Where(CounterID, "==", c.id).Where(CounterPeriod, "==", HierarchyAll).Documents(ctx)
	return c.countShards(ctx, shards)
}

// CountPeriod return count between from and to. this function not support transation cause it easily cause "Too much contention on these documents"
//
//	count, err = counter.CountAll(ctx)
//
func (c *CounterFirestore) CountPeriod(ctx context.Context, period string, from, to time.Time) (float64, error) {
	metaRef := c.conn.getCollectionRef(MetaDoc)
	shards := metaRef.Where(CounterType, "==", MetaCounter).Where(CounterID, "==", c.id).Where(CounterPeriod, "==", HierarchyAll).Documents(ctx)
	return c.countShards(ctx, shards)
}

// Clear all counter shards
//
//	err = counter.Clear(ctx)
//
func (c *CounterFirestore) Clear(ctx context.Context) error {
	if err := c.assert(ctx); err != nil {
		return err
	}

	batch := c.conn.client.Batch()
	var deleted = false
	metaRef := c.conn.getCollectionRef(MetaDoc)
	shards := metaRef.Where(CounterType, "==", MetaCounter).Where(CounterID, "==", c.id).Documents(ctx)
	for {
		snotshot, err := shards.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return errors.Wrap(err, "failed to iterator shards: "+c.errorID())
		}
		deleted = true
		batch.Delete(snotshot.Ref)
	}
	if deleted {
		_, err := batch.Commit(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to commit batch")
		}
	}
	return nil
}
