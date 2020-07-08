package data

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

// CounterFirestore implement Counter
//
type CounterFirestore struct {
	CounterRef `firestore:"-"`

	// client is firestore client
	//
	client *firestore.Client

	// tx is transaction
	//
	tx *firestore.Transaction

	// shardsRef point to a namespace in database
	//
	shardsRef *firestore.CollectionRef

	// N is number of shards
	//
	N int

	// NameSpace is counter table name
	//
	NameSpace string `firestore:"-"`

	// TableName is counter table name
	//
	TableName string `firestore:"-"`

	// CounterName is counter name
	//
	CounterName string `firestore:"-"`

	// CreateTime is object create time, this is readonly field
	//
	CreateTime time.Time `firestore:"-"`

	// ReadTime is object read time, this is readonly field
	//
	ReadTime time.Time `firestore:"-"`

	// UpdateTime is object update time, this is readonly field
	//
	UpdateTime time.Time `firestore:"-"`
}

// Shard is a single counter, which is used in a group of other shards within Counter
//
type Shard struct {
	// C is shard current count
	//
	C int
}

// errorID return error id
//
func (c *CounterFirestore) errorID() string {
	id := "{root}"
	if c.NameSpace != "" {
		id = "{" + c.NameSpace + "}"
	}
	id = c.TableName + id
	if c.CounterName != "" {
		id += "-" + c.CounterName
	}
	return id
}

// Increment increments a randomly picked shard.
//
//	err = counter.Increment(ctx, 2)
//
func (c *CounterFirestore) Increment(ctx context.Context, value int) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	docID := strconv.Itoa(rand.Intn(c.N))
	shardRef := c.shardsRef.Doc(docID)
	var err error
	if c.tx != nil {
		err = c.tx.Update(shardRef, []firestore.Update{
			{Path: "C", Value: firestore.Increment(value)},
		})

	} else {
		_, err = shardRef.Update(ctx, []firestore.Update{
			{Path: "C", Value: firestore.Increment(value)},
		})
	}
	if err != nil {
		return errors.Wrap(err, "failed to increment counter: "+c.errorID())
	}
	return nil
}

// Count returns a total count across all shards.
//
//	count, err = counter.Count(ctx)
//
func (c *CounterFirestore) Count(ctx context.Context) (int64, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	var total int64
	var shards *firestore.DocumentIterator
	if c.tx != nil {
		shards = c.tx.Documents(c.shardsRef)
	} else {
		shards = c.shardsRef.Documents(ctx)
	}
	defer shards.Stop()

	for {
		snotshot, err := shards.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, errors.Wrap(err, "failed iterator shards documents: "+c.errorID())
		}

		vTotal := snotshot.Data()["C"]
		shardCount, ok := vTotal.(int64)
		if !ok {
			return 0, errors.Wrapf(err, "failed get count on shards, invalid dataType %T, want int64: "+c.errorID(), vTotal)
		}
		total += shardCount
	}
	return total, nil
}

// GetCreateTime return object create time
//
//	id := d.CreateTime()
//
func (c *CounterFirestore) GetCreateTime() time.Time {
	return c.CreateTime
}

// GetReadTime return object create time
//
//	id := d.ReadTime()
//
func (c *CounterFirestore) GetReadTime() time.Time {
	return c.ReadTime
}

// GetUpdateTime return object update time
//
//	id := d.UpdateTime()
//
func (c *CounterFirestore) GetUpdateTime() time.Time {
	return c.UpdateTime
}
