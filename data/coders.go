package data

import (
	"context"
)

// Coders is collection of code
//
type Coders struct {
	CurrentConnection Connection

	//TableName is code table name
	//
	TableName string
}

// Coder return code from database, set numshards 100 times of concurrent usage. for example if you think concurrent use is 10/seconds then set numshards to 1000 to avoid too much retention error
//
//	coders := db.Coders()
//	productCoder,err = coders.Coder("product-code",100)
//
func (c *Coders) Coder(name string, numshards int) Coder {
	return &CoderFirestore{
		ShardsFirestore: ShardsFirestore{
			conn:      c.CurrentConnection.(*ConnectionFirestore),
			tableName: c.TableName,
			id:        name,
			numShards: numshards,
		},
	}
}

// Delete code from database
//
//	coders := db.Coders()
//	err = coders.Delete(ctx, "product-code")
//
func (c *Coders) Delete(ctx context.Context, name string) error {
	shards := ShardsFirestore{
		conn:      c.CurrentConnection.(*ConnectionFirestore),
		tableName: c.TableName,
		id:        name,
		numShards: 0,
	}
	if err := shards.assert(ctx); err != nil {
		return err
	}

	return shards.deleteShards(ctx)
}
