package rmdl

import (
	"context"

	data "github.com/piyuo/libsrv/data"
)

// Database represent regional database
//
type Database struct {
	data.DocDB
}

// NewDatabase create database
//
func NewDatabase(ctx context.Context, databaseName string) (*Database, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	conn, err := data.FirestoreRegionalConnection(ctx, databaseName)
	if err != nil {
		return nil, err
	}

	db := &Database{}
	db.SetConnection(conn)
	return db, nil
}

// StoreTable return store table
//
func (db *Database) StoreTable() data.Table {
	table := &data.DocTable{}
	table.SetConnection(db.Connection())
	table.SetTableName("store")
	return table
}

// LocationTable return location table
//
func (db *Database) LocationTable() data.Table {
	table := &data.DocTable{}
	table.SetConnection(db.Connection())
	table.SetTableName("location")
	return table
}

// Counters return global counters
//
func (db *Database) Counters() *Counters {
	counters := &Counters{}
	counters.SetConnection(db.Connection())
	counters.SetTableName("counter")
	return counters
}

// Serial return serial
//
func (db *Database) Serial() *Serial {
	serial := &Serial{}
	serial.SetConnection(db.Connection())
	serial.SetTableName("serial")
	return serial
}

// Counters is collection of global usage counters
//
type Counters struct {
	data.Counters `firestore:"-"`
}

// LocationTotal return total location count
//
//	id := d.LocationTotal(ctx)
//
func (c *Counters) LocationTotal(ctx context.Context) (data.Counter, error) {
	return c.Counter(ctx, "location-total", 4)
}

// Serial keep serial numbers
//
type Serial struct {
	data.Serial `firestore:"-"`
}

// LocationID generate new location serial id
//
//	id := d.LocationID(ctx)
//
func (s *Serial) LocationID(ctx context.Context) (uint32, error) {
	return s.Number(ctx, "location-id")
}
