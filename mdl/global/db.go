package gmdl

import (
	"context"

	data "github.com/piyuo/libsrv/data"
)

// DB represent global database
//
type DB struct {
	data.DB
}

// NewDB create db instance
//
func NewDB(ctx context.Context) (*DB, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	conn, err := data.FirestoreGlobalConnection(ctx, "")
	if err != nil {
		return nil, err
	}

	db := &DB{
		DB: data.DB{Connection: conn},
	}
	return db, nil
}

// Counters return global counters
//
func (db *DB) Counters() *Counters {
	return &Counters{
		Counters: data.Counters{
			Connection: db.Connection,
			TableName:  "counter",
		},
	}
}

// NewTable return table
//
func (db *DB) newTable(name string, factory func() data.ObjectRef) *data.Table {
	return &data.Table{
		Connection: db.Connection,
		TableName:  name,
		Factory:    factory,
	}
}

// Serial return serial
//
func (db *DB) Serial() *Serial {
	return &Serial{
		Serial: data.Serial{
			Connection: db.Connection,
			TableName:  "serial",
		}}
}

// Counters is collection of global usage counters
//
type Counters struct {
	data.Counters `firestore:"-"`
}

// Serial keep serial numbers
//
type Serial struct {
	data.Serial `firestore:"-"`
}
