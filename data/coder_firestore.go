package data

import (
	"context"
	"math/rand"
	"strconv"

	"cloud.google.com/go/firestore"
	util "github.com/piyuo/libsrv/util"
	"github.com/pkg/errors"
)

// CoderFirestore generate code from firestore
//
type CoderFirestore struct {
	CoderRef `firestore:"-"`

	ShardsFirestore `firestore:"-"`

	numberCallRX bool

	numberShardIndex int

	numberCanCreateShard bool

	numberCanIncrementShard bool
}

// CodeRX encode uint32 number into string, must used it in transaction with CodeWX()
//
//	err = db.Transaction(ctx, func(ctx context.Context) error {
//		code, err:= coder.CodeRX()
//		So(err, ShouldBeNil)
//		So(code, ShouldNotBeEmpty)
//		err := coder.CodeWX()
//	})
//
func (c *CoderFirestore) CodeRX() (string, error) {
	number, err := c.NumberRX()
	if err != nil {
		return "", err
	}
	return util.SerialID32(uint32(number)), nil
}

// CodeWX commit CodeRX()
//
//	err = db.Transaction(ctx, func(ctx context.Context) error {
//		code, err:= coder.Code16RX()
//		So(err, ShouldBeNil)
//		So(code, ShouldNotBeEmpty)
//		err := coder.Code16WX()
//	})
//
func (c *CoderFirestore) CodeWX() error {
	return c.NumberWX()
}

// Code16RX encode uint16 number into string, must used it in transaction with CodeWX()
//
//	err = db.Transaction(ctx, func(ctx context.Context) error {
//		code, err:= coder.Code16RX()
//		So(err, ShouldBeNil)
//		So(code, ShouldNotBeEmpty)
//		err := coder.Code16WX()
//	})
//
func (c *CoderFirestore) Code16RX() (string, error) {
	number, err := c.NumberRX()
	if err != nil {
		return "", err
	}
	return util.SerialID16(uint16(number)), nil
}

// Code16WX commit Code16RX()
//
//	err = db.Transaction(ctx, func(ctx context.Context) error {
//		code, err:= coder.Code16RX()
//		So(err, ShouldBeNil)
//		So(code, ShouldNotBeEmpty)
//		err := coder.Code16WX()
//	})
//
func (c *CoderFirestore) Code16WX() error {
	return c.NumberWX()
}

// Code64RX encode uint32 number into string, must used it in transaction with Code64WX()
//
//	err = db.Transaction(ctx, func(ctx context.Context) error {
//		code, err:= coder.Code64RX()
//		So(err, ShouldBeNil)
//		So(code, ShouldNotBeEmpty)
//		err := coder.Code64WX()
//	})
//
func (c *CoderFirestore) Code64RX() (string, error) {
	number, err := c.NumberRX()
	if err != nil {
		return "", err
	}
	return util.SerialID64(uint64(number)), nil
}

// Code64WX commit with Code64RX()
//
//	err = db.Transaction(ctx, func(ctx context.Context) error {
//		code, err:= coder.Code64RX()
//		So(err, ShouldBeNil)
//		So(code, ShouldNotBeEmpty)
//		err := coder.Code64WX()
//	})
//
func (c *CoderFirestore) Code64WX() error {
	return c.NumberWX()
}

// NumberRX prepare return unique but not serial number, must used it in transaction with NumberWX()
//
//	err = db.Transaction(ctx, func(ctx context.Context) error {
//		num, err:= coder.NumberRX()
//		So(err, ShouldBeNil)
//		So(num > 0, ShouldBeTrue)
//		err := code.NumberWX()
//	})
//
func (c *CoderFirestore) NumberRX() (int64, error) {
	if c.conn.tx == nil {
		return 0, errors.New("this function must run in transaction")
	}

	_, shardsRef := c.getRef()
	c.numberCallRX = true
	c.numberCanCreateShard = false
	c.numberCanIncrementShard = false
	c.numberShardIndex = rand.Intn(c.numShards)
	shardID := strconv.Itoa(c.numberShardIndex)

	snapshot, err := c.conn.tx.Get(shardsRef.Doc(shardID))
	if snapshot != nil && !snapshot.Exists() {
		c.numberCanCreateShard = true
		return int64(1*c.numShards + c.numberShardIndex), nil
	}

	if err != nil {
		return 0, errors.Wrap(err, "failed to get number: "+c.errorID())
	}
	idRef, err := snapshot.DataAt("N")
	if err != nil {
		return 0, errors.Wrap(err, "failed to get value from number: "+c.errorID())
	}
	id := idRef.(int64)
	c.numberCanIncrementShard = true
	value := (id+1)*int64(c.numShards) + int64(c.numberShardIndex)
	return value, nil
}

// NumberWX commit NumberRX()
//
//	err = db.Transaction(ctx, func(ctx context.Context) error {
//		num, err:= coder.NumberRX()
//		So(err, ShouldBeNil)
//		So(num > 0, ShouldBeTrue)
//		err := code.NumberWX()
//	})
//
func (c *CoderFirestore) NumberWX() error {
	if c.conn.tx == nil {
		return errors.New("This function must run in transaction")
	}
	if c.numberCallRX == false {
		return errors.New("WX() function need call NumberRX() first")
	}

	docRef, shardsRef := c.getRef()
	shardID := strconv.Itoa(c.numberShardIndex)
	shardRef := shardsRef.Doc(shardID)
	if c.numberCanCreateShard {
		// create shards document
		err := c.conn.tx.Set(docRef, &struct{}{}) //put empty struct
		if err != nil {
			return errors.Wrap(err, "failed to create shards document: "+c.errorID())
		}

		// create shard
		err = c.conn.tx.Set(shardRef, map[string]interface{}{"N": 1}, firestore.MergeAll)
		if err != nil {
			return errors.Wrap(err, "failed to create shard: "+c.errorID())
		}
	}

	if c.numberCanIncrementShard {
		err := c.conn.tx.Update(shardRef, []firestore.Update{
			{Path: "N", Value: firestore.Increment(1)},
		})
		if err != nil {
			return errors.Wrap(err, "failed to increment shard: "+c.errorID())
		}
	}
	c.numberCallRX = false
	c.numberCanCreateShard = false
	c.numberCanIncrementShard = false
	c.numberShardIndex = -1
	return nil
}

// Reset reset code
//
//	err = db.Transaction(ctx, func(ctx context.Context) error {
//		err:= coder.Reset(ctx)
//	})
//
func (c *CoderFirestore) Reset(ctx context.Context) error {
	if err := c.assert(ctx); err != nil {
		return err
	}
	return c.deleteShards(ctx)
}
