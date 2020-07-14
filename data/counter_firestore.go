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
	CounterRef `firestore:"-"`

	ShardsFirestore `firestore:"-"`
}

// Increment increments a randomly picked shard. this function is slow than FastIncrement() but you don't need to create all shards first.
//
//	err = counter.Increment(ctx, 1)
//
func (c *CounterFirestore) Increment(ctx context.Context, value interface{}) error {
	return c.increment(ctx, value, false)
}

// FastIncrement increments a randomly picked shard. before use this function you must use createShard to create all necessary shard
//
//	err = counter.Increment(ctx, 1)
//
func (c *CounterFirestore) FastIncrement(ctx context.Context, value interface{}) error {
	return c.increment(ctx, value, true)
}

// increments a randomly picked shard. if shardsCreated is false we need extra Get() to check is shard already exist?
//
//	err = counter.increment(ctx, 1, false)
//
func (c *CounterFirestore) increment(ctx context.Context, value interface{}, shardsCreated bool) error {
	if err := c.assert(ctx); err != nil {
		return err
	}
	docRef, shardsRef := c.getRef()
	shardID := strconv.Itoa(rand.Intn(c.numShards))
	shardRef := shardsRef.Doc(shardID)
	if c.conn.tx != nil {

		if err := c.ensureShardTx(c.conn.tx, docRef, shardRef); err != nil {
			return err
		}

		err := c.conn.tx.Update(shardRef, []firestore.Update{
			{Path: "C", Value: firestore.Increment(value)},
		})

		if err != nil {
			return errors.Wrap(err, "failed to update shard in transaction, you may need CreateShards() first: "+c.errorID()+"["+shardID+"]")
		}
		return nil
	}

	if err := c.ensureShard(ctx, docRef, shardRef); err != nil {
		return err
	}

	_, err := shardRef.Update(ctx, []firestore.Update{
		{Path: "C", Value: firestore.Increment(value)},
	})

	if err != nil {
		return errors.Wrap(err, "failed to update shard, you may need CreateShards() first: "+c.errorID()+"["+shardID+"]")
	}
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

	var shardCount float64
	var total float64
	shardCount = 0
	for {
		snotshot, err := shards.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, errors.Wrap(err, "failed to iterator shards: "+c.errorID())
		}

		shardCount++
		iTotal := snotshot.Data()["N"]
		shardCount, err := util.ToFloat64(iTotal)
		if err != nil {
			return 0, errors.Wrapf(err, "failed to get count on shards, invalid dataType %T, want float64: "+c.errorID(), iTotal)
		}
		total += shardCount
	}
	if shardCount == 0 {
		return 0, errors.New("failed to get any shards, you may need CreateShards() first: " + c.errorID())
	}
	return total, nil
}

// CreateShards create shards document and collection, it is safe to create shards as many time as you want, normally we recreate shards when we need more shards
//
//	err = counter.CreateShards(ctx)
//
func (c *CounterFirestore) CreateShards(ctx context.Context) error {
	return c.createShards(ctx)
}
