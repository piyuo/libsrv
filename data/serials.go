package data

import (
	"context"
)

// Serials is collection of serial
//
type Serials struct {
	Connection Connection

	//TableName is serial table name
	//
	TableName string
}

// Serial return serial from database, create one if not exist, please be aware Serial can only generate 1 number per second, use serial with high frequency will cause too much retention error
//
//	serials := db.Serials()
//	productNo,err = serials.Serial("product-no")
//
func (c *Serials) Serial(name string) Serial {
	return &SerialFirestore{
		ShardsFirestore: ShardsFirestore{
			conn:      c.Connection.(*ConnectionFirestore),
			tableName: c.TableName,
			id:        name,
			numShards: 0,
		},
	}
}

// Delete serial from database
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
