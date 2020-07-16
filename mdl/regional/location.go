package rmdl

import (
	data "github.com/piyuo/libsrv/data"
)

// Location represent single location
//
type Location struct {
	data.Object `firestore:"-"`
}

// LocationTable return location table
//
//	counter := db.LocationTable()
//
func (db *DB) LocationTable() *data.Table {

	return &data.Table{
		Connection: db.Connection,
		TableName:  "location",
		Factory: func() data.ObjectRef {
			return &Location{}
		},
	}
}

// LocationTotal return total location count
//
//	id := d.LocationTotal(ctx)
//
func (c *Counters) LocationTotal() data.CounterRef {
	return c.Counter("LocationTotal", 10)
}

// LocationID return location id coder
//
//	coder := d.LocationID()
//
func (c *Coders) LocationID() data.CoderRef {
	return c.Coder("LocationID", 100)
}
