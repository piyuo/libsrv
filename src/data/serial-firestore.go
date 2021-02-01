package data

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
)

// SerialFirestore generial serial from firestore
//
type SerialFirestore struct {
	Serial `firestore:"-"`

	MetaFirestore `firestore:"-"`

	callRX bool

	shardExist bool
}

// getRef return docRef and shardsRef
//
func (c *SerialFirestore) getRef() *firestore.DocumentRef {
	return c.conn.getDocRef(c.tableName, c.id)
}

// NumberRX return sequence number, number is unique and serial, please be aware serial can only generate one sequence per second, use it with high frequency will cause error and  must used it in transaction with NumberWX()
//
//	err = db.Transaction(ctx, func(ctx context.Context) error {
//		num, err:= serial.NumberRX()
//		So(err, ShouldBeNil)
//		So(num, ShouldEqual,1)
//		err := serial.NumberWX()
//	})
//
func (c *SerialFirestore) NumberRX() (int64, error) {
	if c.conn.tx == nil {
		return 0, errors.New("NumberRX() must run in transaction")
	}

	c.callRX = true
	snapshot, err := c.conn.tx.Get(c.getRef())
	if snapshot != nil && !snapshot.Exists() {
		c.shardExist = false
		return 1, nil
	}

	if err != nil {
		return 0, errors.Wrap(err, "failed to get serial: "+c.errorID())
	}

	idRef, err := snapshot.DataAt(MetaValue)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get value from serial: "+c.errorID())
	}
	c.shardExist = true
	id := idRef.(int64)
	return id + 1, nil
}

// NumberWX commit NumberRX
//
//	err = db.Transaction(ctx, func(ctx context.Context) error {
//		num, err:= serial.NumberRX()
//		So(err, ShouldBeNil)
//		So(num, ShouldEqual,1)
//		err := serial.NumberWX()
//	})
//
func (c *SerialFirestore) NumberWX() error {
	if c.conn.tx == nil {
		return errors.New("NumberWX() must run in transaction")
	}
	if c.callRX == false {
		return errors.New("NumberWX() need call NumberRX() first")
	}

	if c.shardExist {
		if err := c.incrementShard(c.getRef(), 1); err != nil {
			return nil
		}
	} else {
		shard := map[string]interface{}{
			MetaID:    c.id,
			MetaValue: 1,
		}
		if err := c.createShard(c.getRef(), shard); err != nil {
			return err
		}
	}
	c.callRX = false
	c.shardExist = false
	return nil
}

// Clear all shards
//
//	err = c.Clear(ctx)
//
func (c *SerialFirestore) Clear(ctx context.Context) error {
	return c.clear(ctx)
}

// ShardsCount returns shards count
//
//	count, err = coder.ShardsCount(ctx)
//
func (c *SerialFirestore) ShardsCount(ctx context.Context) (int, error) {
	return c.shardsCount(ctx)
}
