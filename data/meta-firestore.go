package data

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/piyuo/libsrv/util"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

// MetaFirestore is parent of counter/coder/serial, provide basic function like createShard() and incrementShard()
//
type MetaFirestore struct {

	// conn is current firestore connection
	//
	conn *ConnectionFirestore

	// numShards is number of shards
	//
	numShards int

	// tableName  table name
	//
	tableName string `firestore:"-"`

	// id is document id in table
	//
	id string `firestore:"-"`
}

func (c *MetaFirestore) errorID() string {
	return c.conn.errorID(c.tableName, c.id)
}

// assert check ctx, table name, id are valid
//
func (c *MetaFirestore) assert(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if c.tableName == "" {
		return errors.New("table name can not be empty")
	}
	if c.id == "" {
		return errors.New("id can not be empty")
	}
	return nil
}

// Clear all shards
//
//	err = c.Clear(ctx)
//
func (c *MetaFirestore) clear(ctx context.Context) error {
	if err := c.assert(ctx); err != nil {
		return err
	}
	tableRef := c.conn.getCollectionRef(c.tableName)
	shards := tableRef.Where(MetaID, "==", c.id).Documents(ctx)

	batch := c.conn.client.Batch()
	var deleted = false
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

// shardsCount returns shards count
//
//	count, err = shardsCount(TypeCoder)
//
func (c *MetaFirestore) shardsCount(ctx context.Context) (int, error) {
	tableRef := c.conn.getCollectionRef(c.tableName)
	shards := tableRef.Where(MetaID, "==", c.id).Documents(ctx)
	defer shards.Stop()
	count := 0
	for {
		_, err := shards.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, errors.Wrap(err, "failed to iterator shards: "+c.errorID())
		}
		count++
	}
	return count, nil
}

// countValue returns a total value count on given shards
//
//	count, err = counter.countValue()
//
func (c *MetaFirestore) countValue(shards *firestore.DocumentIterator) (float64, error) {
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

		iTotal := snotshot.Data()[MetaValue]
		shardCount, err := util.ToFloat64(iTotal)
		if err != nil {
			return 0, errors.Wrapf(err, "failed to get count on shards, invalid dataType %T, want float64: "+c.errorID(), iTotal)
		}
		total += shardCount
	}
	return total, nil
}

// createShard create a shard
//
func (c *MetaFirestore) createShard(ref *firestore.DocumentRef, shard map[string]interface{}) error {
	err := c.conn.tx.Set(ref, shard, firestore.MergeAll)
	if err != nil {
		return errors.Wrap(err, "failed to create shard: "+c.errorID())
	}
	return nil
}

// incrementShard increment shard count
//
func (c *MetaFirestore) incrementShard(ref *firestore.DocumentRef, value interface{}) error {
	err := c.conn.tx.Update(ref, []firestore.Update{
		{Path: MetaValue, Value: firestore.Increment(value)},
	})
	if err != nil {
		return errors.Wrap(err, "failed to increment shard: "+c.errorID())
	}
	return nil
}
