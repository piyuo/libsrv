package data

import (
	"context"
)

// Coders is collection of code
//
type Coders struct {
	Connection ConnectionRef

	//TableName is code table name
	//
	TableName string
}

// Coder return code from database, the numShards must be multiple by 10
//
//	coders := db.Coders()
//	productCoder,err = coders.Coder("product-code",10)
//
func (c *Coders) Coder(name string, numshards int) CoderRef {
	return &CoderFirestore{
		ShardsFirestore: ShardsFirestore{
			conn:      c.Connection.(*ConnectionFirestore),
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
