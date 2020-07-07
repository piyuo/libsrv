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
	Counter `firestore:"-"`

	// N is number of shards
	//
	N int
	// nsRef point to a namespace in database
	//
	nsRef *firestore.DocumentRef

	tablename   string
	countername string
	client      *firestore.Client
	docRef      *firestore.DocumentRef
	tx          *firestore.Transaction
	// createTime is object create time
	//
	createTime time.Time

	// readTime is object read time
	//
	readTime time.Time

	// updateTime is object update time
	//
	updateTime time.Time `firestore:"-"`
}

// Shard is a single counter, which is used in a group of other shards within Counter
//
type Shard struct {
	// C is shard current count
	//
	C int
}

func (c *CounterFirestore) errorID() string {
	id := "{root}"
	if c.nsRef != nil {
		id = "{" + c.nsRef.ID + "}"
	}
	id = c.tablename + id
	if c.countername != "" {
		id += "-" + c.countername
	}
	return id
}

// create counter in transaction with a given number of shards as subcollection of specified document.
//
//	err := db.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
//		return counter.create(ctx, db.tx, docRef, counter, numShards)
//	})
//
func (c *CounterFirestore) create(ctx context.Context, tx *firestore.Transaction, docRef *firestore.DocumentRef, counter *CounterFirestore, numShards int) error {
	snapshot, err := tx.Get(docRef)
	if snapshot != nil && !snapshot.Exists() {
		counter.N = numShards
		err = tx.Set(docRef, counter)
		if err != nil {
			return errors.Wrap(err, "failed to set counter")
		}
		colRef := docRef.Collection("shards")
		// Initialize each shard with count=0
		for num := 0; num < numShards; num++ {
			shard := &Shard{C: 0}
			sharedRef := colRef.Doc(strconv.Itoa(num))
			err = tx.Set(sharedRef, shard)
			if err != nil {
				return errors.Wrapf(err, "failed init counter shared:%v ", num)
			}
		}
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "failed to get counter")
	}
	return nil
}

// Increment increments a randomly picked shard.
//
//	err = counter.Increment(ctx, 2)
//
func (c *CounterFirestore) Increment(ctx context.Context, value int) error {
	if c.N == 0 {
		return errors.New("NumShards can not be empty")
	}

	docID := strconv.Itoa(rand.Intn(c.N))
	shardRef := c.docRef.Collection("shards").Doc(docID)

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
	var total int64
	collectionRef := c.docRef.Collection("shards")
	var shards *firestore.DocumentIterator
	if c.tx != nil {
		shards = c.tx.Documents(collectionRef)
	} else {
		shards = collectionRef.Documents(ctx)
	}
	defer shards.Stop()

	for {
		doc, err := shards.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, errors.Wrap(err, "failed iterator shards documents: "+c.errorID())
		}

		vTotal := doc.Data()["C"]
		shardCount, ok := vTotal.(int64)
		if !ok {
			return 0, errors.Wrapf(err, "failed get count on shards, invalid dataType %T, want int64: "+c.errorID(), vTotal)
		}
		total += shardCount
	}
	return total, nil
}

// Delete counter and all shards.
//
//	err = counter.Delete(ctx)
//
func (c *CounterFirestore) Delete(ctx context.Context) error {
	if c.client == nil {
		return nil
	}

	shards := c.docRef.Collection("shards").Documents(ctx)
	defer shards.Stop()

	if c.tx != nil {
		for {
			doc, err := shards.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return errors.Wrap(err, "failed to iterator shards: "+c.errorID())
			}

			if err = c.tx.Delete(doc.Ref); err != nil {
				return errors.Wrap(err, "failed to delete shards in transaction: "+c.errorID())
			}

		}
		if err := c.tx.Delete(c.docRef); err != nil {
			return errors.Wrap(err, "failed to delete counter in transaction: "+c.errorID())
		}
	} else {
		shardsBatch := c.client.Batch()
		for {
			doc, err := shards.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return errors.Wrap(err, "failed to iterator shards: "+c.errorID())
			}
			shardsBatch.Delete(doc.Ref)
		}
		shardsBatch.Delete(c.docRef)

		_, err := shardsBatch.Commit(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to commit delete batch: "+c.errorID())
		}
	}

	c.nsRef = nil
	c.tablename = ""
	c.countername = ""
	c.client = nil
	c.docRef = nil
	c.tx = nil
	return nil
}

// CreateTime return object create time
//
//	id := d.CreateTime()
//
func (c *CounterFirestore) CreateTime() time.Time {
	return c.createTime
}

// SetCreateTime set object create time
//
//	id := d.SetCreateTime(time.Now())
//
func (c *CounterFirestore) SetCreateTime(t time.Time) {
	c.createTime = t
}

// ReadTime return object create time
//
//	id := d.ReadTime()
//
func (c *CounterFirestore) ReadTime() time.Time {
	return c.readTime
}

// SetReadTime set object read time
//
//	id := d.SetReadTime(time.Now())
//
func (c *CounterFirestore) SetReadTime(t time.Time) {
	c.readTime = t
}

// UpdateTime return object update time
//
//	id := d.UpdateTime()
//
func (c *CounterFirestore) UpdateTime() time.Time {
	return c.updateTime
}

// SetUpdateTime set object update time
//
//	id := d.SetUpdateTime(time.Now())
//
func (c *CounterFirestore) SetUpdateTime(t time.Time) {
	c.updateTime = t
}
