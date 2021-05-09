package gdb

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/piyuo/libsrv/db"
	"github.com/piyuo/libsrv/log"
	"github.com/pkg/errors"
)

// CounterFirestore implement Counter
//
type CounterFirestore struct {
	db.Counter `firestore:"-"`

	MetaFirestore `firestore:"-"`

	callRX bool

	// pickedShard is a shard random picked
	//
	pickedShard string

	// shardExists return true if shard exists
	//
	shardExists bool
}

// shardAllRef return picked all period ref
//
func (c *CounterFirestore) shardAllRef() *firestore.DocumentRef {
	return c.client.getDocRef(c.collection, c.id+string(db.HierarchyTotal)+"_"+c.pickedShard)
}

// IncrementRX increments a randomly picked shard. must used it in transaction with IncrementWX()
//
//	err := Transaction(ctx, func(ctx context.Context,tx db.Transaction) error {
//		err = counter.IncrementRX(ctx,transaction)
//		err = counter.IncrementWX(ctx,transaction,1)
//	})
//
func (c *CounterFirestore) IncrementRX(ctx context.Context, transaction db.Transaction) error {
	tx := transaction.(*TransactionFirestore)
	c.callRX = true
	c.pickedShard = strconv.Itoa(rand.Intn(c.numShards)) //random pick a shard
	log.Debug(ctx, "counter pick %v from %v shards", c.pickedShard, c.numShards)
	var err error
	c.shardExists, err = tx.isShardExists(ctx, c.shardAllRef())
	if err != nil {
		return errors.Wrap(err, "all")
	}
	return nil
}

// IncrementWX commit IncrementRX()
//
//	err := Transaction(ctx, func(ctx context.Context,tx db.Transaction) error {
//		err = counter.IncrementRX(ctx,transaction)
//		err = counter.IncrementWX(ctx,transaction,1)
//	})
//
func (c *CounterFirestore) IncrementWX(ctx context.Context, transaction db.Transaction, value interface{}) error {
	tx := transaction.(*TransactionFirestore)
	if !c.callRX {
		return errors.New("must call RX first")
	}

	utcNow := time.Now().UTC()
	shard := map[string]interface{}{
		db.MetaID:      c.id,
		db.MetaN:       value,
		db.CounterTime: utcNow,
	}

	if c.shardExists {
		if err := tx.incrementShard(c.shardAllRef(), value); err != nil {
			return errors.Wrap(err, "inc all")
		}
	} else {
		shard[db.CounterDateLevel] = db.HierarchyTotal
		if err := tx.createShard(c.shardAllRef(), shard); err != nil {
			return errors.Wrap(err, "create all")
		}

	}
	c.callRX = false
	return nil
}

// CountAll return a total count across all period. this function not support transation cause it easily cause "Too much contention on these documents"
//
//	count, err = counter.CountAll(ctx)
//
func (c *CounterFirestore) CountAll(ctx context.Context) (float64, error) {
	tableRef := c.client.getCollectionRef(c.collection)
	shards := tableRef.Where(db.MetaID, "==", c.id).Where(db.CounterDateLevel, "==", db.HierarchyTotal).Documents(ctx)
	defer shards.Stop()
	return c.countValue(shards)
}

// CountPeriod return count between from and to. this function not support transation cause it easily cause "Too much contention on these documents"
//
//	from := time.Date(now.Year()-1, 01, 01, 0, 0, 0, 0, time.UTC)
//	to := time.Date(now.Year()+1, 01, 01, 0, 0, 0, 0, time.UTC)
//	count, err := counter.CountPeriod(ctx, HierarchyYear, from, to)
//
func (c *CounterFirestore) CountPeriod(ctx context.Context, hierarchy db.Hierarchy, from, to time.Time) (float64, error) {
	tableRef := c.client.getCollectionRef(c.collection)
	shards := tableRef.Where(db.MetaID, "==", c.id).Where(db.CounterDateLevel, "==", string(hierarchy)).Where(db.CounterTime, ">=", from).Where(db.CounterTime, "<=", to).Documents(ctx)
	defer shards.Stop()
	return c.countValue(shards)
}

// Delete delete counter
//
//	err = Delete(ctx)
//
func (c *CounterFirestore) Delete(ctx context.Context) error {
	return c.deleteShards(ctx)
}

// ShardsCount returns shards count
//
//	count, err = ShardsCount(ctx)
//
func (c *CounterFirestore) ShardsCount(ctx context.Context) (int, error) {
	return c.shardsCount(ctx)
}
