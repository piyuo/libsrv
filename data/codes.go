package data

import (
	"context"
)

// Codes is collection of code
//
type Codes struct {
	Connection ConnectionRef

	//TableName is code table name
	//
	TableName string
}

// Code return code from database, the numShards must be multiple by 10
//
//	codes := db.Codes()
//	productCode,err = codes.Code("product-code",10)
//
func (c *Codes) Code(name string, numshards int) CodeRef {
	return &CodeFirestore{
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
//	codes := db.Codes()
//	err = codes.Delete(ctx, "product-code")
//
func (c *Codes) Delete(ctx context.Context, name string) error {
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
