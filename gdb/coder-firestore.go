package gdb

import (
	"context"
	"math/rand"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/piyuo/libsrv/db"
	identifier "github.com/piyuo/libsrv/identifier"
	"github.com/piyuo/libsrv/log"
	"github.com/pkg/errors"
)

// CoderFirestore generate code from firestore
//
type CoderFirestore struct {
	db.Coder `firestore:"-"`

	MetaFirestore `firestore:"-"`

	callRX bool

	shardPick int

	shardExist bool
}

// getPickedRef return picked all period ref
//
func (c *CoderFirestore) getPickedRef() *firestore.DocumentRef {
	return c.client.getDocRef(c.collection, c.id+"_"+strconv.Itoa(c.shardPick))
}

// pickShard random pick a shard, return shardIndex, isShardExist, error
//
func (c *CoderFirestore) pickShard(ctx context.Context, transaction db.Transaction) (bool, int64, error) {
	tx := transaction.(*TransactionFirestore)
	snapshot, err := tx.snapshot(ctx, c.getPickedRef())
	if err != nil {
		return false, 0, errors.Wrapf(err, "get coder snapshot %v-%v", c.collection, c.id)
	}
	if snapshot == nil {
		// value format is incrementValue+shardIndex, e.g. 12 , 1= increment value, 2=shard index
		value := int64(c.numShards + c.shardPick)
		return false, value, nil
	}

	idRef, err := snapshot.DataAt(db.MetaN)
	if err != nil {
		return false, 0, errors.Wrapf(err, "get data at snapshot %v-%v", c.collection, c.id)
	}
	id := idRef.(int64)
	value := (id+1)*int64(c.numShards) + int64(c.shardPick)
	return true, value, nil
}

// CodeRX encode uint32 number into string, must used it in transaction with CodeWX()
//
//	err := Transaction(ctx, func(ctx context.Context,tx db.Transaction) error {
//		code, err:= coder.CodeRX(ctx,tx)
//		err := coder.CodeWX(ctx,tx)
//	})
//
func (c *CoderFirestore) CodeRX(ctx context.Context, transaction db.Transaction) (string, error) {
	number, err := c.NumberRX(ctx, transaction)
	if err != nil {
		return "", err
	}
	return identifier.SerialID32(uint32(number)), nil
}

// CodeWX commit CodeRX()
//
//	err := Transaction(ctx, func(ctx context.Context,tx db.Transaction) error {
//		code, err:= coder.CodeRX(ctx,tx)
//		err := coder.CodeWX(ctx,tx)
//	})
//
func (c *CoderFirestore) CodeWX(ctx context.Context, transaction db.Transaction) error {
	return c.NumberWX(ctx, transaction)
}

// Code16RX encode uint16 number into string, must used it in transaction with CodeWX()
//
//	err := Transaction(ctx, func(ctx context.Context,tx db.Transaction) error {
//		code, err:= coder.Code16RX(ctx,tx)
//		err := coder.Code16WX(ctx,tx)
//	})
//
func (c *CoderFirestore) Code16RX(ctx context.Context, transaction db.Transaction) (string, error) {
	number, err := c.NumberRX(ctx, transaction)
	if err != nil {
		return "", err
	}
	return identifier.SerialID16(uint16(number)), nil
}

// Code16WX commit Code16RX()
//
//	err := Transaction(ctx, func(ctx context.Context,tx db.Transaction) error {
//		code, err:= coder.Code16RX(ctx,tx)
//		err := coder.Code16WX(ctx,tx)
//	})
//
func (c *CoderFirestore) Code16WX(ctx context.Context, transaction db.Transaction) error {
	return c.NumberWX(ctx, transaction)
}

// Code64RX encode uint32 number into string, must used it in transaction with Code64WX()
//
//	err := Transaction(ctx, func(ctx context.Context,tx db.Transaction) error {
//		code, err:= coder.Code64RX(ctx,tx)
//		err := coder.Code64WX(ctx,tx)
//	})
//
func (c *CoderFirestore) Code64RX(ctx context.Context, transaction db.Transaction) (string, error) {
	number, err := c.NumberRX(ctx, transaction)
	if err != nil {
		return "", err
	}
	return identifier.SerialID64(uint64(number)), nil
}

// Code64WX commit with Code64RX()
//
//	err := Transaction(ctx, func(ctx context.Context,tx db.Transaction) error {
//		code, err:= coder.Code64RX(ctx,tx)
//		err := coder.Code64WX(ctx,tx)
//	})
//
func (c *CoderFirestore) Code64WX(ctx context.Context, transaction db.Transaction) error {
	return c.NumberWX(ctx, transaction)
}

// NumberRX prepare return unique but not serial number, must used it in transaction with NumberWX()
//
//	err := Transaction(ctx, func(ctx context.Context,tx db.Transaction) error {
//		num, err:= coder.NumberRX(ctx,tx)
//		err := coder.NumberWX(ctx,tx)
//	})
//
func (c *CoderFirestore) NumberRX(ctx context.Context, transaction db.Transaction) (int64, error) {
	c.callRX = true
	c.shardPick = rand.Intn(c.numShards) //random pick a shard
	log.Debug(ctx, "coder pick %v from %v shards", c.shardPick, c.numShards)

	exist, value, err := c.pickShard(ctx, transaction)
	if err != nil {
		return 0, err
	}
	c.shardExist = exist
	return value, nil
}

// NumberWX commit NumberRX()
//
//	err := Transaction(ctx, func(ctx context.Context,tx db.Transaction) error {
//		num, err:= coder.NumberRX(ctx,tx)
//		err := coder.NumberWX(ctx,tx)
//	})
//
func (c *CoderFirestore) NumberWX(ctx context.Context, transaction db.Transaction) error {
	tx := transaction.(*TransactionFirestore)
	if !c.callRX {
		return errors.New("must call RX first")
	}

	if c.shardExist {
		if err := tx.incrementShard(c.getPickedRef(), 1); err != nil {
			return errors.Wrap(err, "inc shard")
		}
	} else {
		shard := map[string]interface{}{
			db.MetaID: c.id,
			db.MetaN:  1,
		}
		if err := tx.createShard(c.getPickedRef(), shard); err != nil {
			return errors.Wrap(err, "new shard")
		}
	}
	c.callRX = false
	c.shardExist = false
	c.shardPick = -1
	return nil
}

// Delete delete coder
//
//	err = Delete(ctx)
//
func (c *CoderFirestore) Delete(ctx context.Context) error {
	return c.deleteShards(ctx)
}

// ShardsCount returns shards count
//
//	count, err = ShardsCount(ctx)
//
func (c *CoderFirestore) ShardsCount(ctx context.Context) (int, error) {
	return c.shardsCount(ctx)
}
