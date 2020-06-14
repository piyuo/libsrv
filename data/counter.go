package data

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// Counter is a collection of documents (shards) to realize counter with high frequency.
//
type Counter struct {
	NumShards int
	docRef    *firestore.DocumentRef
}

// Shard is a single counter, which is used in a group of other shards within Counter.
//
type Shard struct {
	Count int
}

// initCounter creates a given number of shards as
// subcollection of specified document.
func (c *Counter) init(ctx context.Context) error {
	colRef := c.docRef.Collection("shards")

	// Initialize each shard with count=0
	for num := 0; num < c.NumShards; num++ {
		shard := Shard{0}

		if _, err := colRef.Doc(strconv.Itoa(num)).Set(ctx, shard); err != nil {
			return fmt.Errorf("Set: %v", err)
		}
	}
	return nil
}

// Increment increments a randomly picked shard.
//
func (c *Counter) Increment(ctx context.Context, value int) error {
	rand.Seed(time.Now().UTC().UnixNano())
	docID := strconv.Itoa(rand.Intn(c.NumShards))
	shardRef := c.docRef.Collection("shards").Doc(docID)
	shardRef.Update(ctx, []firestore.Update{
		{Path: "Count", Value: firestore.Increment(value)},
	})
	return nil
}

// Count returns a total count across all shards.
//
func (c *Counter) Count(ctx context.Context) (int64, error) {
	var total int64
	shards := c.docRef.Collection("shards").Documents(ctx)
	for {
		doc, err := shards.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, fmt.Errorf("Next: %v", err)
		}

		vTotal := doc.Data()["Count"]
		shardCount, ok := vTotal.(int64)
		if !ok {
			return 0, fmt.Errorf("firestore: invalid dataType %T, want int64", vTotal)
		}
		total += shardCount
	}
	return total, nil
}
