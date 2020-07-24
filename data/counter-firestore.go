package data

import (
	"context"
	"math/rand"
	"strconv"

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

	incrementCallRX bool

	incrementValue interface{}

	incrementShardIndex int

	incrementCanCreateShard bool

	incrementCanIncrementShard bool
}

// IncrementRX increments a randomly picked shard. must used it in transaction with IncrementWX()
//
//	err = counter.IncrementRX(1)
//
func (c *CounterFirestore) IncrementRX(value interface{}) error {
	if c.conn.tx == nil {
		return errors.New("This function must run in transaction")
	}

	_, shardsRef := c.getRef()
	c.incrementCallRX = true
	c.incrementValue = value
	c.incrementCanCreateShard = false
	c.incrementCanIncrementShard = false
	c.incrementShardIndex = rand.Intn(c.numShards)
	shardID := strconv.Itoa(c.incrementShardIndex)
	shardRef := shardsRef.Doc(shardID)

	snapshot, err := c.conn.tx.Get(shardRef)
	if snapshot != nil && !snapshot.Exists() {
		c.incrementCanCreateShard = true
		return nil
	}

	if err != nil {
		return errors.Wrap(err, "failed to get shard: "+c.errorID())
	}

	c.incrementCanIncrementShard = true
	return nil

}

// IncrementWX commit IncrementRX()
//
//	err = counter.IncrementWX()
//
func (c *CounterFirestore) IncrementWX() error {
	if c.conn.tx == nil {
		return errors.New("This function must run in transaction")
	}
	if c.incrementCallRX == false {
		return errors.New("WX() function need call NumberRX() first")
	}

	docRef, shardsRef := c.getRef()
	shardID := strconv.Itoa(c.incrementShardIndex)
	shardRef := shardsRef.Doc(shardID)
	if c.incrementCanCreateShard {
		// create shards document
		err := c.conn.tx.Set(docRef, &struct{}{}) //put empty struct
		if err != nil {
			return errors.Wrap(err, "failed to create shards document: "+c.errorID())
		}

		// create shard
		err = c.conn.tx.Set(shardRef, map[string]interface{}{"N": c.incrementValue}, firestore.MergeAll)
		if err != nil {
			return errors.Wrap(err, "failed to create shard: "+c.errorID())
		}
	}

	if c.incrementCanIncrementShard {
		err := c.conn.tx.Update(shardRef, []firestore.Update{
			{Path: "N", Value: firestore.Increment(c.incrementValue)},
		})
		if err != nil {
			return errors.Wrap(err, "failed to increment shard: "+c.errorID())
		}
	}
	c.incrementCallRX = false
	c.incrementValue = 0
	c.incrementCanCreateShard = false
	c.incrementCanIncrementShard = false
	c.incrementShardIndex = -1
	return nil

}

// Count returns a total count across all shards. avoid use this function in transation it easily cause "Too much contention on these documents"
//
//	count, err = counter.Count(ctx)
//
func (c *CounterFirestore) Count(ctx context.Context) (float64, error) {
	if err := c.assert(ctx); err != nil {
		return 0, err
	}

	_, shardsRef := c.getRef()
	var shards *firestore.DocumentIterator
	if c.conn.tx != nil {
		shards = c.conn.tx.Documents(shardsRef)
	} else {
		shards = shardsRef.Documents(ctx)
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

		iTotal := snotshot.Data()["N"]
		shardCount, err := util.ToFloat64(iTotal)
		if err != nil {
			return 0, errors.Wrapf(err, "failed to get count on shards, invalid dataType %T, want float64: "+c.errorID(), iTotal)
		}
		total += shardCount
	}
	return total, nil
}

// Reset reset counter
//
//	err = db.Transaction(ctx, func(ctx context.Context) error {
//		err:= counter.Reset(ctx)
//	})
//
func (c *CounterFirestore) Reset(ctx context.Context) error {
	if err := c.assert(ctx); err != nil {
		return err
	}
	return c.deleteShards(ctx)
}
