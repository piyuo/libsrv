package data

import (
	"context"
	"math/rand"
	"strconv"

	"cloud.google.com/go/firestore"
	identifier "github.com/piyuo/libsrv/identifier"
	"github.com/pkg/errors"
)

// CoderFirestore generate code from firestore
//
type CoderFirestore struct {
	Coder `firestore:"-"`

	MetaFirestore `firestore:"-"`

	callRX bool

	shardPick int

	shardExist bool
}

// CodeRX encode uint32 number into string, must used it in transaction with CodeWX()
//
//	err = db.Transaction(ctx, func(ctx context.Context) error {
//		code, err:= coder.CodeRX(ctx)
//		So(err, ShouldBeNil)
//		So(code, ShouldNotBeEmpty)
//		err := coder.CodeWX(ctx)
//	})
//
func (c *CoderFirestore) CodeRX(ctx context.Context) (string, error) {
	number, err := c.NumberRX(ctx)
	if err != nil {
		return "", err
	}
	return identifier.SerialID32(uint32(number)), nil
}

// CodeWX commit CodeRX()
//
//	err = db.Transaction(ctx, func(ctx context.Context) error {
//		code, err:= coder.Code16RX(ctx)
//		So(err, ShouldBeNil)
//		So(code, ShouldNotBeEmpty)
//		err := coder.Code16WX(ctx)
//	})
//
func (c *CoderFirestore) CodeWX(ctx context.Context) error {
	return c.NumberWX(ctx)
}

// Code16RX encode uint16 number into string, must used it in transaction with CodeWX()
//
//	err = db.Transaction(ctx, func(ctx context.Context) error {
//		code, err:= coder.Code16RX(ctx)
//		So(err, ShouldBeNil)
//		So(code, ShouldNotBeEmpty)
//		err := coder.Code16WX(ctx)
//	})
//
func (c *CoderFirestore) Code16RX(ctx context.Context) (string, error) {
	number, err := c.NumberRX(ctx)
	if err != nil {
		return "", err
	}
	return identifier.SerialID16(uint16(number)), nil
}

// Code16WX commit Code16RX()
//
//	err = db.Transaction(ctx, func(ctx context.Context) error {
//		code, err:= coder.Code16RX(ctx)
//		So(err, ShouldBeNil)
//		So(code, ShouldNotBeEmpty)
//		err := coder.Code16WX(ctx)
//	})
//
func (c *CoderFirestore) Code16WX(ctx context.Context) error {
	return c.NumberWX(ctx)
}

// Code64RX encode uint32 number into string, must used it in transaction with Code64WX()
//
//	err = db.Transaction(ctx, func(ctx context.Context) error {
//		code, err:= coder.Code64RX(ctx)
//		So(err, ShouldBeNil)
//		So(code, ShouldNotBeEmpty)
//		err := coder.Code64WX(cts)
//	})
//
func (c *CoderFirestore) Code64RX(ctx context.Context) (string, error) {
	number, err := c.NumberRX(ctx)
	if err != nil {
		return "", err
	}
	return identifier.SerialID64(uint64(number)), nil
}

// Code64WX commit with Code64RX()
//
//	err = db.Transaction(ctx, func(ctx context.Context) error {
//		code, err:= coder.Code64RX(ctx)
//		So(err, ShouldBeNil)
//		So(code, ShouldNotBeEmpty)
//		err := coder.Code64WX(ctx)
//	})
//
func (c *CoderFirestore) Code64WX(ctx context.Context) error {
	return c.NumberWX(ctx)
}

// NumberRX prepare return unique but not serial number, must used it in transaction with NumberWX()
//
//	err = db.Transaction(ctx, func(ctx context.Context) error {
//		num, err:= coder.NumberRX(ctx)
//		So(err, ShouldBeNil)
//		So(num > 0, ShouldBeTrue)
//		err := code.NumberWX(ctx)
//	})
//
func (c *CoderFirestore) NumberRX(ctx context.Context) (int64, error) {
	if c.conn.tx == nil {
		return 0, errors.New("NumberRX() must run in transaction")
	}

	c.callRX = true
	c.shardPick = rand.Intn(c.numShards) //random pick a shard
	//	fmt.Printf("coder pick:" + strconv.Itoa(c.shardPick) + "\n")

	exist, value, err := c.pickShard(ctx)
	if err != nil {
		return 0, err
	}
	c.shardExist = exist
	return value, nil
}

// getPickedRef return picked all period ref
//
func (c *CoderFirestore) getPickedRef() *firestore.DocumentRef {
	return c.conn.getDocRef(c.tableName, c.id+"."+strconv.Itoa(c.shardPick))
}

// pickShard random pick a shard, return shardIndex, isShardExist, error
//
func (c *CoderFirestore) pickShard(ctx context.Context) (bool, int64, error) {
	snapshot, err := c.conn.tx.Get(c.getPickedRef())
	if snapshot != nil && !snapshot.Exists() {
		// value format is incrementValue+shardIndex, e.g. 12 , 1= increment value, 2=shard index
		value := int64(c.numShards + c.shardPick)
		return false, value, nil
	}

	if err != nil {
		return false, 0, err
	}

	idRef, err := snapshot.DataAt(MetaValue)
	if err != nil {
		return false, 0, errors.Wrap(err, "failed to get value from number: "+c.errorID())
	}
	id := idRef.(int64)
	value := (id+1)*int64(c.numShards) + int64(c.shardPick)
	return true, value, nil
}

// NumberWX commit NumberRX()
//
//	err = db.Transaction(ctx, func(ctx context.Context) error {
//		num, err:= coder.NumberRX(ctx)
//		So(err, ShouldBeNil)
//		So(num > 0, ShouldBeTrue)
//		err := code.NumberWX(ctx)
//	})
//
func (c *CoderFirestore) NumberWX(ctx context.Context) error {
	if c.conn.tx == nil {
		return errors.New("NumberWX() must run in transaction")
	}
	if c.callRX == false {
		return errors.New("NumberWX() need call NumberRX() first")
	}

	if c.shardExist {
		if err := c.incrementShard(c.getPickedRef(), 1); err != nil {
			return err
		}
	} else {
		shard := map[string]interface{}{
			MetaID:    c.id,
			MetaValue: 1,
		}
		if err := c.createShard(c.getPickedRef(), shard); err != nil {
			return err
		}
	}
	c.callRX = false
	c.shardExist = false
	c.shardPick = -1
	return nil
}

// Clear all shards
//
//	err = c.Clear(ctx)
//
func (c *CoderFirestore) Clear(ctx context.Context) error {
	return c.clear(ctx)
}

// ShardsCount returns shards count
//
//	count, err = coder.ShardsCount(ctx)
//
func (c *CoderFirestore) ShardsCount(ctx context.Context) (int, error) {
	return c.shardsCount(ctx)
}
