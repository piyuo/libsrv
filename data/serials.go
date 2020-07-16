package data

import (
	"context"
)

// Serials is collection of serial
//
type Serials struct {
	Connection ConnectionRef

	//TableName is serial table name
	//
	TableName string
}

// Serial return serial from database, create one if not exist
//
//	serials := db.Serials()
//	productNo,err = serials.Serial("product-no")
//
func (c *Serials) Serial(name string) SerialRef {
	return &SerialFirestore{
		ShardsFirestore: ShardsFirestore{
			conn:      c.Connection.(*ConnectionFirestore),
			tableName: c.TableName,
			id:        name,
			numShards: 0,
		},
	}
}

// Delete code from database
//
//	serials := db.Serials()
//	err = serials.Delete(ctx, "product-no")
//
func (c *Serials) Delete(ctx context.Context, name string) error {

	shards := ShardsFirestore{
		conn:      c.Connection.(*ConnectionFirestore),
		tableName: c.TableName,
		id:        name,
		numShards: 0,
	}
	if err := shards.assert(ctx); err != nil {
		return err
	}
	return shards.deleteDoc(ctx)
}
