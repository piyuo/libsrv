package gdb

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/piyuo/libsrv/src/db"
	"github.com/piyuo/libsrv/src/util"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

// MetaFirestore is parent of counter/coder/serial, provide basic function like createShard() and incrementShard()
//
type MetaFirestore struct {

	// client is db client
	//
	client *ClientFirestore

	// numShards is number of shards
	//
	numShards int

	// collection is collection name
	//
	collection string `firestore:"-"`

	// id is document id in collection
	//
	id string `firestore:"-"`
}

// check check ctx, table name, id are valid
//
func (c *MetaFirestore) check(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if c.collection == "" {
		return errors.New("collection must not empty")
	}
	if c.id == "" {
		return errors.New("id must no empty")
	}
	return nil
}

// deleteShards delete all shards
//
//	err = deleteShards(ctx)
//
func (c *MetaFirestore) deleteShards(ctx context.Context) error {
	if err := c.check(ctx); err != nil {
		return err
	}
	tableRef := c.client.getCollectionRef(c.collection)
	shardsIter := tableRef.Where(db.MetaID, "==", c.id).Documents(ctx)
	done, err := c.client.deleteByIterator(ctx, c.numShards+1, shardsIter)
	if err != nil {
		return errors.Wrap(err, "delete shards "+c.collection)
	}
	if done != true {
		return errors.Wrap(err, "delete shards not done "+c.collection)
	}
	return nil
}

// shardsCount returns shards count
//
//	count, err = shardsCount(ctx)
//
func (c *MetaFirestore) shardsCount(ctx context.Context) (int, error) {
	collectionRef := c.client.getCollectionRef(c.collection)
	shards := collectionRef.Where(db.MetaID, "==", c.id).Documents(ctx)
	defer shards.Stop()
	count := 0
	for {
		_, err := shards.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, errors.Wrapf(err, "iter next %v-%v", c.collection, c.id)
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
	var total float64
	for {
		snotshot, err := shards.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, errors.Wrapf(err, "iter next %v-%v", c.collection, c.id)
		}

		iTotal := snotshot.Data()[db.MetaN]
		shardCount, err := util.ToFloat64(iTotal)
		if err != nil {
			return 0, errors.Wrapf(err, "invalid dataType %T want float64 %v-%v", iTotal, c.collection, c.id)
		}
		total += shardCount
	}
	return total, nil
}
