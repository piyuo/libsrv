package data

import (
	"context"
	"math/rand"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

// CounterFirestore implement Counter
//
type CounterFirestore struct {
	CounterRef `firestore:"-"`

	// connection is current firestore connection
	//
	connection *ConnectionFirestore

	// numShards is number of shards
	//
	numShards int

	// cableName is counter table name
	//
	tableName string `firestore:"-"`

	// counterName is counter name
	//
	counterName string `firestore:"-"`
}

// Shard is a single counter, which is used in a group of other shards within Counter
//
type Shard struct {
	// C is shard current count
	//
	C int
}

func (c *CounterFirestore) errorID() string {
	return c.connection.errorID(c.tableName, c.counterName)
}

// getRef return docRef and shardsRef to do counter database operation
//
func (c *CounterFirestore) getRef() (*firestore.DocumentRef, *firestore.CollectionRef) {
	docRef := c.connection.getDocRef(c.tableName, c.counterName)
	shardsRef := docRef.Collection("shards")
	return docRef, shardsRef
}

// Increment increments a randomly picked shard.
//
//	err = counter.Increment(ctx, 1)
//
func (c *CounterFirestore) Increment(ctx context.Context, value int) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if c.tableName == "" {
		return errors.New("table name can not be empty")
	}
	if c.counterName == "" {
		return errors.New("counter name can not be empty")
	}

	_, shardsRef := c.getRef()
	shardID := strconv.Itoa(rand.Intn(c.numShards))
	shardRef := shardsRef.Doc(shardID)
	var err error
	if c.connection.tx != nil {
		err = c.connection.tx.Update(shardRef, []firestore.Update{
			{Path: "C", Value: firestore.Increment(value)},
		})

		if err != nil {
			return errors.Wrap(err, "failed to update shard in transaction, you may need CreateShards() first: "+c.errorID()+"["+shardID+"]")
		}
		return nil
	}

	_, err = shardRef.Update(ctx, []firestore.Update{
		{Path: "C", Value: firestore.Increment(value)},
	})

	if err != nil {
		return errors.Wrap(err, "failed to update shard, you may need CreateShards() first: "+c.errorID()+"["+shardID+"]")
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
	if c.tableName == "" {
		return 0, errors.New("table name can not be empty")
	}
	if c.counterName == "" {
		return 0, errors.New("counter name can not be empty")
	}
	_, shardsRef := c.getRef()
	var total int64
	var shards *firestore.DocumentIterator
	if c.connection.tx != nil {
		shards = c.connection.tx.Documents(shardsRef)
	} else {
		shards = shardsRef.Documents(ctx)
	}
	defer shards.Stop()

	shardCount := 0
	for {
		snotshot, err := shards.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, errors.Wrap(err, "failed to iterator shards: "+c.errorID())
		}

		shardCount++
		vTotal := snotshot.Data()["C"]
		shardCount, ok := vTotal.(int64)
		if !ok {
			return 0, errors.Wrapf(err, "failed to get count on shards, invalid dataType %T, want int64: "+c.errorID(), vTotal)
		}
		total += shardCount
	}
	if shardCount == 0 {
		return 0, errors.New("failed to get any shards, you may need CreateShards() first: " + c.errorID())
	}
	return total, nil
}

// CreateShards create counter and all shards, it is safe to create shards as many time as you want, normally we recreate shards when we need more shards
//
//	err = counter.CreateShards(ctx)
//
func (c *CounterFirestore) CreateShards(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	shard := &Shard{C: 0}
	docRef, shardsRef := c.getRef()

	if c.connection.tx != nil {
		err := c.connection.tx.Set(docRef, c)
		if err != nil {
			return errors.Wrapf(err, "failed create counter in transaction: "+c.errorID())
		}

		for num := 0; num < c.numShards; num++ {
			sharedRef := shardsRef.Doc(strconv.Itoa(num))
			err := c.connection.tx.Set(sharedRef, shard)
			if err != nil {
				return errors.Wrapf(err, "failed create shared:%v", num)
			}
		}
		return nil
	}

	_, err := docRef.Set(ctx, c)
	if err != nil {
		return errors.Wrapf(err, "failed create counter: "+c.errorID())
	}

	for num := 0; num < c.numShards; num++ {
		sharedRef := shardsRef.Doc(strconv.Itoa(num))
		_, err := sharedRef.Set(ctx, shard)
		if err != nil {
			return errors.Wrapf(err, "failed create shared:%v", num)
		}
	}
	return nil
}
