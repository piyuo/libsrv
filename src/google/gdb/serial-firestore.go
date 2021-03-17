package gdb

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/piyuo/libsrv/src/data"
	"github.com/piyuo/libsrv/src/db"
	"github.com/pkg/errors"
)

// SerialFirestore generial serial from firestore
//
type SerialFirestore struct {
	data.Serial `firestore:"-"`

	MetaFirestore `firestore:"-"`

	callRX bool

	shardExist bool
}

// getRef return docRef and shardsRef
//
func (c *SerialFirestore) getRef() *firestore.DocumentRef {
	return c.client.getDocRef(c.collection, c.id)
}

// NumberRX return sequence number, number is unique and serial, please be aware serial can only generate one sequence per second, use it with high frequency will cause error and  must used it in transaction with NumberWX()
//
//	err := Transaction(ctx, func(ctx context.Context,tx db.Transaction) error {
//		num, err:= serial.NumberRX(ctx,tx)
//		err := serial.NumberWX(ctx,tx)
//	})
//
func (c *SerialFirestore) NumberRX(ctx context.Context, transaction db.Transaction) (int64, error) {
	tx := transaction.(*TransactionFirestore)
	c.callRX = true
	snapshot, err := tx.snapshot(ctx, c.getRef())
	if err != nil {
		return 0, errors.Wrapf(err, "get serial snapshot %v-%v", c.collection, c.id)
	}

	if snapshot == nil {
		c.shardExist = false
		return 1, nil
	}

	idRef, err := snapshot.DataAt(data.MetaValue)
	if err != nil {
		return 0, errors.Wrapf(err, "get data at snapshot %v-%v", c.collection, c.id)
	}
	c.shardExist = true
	id := idRef.(int64)
	return id + 1, nil
}

// NumberWX commit NumberRX
//
//	err := Transaction(ctx, func(ctx context.Context,tx db.Transaction) error {
//		num, err:= serial.NumberRX(ctx,tx)
//		err := serial.NumberWX(ctx,tx)
//	})
//
func (c *SerialFirestore) NumberWX(ctx context.Context, transaction db.Transaction) error {
	if c.callRX == false {
		return errors.New("must call RX first")
	}

	tx := transaction.(*TransactionFirestore)
	if c.shardExist {
		if err := tx.incrementShard(c.getRef(), 1); err != nil {
			return nil
		}
	} else {
		shard := map[string]interface{}{
			data.MetaID:    c.id,
			data.MetaValue: 1,
		}
		if err := tx.createShard(c.getRef(), shard); err != nil {
			return err
		}
	}
	c.callRX = false
	c.shardExist = false
	return nil
}

// Clear delete all shard in collection. delete max doc count. return true if collection is cleared
//
//	err = Clear(ctx,100)
//
func (c *SerialFirestore) Clear(ctx context.Context, max int) (bool, error) {
	return c.clear(ctx, max)
}

// ShardsCount returns shards count
//
//	count, err = ShardsCount(ctx)
//
func (c *SerialFirestore) ShardsCount(ctx context.Context) (int, error) {
	return c.shardsCount(ctx)
}
