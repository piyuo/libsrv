package data

import (
	"context"
	"math/rand"
	"strconv"

	"cloud.google.com/go/firestore"
	util "github.com/piyuo/libsrv/util"
	"github.com/pkg/errors"
)

// CodeFirestore generate code from firestore
//
type CodeFirestore struct {
	CodeRef `firestore:"-"`

	ShardsFirestore `firestore:"-"`
}

// Code encode uint32 number into string, please be aware serial can only generate one number per second
//
//	code, err := code.Code(ctx)
//	So(c, ShouldBeEmpty)
//
func (c *CodeFirestore) Code(ctx context.Context) (string, error) {
	number, err := c.Number(ctx)
	if err != nil {
		return "", err
	}
	return util.SerialID32(uint32(number)), nil
}

// Code16 encode uint16 number into string, please be aware serial can only generate one number per second
//
//	c, err := code.Code16(ctx)
//	So(c, ShouldBeEmpty)
//
func (c *CodeFirestore) Code16(ctx context.Context) (string, error) {
	number, err := c.Number(ctx)
	if err != nil {
		return "", err
	}
	return util.SerialID16(uint16(number)), nil
}

// Code64 encode uint64 serial number to string, please be aware serial can only generate one number per second
//
//	c, err := code.Code64(ctx)
//	So(c, ShouldBeEmpty)
//
func (c *CodeFirestore) Code64(ctx context.Context) (string, error) {
	number, err := c.Number(ctx)
	if err != nil {
		return "", err
	}
	return util.SerialID64(uint64(number)), nil
}

// Number return code number, number is unique but not serial
//
//	n, err := code.Number(ctx)
//
func (c *CodeFirestore) Number(ctx context.Context) (int64, error) {
	if err := c.assert(ctx); err != nil {
		return 0, err
	}

	if c.conn.tx != nil {
		return c.getTx(ctx, c.conn.tx)
	}

	var id int64
	var err error
	err = c.conn.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		id, err = c.getTx(ctx, tx)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, errors.Wrap(err, "failed to commit serial transaction: "+c.errorID())
	}
	return id, nil
}

// getTx generate code number in transaction, return number is not serial
//
//	num, err := s.getTx(ctx)
//	So(num, ShouldEqual, 1)
//
func (c *CodeFirestore) getTx(ctx context.Context, tx *firestore.Transaction) (int64, error) {
	docRef, shardsRef := c.getRef()
	shardIndex := rand.Intn(c.numShards)
	shardID := strconv.Itoa(shardIndex)
	shardRef := shardsRef.Doc(shardID)

	snapshot, err := tx.Get(shardRef)
	if snapshot != nil && !snapshot.Exists() {

		if err := c.ensureShardsDocumentTx(tx, docRef); err != nil {
			return 0, errors.Wrap(err, "failed to init code number: "+c.errorID())
		}

		err = tx.Set(shardRef, map[string]interface{}{"N": 1}, firestore.MergeAll)
		if err != nil {
			return 0, errors.Wrap(err, "failed to init code number: "+c.errorID())
		}
		return int64(1*c.numShards + shardIndex), nil
	}

	if err != nil {
		return 0, errors.Wrap(err, "failed to get code number: "+c.errorID())
	}
	idRef, err := snapshot.DataAt("N")
	if err != nil {
		return 0, errors.Wrap(err, "failed to get value from code number: "+c.errorID())
	}
	id := idRef.(int64)
	err = tx.Update(shardRef, []firestore.Update{
		{Path: "N", Value: firestore.Increment(1)},
	})
	if err != nil {
		return 0, errors.Wrap(err, "failed to increment code number: "+c.errorID())
	}
	return (id+1)*int64(c.numShards) + int64(shardIndex), nil
}

// CreateShards create shards document and collection, it is safe to create shards as many time as you want, normally we recreate shards when we need more shards
//
//	err = code.CreateShards(ctx)
//
func (c *CodeFirestore) CreateShards(ctx context.Context) error {
	return c.createShards(ctx)
}
