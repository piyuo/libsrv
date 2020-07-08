package rmdl

import (
	"context"

	data "github.com/piyuo/libsrv/data"
)

// Location represent single location
//
type Location struct {
	data.Object `firestore:"-"`
}

// LocationTable return location table
//
func (db *DB) LocationTable() *data.Table {
	return db.newTable("location", func() data.ObjectRef {
		return &Location{}
	})
}

// LocationTotal return total location count
//
//	id := d.LocationTotal(ctx)
//
func (c *Counters) LocationTotal(ctx context.Context) (data.CounterRef, error) {
	return c.Counter(ctx, "location-total", 4)
}

// LocationID generate new location serial id
//
//	id := d.LocationID(ctx)
//
func (s *Serial) LocationID(ctx context.Context) (uint32, error) {
	return s.Number(ctx, "location-id")
}
