package data

import (
	"context"
)

// Counters is collection of counter
//
type Counters struct {
	Connection ConnectionRef

	//TableName is counter table name
	//
	TableName string
}

// Counter return counter from database, create one if not exist
//
//	counters := db.Counters()
//	orderCountCounter,err = counters.Counter("order-count",10)
//
func (c *Counters) Counter(name string, numshards int) CounterRef {

	if numshards <= 0 {
		numshards = 10
	}

	return &CounterFirestore{
		ShardsFirestore: ShardsFirestore{
			conn:      c.Connection.(*ConnectionFirestore),
			tableName: c.TableName,
			id:        name,
			numShards: numshards,
		},
	}
}

// Delete counter from database
//
//	counters := db.Counters()
//	err = counters.Delete(ctx, "myCounter")
//
func (c *Counters) Delete(ctx context.Context, name string) error {
	shards := ShardsFirestore{
		conn:      c.Connection.(*ConnectionFirestore),
		tableName: c.TableName,
		id:        name,
		numShards: 0,
	}
	if err := shards.assert(ctx); err != nil {
		return err
	}

	return shards.deleteShards(ctx)
}
