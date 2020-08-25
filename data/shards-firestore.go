package data

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

// ShardsFirestore is root for counter/coder/serial, provide basic function like deleteDoc() and deleteShards()
//
type ShardsFirestore struct {

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

func (c *ShardsFirestore) errorID() string {
	return c.conn.errorID(c.tableName, c.id)
}

// assert check ctx, table name, id are valid
//
func (c *ShardsFirestore) assert(ctx context.Context) error {
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

// getRef return docRef and shardsRef
//
func (c *ShardsFirestore) getRef() (*firestore.DocumentRef, *firestore.CollectionRef) {
	docRef := c.conn.getDocRef(c.tableName, c.id)
	shardsRef := docRef.Collection("shards")
	return docRef, shardsRef
}

// deleteDoc delete shards document
//
//	err = c.deleteDoc(ctx)
//
func (c *ShardsFirestore) deleteDoc(ctx context.Context) error {
	docRef, _ := c.getRef()
	if c.conn.tx != nil {
		if err := c.conn.tx.Delete(docRef); err != nil {
			return err
		}
		return nil
	}
	if _, err := docRef.Delete(ctx); err != nil {
		return err
	}
	return nil
}

// deleteShards delete counter
//
//	err = c.deleteShards(ctx)
//
func (c *ShardsFirestore) deleteShards(ctx context.Context) error {
	if c.conn.tx != nil {
		err := c.deleteShardsTx(ctx, c.conn.tx)
		if err != nil {
			return errors.Wrap(err, "failed to delete shards in transaction: "+c.errorID())
		}
		return nil
	}
	err := c.conn.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		return c.deleteShardsTx(ctx, tx)
	})
	if err != nil {
		return errors.Wrap(err, "failed to delete shards: "+c.errorID())
	}
	return nil
}

// deleteShardsTx will delete shards collection and document
//
//	err = c.deleteShardsTx(ctx, tx)
//
func (c *ShardsFirestore) deleteShardsTx(ctx context.Context, tx *firestore.Transaction) error {
	docRef, shardsRef := c.getRef()
	shards := tx.Documents(shardsRef)
	defer shards.Stop()
	for {
		shard, err := shards.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		if err = tx.Delete(shard.Ref); err != nil {
			return err
		}

	}
	if err := tx.Delete(docRef); err != nil {
		return err
	}
	return nil
}

// shardsInfo is a debug function it return shards document and collection count
//
//	docCount,shardsCount,err = c.shardsInfo(ctx)
//
func (c *ShardsFirestore) shardsInfo(ctx context.Context) (int, int, error) {
	docRef, shardsRef := c.getRef()

	docCount := 0
	snapshot, err := docRef.Get(ctx)
	if snapshot != nil && !snapshot.Exists() {
		return 0, 0, nil
	}
	if err != nil {
		return 0, 0, errors.Wrap(err, "failed to get document: "+c.errorID())
	}
	docCount++
	shardsCount := 0
	shards := shardsRef.Documents(ctx)
	defer shards.Stop()
	for {
		_, err := shards.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, 0, err
		}
		shardsCount++
	}
	return docCount, shardsCount, nil
}
